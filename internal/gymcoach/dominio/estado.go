package dominio

import (
	"errors"
	"slices"
)

var (
	ErrRutinaNoEncontrada = errors.New("rutina no encontrada")
	ErrClienteInvalido    = errors.New("cliente no encontrado")
	ErrEstadoModificado   = errors.New("el estado de la rutina fue modificado")
)

var EstadosValidos = map[string]bool{
	"pendiente":   true,
	"en_progreso": true,
	"completada":  true,
}

var siguientes = map[string][]string{
	"pendiente":   {"en_progreso"},
	"en_progreso": {"completada"},
	"completada":  {},
}

func EstadoValido(estado string) bool {
	return EstadosValidos[estado]
}

func PuedeTransicionar(de, a string) bool {
	return slices.Contains(siguientes[de], a)
}
