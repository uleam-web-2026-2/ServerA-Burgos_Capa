package aplicacion

import (
	"time"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/dominio"
)

type RutinaEntradaDTO struct {
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	DiasSemana  int    `json:"dias_semana"`
	ClienteID   uint   `json:"cliente_id"`
}

type ActualizarRutinaDTO struct {
	RutinaEntradaDTO
	Estado string `json:"estado"`
}

type ClienteDTO struct {
	ID           uint      `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Nombre       string    `json:"nombre"`
	Correo       string    `json:"correo"`
	Objetivo     string    `json:"objetivo"`
	EntrenadorID uint      `json:"entrenador_id"`
}

type RutinaDTO struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Titulo      string     `json:"titulo"`
	Descripcion string     `json:"descripcion"`
	DiasSemana  int        `json:"dias_semana"`
	Estado      string     `json:"estado"`
	ClienteID   uint       `json:"cliente_id"`
	Cliente     *ClienteDTO `json:"cliente,omitempty"`
}

func rutinaDTO(rutina dominio.Rutina) RutinaDTO {
	resultado := RutinaDTO{
		ID:          rutina.ID,
		CreatedAt:   rutina.CreatedAt,
		UpdatedAt:   rutina.UpdatedAt,
		Titulo:      rutina.Titulo,
		Descripcion: rutina.Descripcion,
		DiasSemana:  rutina.DiasSemana,
		Estado:      rutina.Estado,
		ClienteID:   rutina.ClienteID,
	}
	if rutina.Cliente.ID != 0 {
		resultado.Cliente = &ClienteDTO{
			ID:           rutina.Cliente.ID,
			CreatedAt:    rutina.Cliente.CreatedAt,
			UpdatedAt:    rutina.Cliente.UpdatedAt,
			Nombre:       rutina.Cliente.Nombre,
			Correo:       rutina.Cliente.Correo,
			Objetivo:     rutina.Cliente.Objetivo,
			EntrenadorID: rutina.Cliente.EntrenadorID,
		}
	}
	return resultado
}
