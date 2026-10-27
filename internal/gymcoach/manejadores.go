package gymcoach

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/aplicacion"
	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/dominio"
	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/respuesta"
	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/sesion"
)

type Manejador struct {
	DB              *gorm.DB
	ServicioRutinas *aplicacion.ServicioRutinas
}

func (m *Manejador) Rutas(r chi.Router, tokens *sesion.Tokens) {
	// Nivel 1 · Pública: cualquiera, sin token.
	// (POST /login se registra en main.go. Aquí no hay rutas públicas.)

	// Nivel 2 · Con sesión: antes del manejador corre el guardia.
	conSesion := r.With(tokens.Autenticar)
	conSesion.Get("/rutinas/{id}", m.verRutina)

	// Nivel 3 · Solo entrenadores: después del guardia se revisa el rol.
	soloEntrenador := conSesion.With(sesion.SoloRol(sesion.RolEntrenador))
	soloEntrenador.Get("/clientes", m.listarClientes)
	soloEntrenador.Post("/clientes", m.crearCliente)
	soloEntrenador.Get("/rutinas", m.listarRutinas)
	soloEntrenador.Post("/rutinas", m.crearRutina)
	soloEntrenador.Put("/rutinas/{id}", m.actualizarRutina)
	soloEntrenador.Delete("/rutinas/{id}", m.borrarRutina)

	// Nivel 3 · Solo clientes: el cambio de estado es de quien sigue la rutina.
	soloCliente := conSesion.With(sesion.SoloRol(sesion.RolCliente))
	soloCliente.Patch("/rutinas/{id}", m.cambiarEstado)
}

func decodificar(r *http.Request, destino any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destino); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("el cuerpo debe contener un solo objeto JSON")
	}
	return nil
}

func idDeRuta(w http.ResponseWriter, r *http.Request) (uint, bool) {
	n, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil || n == 0 {
		respuesta.Error(w, http.StatusBadRequest, "id_invalido", "El id debe ser un entero positivo")
		return 0, false
	}
	return uint(n), true
}

func (m *Manejador) crearCliente(w http.ResponseWriter, r *http.Request) {
	var entrada struct {
		Nombre       string `json:"nombre"`
		Correo       string `json:"correo"`
		Objetivo     string `json:"objetivo"`
		EntrenadorID uint   `json:"entrenador_id"`
	}
	if err := decodificar(r, &entrada); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "json_invalido", "El cuerpo debe ser un único JSON válido con campos conocidos")
		return
	}
	entrada.Nombre = strings.TrimSpace(entrada.Nombre)
	entrada.Correo = strings.TrimSpace(entrada.Correo)
	entrada.Objetivo = strings.TrimSpace(entrada.Objetivo)
	if entrada.Nombre == "" || entrada.Correo == "" || entrada.Objetivo == "" || entrada.EntrenadorID == 0 {
		respuesta.Error(w, http.StatusUnprocessableEntity, "datos_requeridos", "Nombre, correo, objetivo y entrenador_id son obligatorios")
		return
	}
	var entrenador Entrenador
	if err := m.DB.First(&entrenador, entrada.EntrenadorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respuesta.Error(w, http.StatusUnprocessableEntity, "entrenador_invalido", "El entrenador indicado no existe")
			return
		}
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "No se pudo validar el entrenador")
		return
	}
	cliente := Cliente{
		Nombre:       entrada.Nombre,
		Correo:       entrada.Correo,
		Objetivo:     entrada.Objetivo,
		EntrenadorID: entrada.EntrenadorID,
	}
	if err := m.DB.Create(&cliente).Error; err != nil {
		respuesta.Error(w, http.StatusUnprocessableEntity, "cliente_no_guardado", "No se pudo guardar el cliente; verifique que el correo no esté registrado")
		return
	}
	respuesta.Exito(w, http.StatusCreated, cliente)
}

func (m *Manejador) listarClientes(w http.ResponseWriter, r *http.Request) {
	var clientes []Cliente
	query := m.DB.Preload("Entrenador").Preload("Rutinas").Order("id ASC")
	if estado := r.URL.Query().Get("estado"); estado != "" {
		if !estadosValidos[estado] {
			respuesta.Error(w, http.StatusUnprocessableEntity, "estado_invalido", "El filtro de estado de rutina no es válido")
			return
		}
		query = m.DB.Preload("Entrenador").Preload("Rutinas", "estado = ?", estado).Order("id ASC")
	}
	if err := query.Find(&clientes).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "No se pudo obtener la lista de clientes")
		return
	}
	respuesta.Exito(w, http.StatusOK, clientes)
}

func (m *Manejador) listarRutinas(w http.ResponseWriter, r *http.Request) {
	estado := r.URL.Query().Get("estado")
	if estado != "" && !estadosValidos[estado] {
		respuesta.Error(w, http.StatusUnprocessableEntity, "estado_invalido",
			"El filtro de estado de rutina no es válido")
		return
	}
	if m.ServicioRutinas == nil {
		respuesta.Error(w, http.StatusInternalServerError, "servicio_no_configurado", "No se pudo obtener la lista de rutinas")
		return
	}
	rutinas, err := m.ServicioRutinas.Listar(r.Context(), estado)
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base",
			"No se pudo obtener la lista de rutinas")
		return
	}
	respuesta.Exito(w, http.StatusOK, rutinas)
}

func (m *Manejador) crearRutina(w http.ResponseWriter, r *http.Request) {
	var entrada struct {
		Titulo      string `json:"titulo"`
		Descripcion string `json:"descripcion"`
		DiasSemana  int    `json:"dias_semana"`
		ClienteID   uint   `json:"cliente_id"`
	}
	if err := decodificar(r, &entrada); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "json_invalido", "El cuerpo debe ser un único JSON válido con campos conocidos")
		return
	}
	entrada.Titulo = strings.TrimSpace(entrada.Titulo)
	entrada.Descripcion = strings.TrimSpace(entrada.Descripcion)
	if entrada.Titulo == "" || entrada.Descripcion == "" || entrada.ClienteID == 0 {
		respuesta.Error(w, http.StatusUnprocessableEntity, "datos_requeridos", "Título, descripción y cliente_id son obligatorios")
		return
	}
	if entrada.DiasSemana < 1 || entrada.DiasSemana > 7 {
		respuesta.Error(w, http.StatusUnprocessableEntity, "dias_semana_invalidos", "dias_semana debe estar entre 1 y 7")
		return
	}
	rutina := Rutina{
		Titulo:      entrada.Titulo,
		Descripcion: entrada.Descripcion,
		DiasSemana:  entrada.DiasSemana,
		Estado:      "pendiente",
		ClienteID:   entrada.ClienteID,
	}
	if m.ServicioRutinas == nil {
		respuesta.Error(w, http.StatusInternalServerError, "servicio_no_configurado", "No se pudo guardar la rutina")
		return
	}
	if err := m.ServicioRutinas.Crear(r.Context(), &rutina); errors.Is(err, dominio.ErrClienteInvalido) {
		respuesta.Error(w, http.StatusUnprocessableEntity, "cliente_invalido", "El cliente indicado no existe")
		return
	} else if err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "No se pudo guardar la rutina")
		return
	}
	respuesta.Exito(w, http.StatusCreated, rutina)
}

func (m *Manejador) verRutina(w http.ResponseWriter, r *http.Request) {
	id, ok := idDeRuta(w, r)
	if !ok {
		return
	}
	if m.ServicioRutinas == nil {
		respuesta.Error(w, http.StatusInternalServerError, "servicio_no_configurado", "Error al buscar la rutina")
		return
	}
	rutina, err := m.ServicioRutinas.Obtener(r.Context(), id)
	if errors.Is(err, dominio.ErrRutinaNoEncontrada) {
		respuesta.Error(w, http.StatusNotFound, "no_encontrada", "Rutina no encontrada")
		return
	}
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "Error al buscar la rutina")
		return
	}
	respuesta.Exito(w, http.StatusOK, rutina)
}

func (m *Manejador) actualizarRutina(w http.ResponseWriter, r *http.Request) {
	id, ok := idDeRuta(w, r)
	if !ok {
		return
	}
	if m.ServicioRutinas == nil {
		respuesta.Error(w, http.StatusInternalServerError, "servicio_no_configurado", "Error al buscar la rutina")
		return
	}
	if _, err := m.ServicioRutinas.Obtener(r.Context(), id); err != nil {
		if errors.Is(err, dominio.ErrRutinaNoEncontrada) {
			respuesta.Error(w, http.StatusNotFound, "no_encontrada", "Rutina no encontrada")
			return
		}
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "Error al buscar la rutina")
		return
	}
	var entrada struct {
		Titulo      string `json:"titulo"`
		Descripcion string `json:"descripcion"`
		DiasSemana  int    `json:"dias_semana"`
		Estado      string `json:"estado"` // se acepta pero se ignora: el estado solo cambia con PATCH
		ClienteID   uint   `json:"cliente_id"`
	}
	if err := decodificar(r, &entrada); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "json_invalido", "El cuerpo debe ser un único JSON válido con campos conocidos")
		return
	}
	entrada.Titulo = strings.TrimSpace(entrada.Titulo)
	entrada.Descripcion = strings.TrimSpace(entrada.Descripcion)
	if entrada.Titulo == "" || entrada.Descripcion == "" || entrada.ClienteID == 0 {
		respuesta.Error(w, http.StatusUnprocessableEntity, "datos_requeridos", "Título, descripción y cliente_id son obligatorios")
		return
	}
	if entrada.DiasSemana < 1 || entrada.DiasSemana > 7 {
		respuesta.Error(w, http.StatusUnprocessableEntity, "dias_semana_invalidos", "dias_semana debe estar entre 1 y 7")
		return
	}
	actualizada, err := m.ServicioRutinas.Actualizar(r.Context(), id, Rutina{
		Titulo:      entrada.Titulo,
		Descripcion: entrada.Descripcion,
		DiasSemana:  entrada.DiasSemana,
		ClienteID:   entrada.ClienteID,
	})
	if errors.Is(err, dominio.ErrClienteInvalido) {
		respuesta.Error(w, http.StatusUnprocessableEntity, "cliente_invalido", "El cliente indicado no existe")
		return
	}
	if errors.Is(err, dominio.ErrRutinaNoEncontrada) {
		respuesta.Error(w, http.StatusNotFound, "no_encontrada", "Rutina no encontrada")
		return
	}
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "No se pudo actualizar la rutina")
		return
	}
	respuesta.Exito(w, http.StatusOK, actualizada)
}

func (m *Manejador) borrarRutina(w http.ResponseWriter, r *http.Request) {
	id, ok := idDeRuta(w, r)
	if !ok {
		return
	}
	if m.ServicioRutinas == nil {
		respuesta.Error(w, http.StatusInternalServerError, "servicio_no_configurado", "No se pudo eliminar la rutina")
		return
	}
	filas, err := m.ServicioRutinas.Eliminar(r.Context(), id)
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "No se pudo eliminar la rutina")
		return
	}
	if filas == 0 {
		respuesta.Error(w, http.StatusNotFound, "no_encontrada", "Rutina no encontrada")
		return
	}
	respuesta.Exito(w, http.StatusOK, map[string]string{"mensaje": "Rutina eliminada"})
}
