package gymcoach

import (
	"errors"
	"net/http"
	"slices"

	"gorm.io/gorm"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/respuesta"
)

// siguientes es la máquina de estados de la Rutina (sección 5 de la
// ficha): desde cada estado, a cuáles se puede pasar. Van TODOS los
// estados como clave; el estado final lleva la lista vacía.
var siguientes = map[string][]string{
	"pendiente":   {"en_progreso"},
	"en_progreso": {"completada"},
	"completada":  {},
}

// puedePasar dice si la máquina de estados permite ir de un estado a otro.
func puedePasar(de, a string) bool {
	return slices.Contains(siguientes[de], a)
}

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
	if _, existe := siguientes[pedido.Estado]; !existe {
		respuesta.Error(w, http.StatusUnprocessableEntity, "estado_invalido",
			"estado no permitido: "+pedido.Estado)
		return
	}

	// ---- desde aquí hace falta la base de datos ----
	var rutina Rutina
	err := m.DB.First(&rutina, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respuesta.Error(w, http.StatusNotFound, "no_encontrada", "Rutina no encontrada")
		return
	}
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "Error al buscar la rutina")
		return
	}

	// La regla del negocio: se compara el estado que TIENE con el que PIDEN.
	if !puedePasar(rutina.Estado, pedido.Estado) {
		respuesta.Error(w, http.StatusUnprocessableEntity, "transicion_prohibida",
			"de "+rutina.Estado+" no se puede pasar a "+pedido.Estado)
		return
	}

	// Solo se guarda si el estado sigue siendo el que se leyó.
	cambio := m.DB.Model(&Rutina{}).
		Where("id = ? AND estado = ?", rutina.ID, rutina.Estado).
		Update("estado", pedido.Estado)
	if cambio.Error != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "No se pudo guardar el estado")
		return
	}
	if cambio.RowsAffected == 0 {
		respuesta.Error(w, http.StatusConflict, "estado_modificado",
			"La rutina cambió mientras se procesaba; vuelva a consultarla")
		return
	}
	if err := m.DB.Preload("Cliente").First(&rutina, id).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "No se pudo consultar la rutina actualizada")
		return
	}
	respuesta.Exito(w, http.StatusOK, rutina)
}
