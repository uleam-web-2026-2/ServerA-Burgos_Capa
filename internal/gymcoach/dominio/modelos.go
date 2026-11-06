package dominio

import (
	"time"

	"gorm.io/gorm"
)

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
	Titulo      string  `gorm:"not null"`
	Descripcion string  `gorm:"not null"`
	DiasSemana  int     `gorm:"not null"`
	Estado      string  `gorm:"not null;default:pendiente"`
	ClienteID   uint    `gorm:"not null"`
	Cliente     Cliente `gorm:"foreignKey:ClienteID"`
}
