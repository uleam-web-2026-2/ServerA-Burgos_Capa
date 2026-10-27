package gymcoach

import (
	"errors"
	"net/http"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/aplicacion"
	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/dominio"
	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/respuesta"
)

// cambiarEstado: PATCH /rutinas/{id} con el cuerpo {"estado": "..."}.
// 400 id o JSON · 422 estado que no existe · 404 no existe
// · 422 transición prohibida · 200 cambiado.
func (m *Manejador) cambiarEstado(w http.ResponseWriter, r *http.Request) {
	id, ok := idDeRuta(w, r)
	if !ok {
		return
	}
	var pedido struct {
		Estado string `json:"estado"`
	}
	if err := decodificar(r, &pedido); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "json_invalido",
			"El cuerpo debe ser un único JSON válido con campos conocidos")
		return
	}
	if !dominio.EstadoValido(pedido.Estado) {
		respuesta.Error(w, http.StatusUnprocessableEntity, "estado_invalido",
			"estado no permitido: "+pedido.Estado)
		return
	}
	if m.ServicioRutinas == nil {
		respuesta.Error(w, http.StatusInternalServerError, "servicio_no_configurado",
			"No se pudo cambiar el estado de la rutina")
		return
	}

	rutina, err := m.ServicioRutinas.Cambiar(r.Context(), id, pedido.Estado)
	var transicionProhibida *aplicacion.TransicionProhibidaError
	switch {
	case errors.Is(err, dominio.ErrRutinaNoEncontrada):
		respuesta.Error(w, http.StatusNotFound, "no_encontrada", "Rutina no encontrada")
		return
	case errors.As(err, &transicionProhibida):
		respuesta.Error(w, http.StatusUnprocessableEntity, "transicion_prohibida",
			transicionProhibida.Error())
		return
	case errors.Is(err, dominio.ErrEstadoModificado):
		respuesta.Error(w, http.StatusConflict, "estado_modificado",
			"La rutina cambió mientras se procesaba; vuelva a consultarla")
		return
	case errors.Is(err, aplicacion.ErrEstadoInvalido):
		respuesta.Error(w, http.StatusUnprocessableEntity, "estado_invalido",
			"estado no permitido: "+pedido.Estado)
		return
	case err != nil:
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "No se pudo cambiar el estado de la rutina")
		return
	}
	respuesta.Exito(w, http.StatusOK, rutina)
}
