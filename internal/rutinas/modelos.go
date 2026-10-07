package rutinas

import "gorm.io/gorm"

// Estados válidos para una rutina (Fase 2b)
var estadosValidos = map[string]bool{
	"pendiente":   true,
	"en_progreso": true,
	"completada":  true,
}

// Cliente es la entidad del lado del uno
type Cliente struct {
	gorm.Model
	Nombre  string
	Rutinas []Rutina // Relación 1 a N
}

// Rutina es la entidad con estados (lado de los muchos)
type Rutina struct {
	gorm.Model
	Titulo    string
	Estado    string
	ClienteID uint // Clave foránea obligatoria para evitar "invalid field found for struct"
}
