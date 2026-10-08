package gymcoach

import "gorm.io/gorm"

func Sembrar(db *gorm.DB) error {
	var total int64
	if err := db.Model(&Entrenador{}).Count(&total).Error; err != nil {
		return err
	}
	if total > 0 {
		return nil
	}

	entrenador := Entrenador{Nombre: "Entrenador de ejemplo", Correo: "entrenador@example.test"}
	if err := db.Create(&entrenador).Error; err != nil {
		return err
	}
	clientes := []Cliente{
		{
			Nombre:       "Daniela Vera",
			Correo:       "daniela@example.test",
			Objetivo:     "Mejorar fuerza general",
			EntrenadorID: entrenador.ID,
		},
		{
			Nombre:       "Luis Mera",
			Correo:       "luis@example.test",
			Objetivo:     "Aumentar resistencia",
			EntrenadorID: entrenador.ID,
		},
	}
	if err := db.Create(&clientes).Error; err != nil {
		return err
	}
	rutinas := []Rutina{
		{
			Titulo:      "Fuerza inicial",
			Descripcion: "Rutina de cuerpo completo con cargas progresivas.",
			DiasSemana:  3,
			Estado:      "pendiente",
			ClienteID:   clientes[0].ID,
		},
		{
			Titulo:      "Resistencia base",
			Descripcion: "Trabajo de acondicionamiento y movilidad.",
			DiasSemana:  4,
			Estado:      "en_progreso",
			ClienteID:   clientes[1].ID,
		},
	}
	return db.Create(&rutinas).Error
}
