package rutinas

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/respuesta"
)

// Manejador guarda la conexión. Todas las rutas la leen con m.DB.
type Manejador struct {
	DB *gorm.DB
}

// Rutas registra las rutas del paquete rutinas
func (m *Manejador) Rutas(r chi.Router) {
	r.Post("/rutinas", m.crear)
	r.Get("/rutinas", m.listar)
	r.Get("/rutinas/{id}", m.verUno)
	r.Put("/rutinas/{id}", m.actualizar)
	r.Delete("/rutinas/{id}", m.borrar)

	// Listado del lado del uno (Clientes) con Preload para evitar N+1 (Fase 2c)
	r.Get("/clientes", m.listarClientesConRutinas)
}

// leerID saca el id de la URL y lo convierte a número de forma segura.
func leerID(w http.ResponseWriter, r *http.Request) (uint, bool) {
	n, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || n <= 0 {
		respuesta.Error(w, http.StatusBadRequest, "id_invalido", "El id debe ser un número positivo")
		return 0, false
	}
	return uint(n), true
}

// crear (POST /rutinas)
func (m *Manejador) crear(w http.ResponseWriter, r *http.Request) {
	var rut Rutina
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&rut); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "json_invalido", "El cuerpo no es un JSON válido")
		return
	}
	rut.ID = 0 // Evitar que el usuario inyecte su propio ID

	// Validar estado (Fase 2b)
	if !estadosValidos[rut.Estado] {
		respuesta.Error(w, http.StatusUnprocessableEntity, "estado_invalido", "Estado no válido")
		return
	}

	// Regla de negocio extra: El título no puede estar vacío (Fase 2b)
	if rut.Titulo == "" {
		respuesta.Error(w, http.StatusUnprocessableEntity, "titulo_requerido", "El título de la rutina es obligatorio")
		return
	}

	// Guardar en la base de datos con depuración SQL
	if err := m.DB.Debug().Create(&rut).Error; err != nil {
		// Si mandan un ClienteID que no existe, GORM arroja foreign key constraint violation -> 422
		respuesta.Error(w, http.StatusUnprocessableEntity, "error_base", "No se pudo guardar, verifique que el cliente exista")
		return
	}

	respuesta.Exito(w, http.StatusCreated, rut)
}

// listar (GET /rutinas con soporte para filtrar por ?estado de forma segura contra inyección)
func (m *Manejador) listar(w http.ResponseWriter, r *http.Request) {
	var rutinas []Rutina
	estado := r.URL.Query().Get("estado")

	dbQuery := m.DB.Debug()
	if estado != "" {
		// Uso de parámetros ? para evitar inyección SQL (Fase 2d)
		dbQuery = dbQuery.Where("estado = ?", estado)
	}

	if err := dbQuery.Find(&rutinas).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "No se pudo obtener la lista")
		return
	}

	respuesta.Exito(w, http.StatusOK, rutinas)
}

// verUno (GET /rutinas/{id})
func (m *Manejador) verUno(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(w, r)
	if !ok {
		return
	}

	var rut Rutina
	err := m.DB.Debug().First(&rut, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respuesta.Error(w, http.StatusNotFound, "no_encontrado", "Rutina no encontrada")
			return
		}
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "Error al buscar la rutina")
		return
	}

	respuesta.Exito(w, http.StatusOK, rut)
}

// actualizar (PUT /rutinas/{id})
func (m *Manejador) actualizar(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(w, r)
	if !ok {
		return
	}

	// 1. Buscar si existe el registro primero (evita crear uno nuevo por error con Save)
	var existente Rutina
	if err := m.DB.First(&existente, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respuesta.Error(w, http.StatusNotFound, "no_encontrado", "Rutina no encontrada")
			return
		}
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "Error de base de datos")
		return
	}

	// 2. Leer el JSON entrante
	var entrada Rutina
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&entrada); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "json_invalido", "El cuerpo no es un JSON válido")
		return
	}

	// Validar estado
	if !estadosValidos[entrada.Estado] {
		respuesta.Error(w, http.StatusUnprocessableEntity, "estado_invalido", "Estado no válido")
		return
	}

	// 3. Copiar los campos permitidos a cambiar
	existente.Titulo = entrada.Titulo
	existente.Estado = entrada.Estado
	existente.ClienteID = entrada.ClienteID

	if err := m.DB.Debug().Save(&existente).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "No se pudo actualizar")
		return
	}

	respuesta.Exito(w, http.StatusOK, existente)
}

// borrar (DELETE /rutinas/{id})
func (m *Manejador) borrar(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(w, r)
	if !ok {
		return
	}

	res := m.DB.Debug().Delete(&Rutina{}, id)
	if res.Error != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "No se pudo eliminar")
		return
	}

	if res.RowsAffected == 0 {
		respuesta.Error(w, http.StatusNotFound, "no_encontrado", "Rutina no encontrada")
		return
	}

	respuesta.Exito(w, http.StatusOK, map[string]string{"mensaje": "Eliminado correctamente"})
}

// listarClientesConRutinas (GET /clientes) - Preload para evitar problema N+1 (Fase 2c)
func (m *Manejador) listarClientesConRutinas(w http.ResponseWriter, r *http.Request) {
	var clientes []Cliente
	// Precarga la relación "Rutinas" haciendo exactamente dos consultas SQL en total
	if err := m.DB.Debug().Preload("Rutinas").Find(&clientes).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "Error al listar clientes")
		return
	}
	respuesta.Exito(w, http.StatusOK, clientes)
}
