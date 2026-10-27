package gymcoach

import (
	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/dominio"
	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/sesion"
)

const (
	RolEntrenador = sesion.RolEntrenador
	RolCliente    = sesion.RolCliente
)

type BaseModel = dominio.BaseModel
type Entrenador = dominio.Entrenador
type Cliente = dominio.Cliente
type Rutina = dominio.Rutina

var estadosValidos = dominio.EstadosValidos

// puedeTransicionar delega en la máquina de estados de estado.go:
// una sola fuente de verdad para las transiciones.
func puedeTransicionar(de, a string) bool {
	return dominio.PuedeTransicionar(de, a)
}
