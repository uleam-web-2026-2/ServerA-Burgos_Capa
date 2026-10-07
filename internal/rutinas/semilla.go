package rutinas

import "gorm.io/gorm"

// Sembrar carga datos de ejemplo solo si la tabla está vacía.
func Sembrar(db *gorm.DB) {
	var total int64
	if err := db.Model(&Cliente{}).Count(&total).Error; err != nil {
		panic(err)
	}
	if total > 0 {
		return
	}

	datos := []Cliente{
		{
			Nombre: "Carlos Mendoza",
			Rutinas: []Rutina{
				{Titulo: "Rutina de Pecho y Tríceps", Estado: "pendiente"},
				{Titulo: "Rutina de Pierna Completa", Estado: "en_progreso"},
			},
		},
		{
			Nombre: "María Fernanda",
			Rutinas: []Rutina{
				{Titulo: "Cardio y Abdomen", Estado: "completada"},
			},
		},
	}

	if err := db.Create(&datos).Error; err != nil {
		panic(err)
	}
}
