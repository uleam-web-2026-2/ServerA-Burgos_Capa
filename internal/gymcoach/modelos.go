package gymcoach

import (
	"time"

	"gorm.io/gorm"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/sesion"
)

const (
	RolEntrenador = sesion.RolEntrenador
	RolCliente    = sesion.RolCliente
)

var estadosValidos = map[string]bool{
	"pendiente":   true,
	"en_progreso": true,
	"completada":  true,
}

type BaseModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Entrenador struct {
	BaseModel
	Nombre   string    `json:"nombre" gorm:"not null"`
	Correo   string    `json:"correo" gorm:"uniqueIndex;not null"`
	Clientes []Cliente `json:"clientes,omitempty" gorm:"foreignKey:EntrenadorID"`
}

func (Entrenador) TableName() string {
	return "entrenadores"
}

type Cliente struct {
	BaseModel
	Nombre       string     `json:"nombre" gorm:"not null"`
	Correo       string     `json:"correo" gorm:"uniqueIndex;not null"`
	Objetivo     string     `json:"objetivo" gorm:"not null"`
	EntrenadorID uint       `json:"entrenador_id" gorm:"not null"`
	Entrenador   Entrenador `json:"entrenador,omitempty" gorm:"foreignKey:EntrenadorID"`
	Rutinas      []Rutina   `json:"rutinas,omitempty" gorm:"foreignKey:ClienteID"`
}

type Rutina struct {
	BaseModel
	Titulo      string  `json:"titulo" gorm:"not null"`
	Descripcion string  `json:"descripcion" gorm:"not null"`
	DiasSemana  int     `json:"dias_semana" gorm:"not null"`
	Estado      string  `json:"estado" gorm:"not null;default:pendiente"`
	ClienteID   uint    `json:"cliente_id" gorm:"not null"`
	Cliente     Cliente `json:"cliente,omitempty" gorm:"foreignKey:ClienteID"`
}

// puedeTransicionar delega en la máquina de estados de estado.go:
// una sola fuente de verdad para las transiciones.
func puedeTransicionar(de, a string) bool {
	return puedePasar(de, a)
}
