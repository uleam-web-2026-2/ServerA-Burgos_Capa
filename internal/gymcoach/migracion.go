package gymcoach

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func MigrarDatosAntiguos(db *gorm.DB) (int64, error) {
	var registrosActualizados int64
	err := db.Transaction(func(tx *gorm.DB) error {
		if tx.Migrator().HasTable(&Cliente{}) {
			if err := tx.Exec(`
				ALTER TABLE "clientes"
					ADD COLUMN IF NOT EXISTS "correo" text,
					ADD COLUMN IF NOT EXISTS "objetivo" text,
					ADD COLUMN IF NOT EXISTS "entrenador_id" bigint
			`).Error; err != nil {
				return err
			}

			for _, consulta := range []string{
				`UPDATE "clientes"
				 SET "correo" = 'cliente-migrado-' || "id"::text || '@example.invalid'
				 WHERE "correo" IS NULL OR BTRIM("correo") = ''`,
				`UPDATE "clientes"
				 SET "objetivo" = 'Objetivo pendiente de actualizar'
				 WHERE "objetivo" IS NULL OR BTRIM("objetivo") = ''`,
			} {
				resultado := tx.Exec(consulta)
				if resultado.Error != nil {
					return resultado.Error
				}
				registrosActualizados += resultado.RowsAffected
			}

			var entrenadoresFaltantes int64
			if err := tx.Raw(`
				SELECT COUNT(*)
				FROM "clientes" c
				LEFT JOIN "entrenadores" e ON e."id" = c."entrenador_id"
				WHERE c."entrenador_id" IS NULL OR e."id" IS NULL
			`).Scan(&entrenadoresFaltantes).Error; err != nil {
				return err
			}
			if entrenadoresFaltantes > 0 {
				entrenador := Entrenador{
					Nombre: "Entrenador pendiente de asignar",
					Correo: "entrenador-migrado@example.invalid",
				}
				if err := tx.Where("correo = ?", entrenador.Correo).
					FirstOrCreate(&entrenador).Error; err != nil {
					return err
				}
				resultado := tx.Exec(`
					UPDATE "clientes" c
					SET "entrenador_id" = ?
					WHERE c."entrenador_id" IS NULL
					   OR NOT EXISTS (
							SELECT 1 FROM "entrenadores" e
							WHERE e."id" = c."entrenador_id"
					   )
				`, entrenador.ID)
				if resultado.Error != nil {
					return resultado.Error
				}
				registrosActualizados += resultado.RowsAffected
			}
		}

		if tx.Migrator().HasTable(&Rutina{}) {
			if err := tx.Exec(`
				ALTER TABLE "rutinas"
					ADD COLUMN IF NOT EXISTS "titulo" text,
					ADD COLUMN IF NOT EXISTS "descripcion" text,
					ADD COLUMN IF NOT EXISTS "dias_semana" integer,
					ADD COLUMN IF NOT EXISTS "estado" text,
					ADD COLUMN IF NOT EXISTS "cliente_id" bigint
			`).Error; err != nil {
				return err
			}

			for _, consulta := range []string{
				`UPDATE "rutinas"
				 SET "titulo" = 'Rutina migrada ' || "id"::text
				 WHERE "titulo" IS NULL OR BTRIM("titulo") = ''`,
				`UPDATE "rutinas"
				 SET "descripcion" = 'Descripción pendiente de actualizar'
				 WHERE "descripcion" IS NULL OR BTRIM("descripcion") = ''`,
				`UPDATE "rutinas" SET "dias_semana" = 3 WHERE "dias_semana" IS NULL OR "dias_semana" < 1 OR "dias_semana" > 7`,
				`UPDATE "rutinas" SET "estado" = 'pendiente' WHERE "estado" IS NULL OR "estado" NOT IN ('pendiente', 'en_progreso', 'completada')`,
			} {
				resultado := tx.Exec(consulta)
				if resultado.Error != nil {
					return resultado.Error
				}
				registrosActualizados += resultado.RowsAffected
			}

			var rutinasSinCliente int64
			if err := tx.Raw(`
				SELECT COUNT(*)
				FROM "rutinas" r
				LEFT JOIN "clientes" c ON c."id" = r."cliente_id"
				WHERE r."cliente_id" IS NULL OR c."id" IS NULL
			`).Scan(&rutinasSinCliente).Error; err != nil {
				return err
			}
			if rutinasSinCliente > 0 {
				var cliente Cliente
				err := tx.Order("id ASC").First(&cliente).Error
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				if errors.Is(err, gorm.ErrRecordNotFound) {
					entrenador := Entrenador{
						Nombre: "Entrenador pendiente de asignar",
						Correo: "entrenador-migrado@example.invalid",
					}
					if err := tx.Where("correo = ?", entrenador.Correo).
						FirstOrCreate(&entrenador).Error; err != nil {
						return err
					}
					cliente = Cliente{
						Nombre:       "Cliente pendiente de asignar",
						Correo:       "cliente-migrado@example.invalid",
						Objetivo:     "Objetivo pendiente de actualizar",
						EntrenadorID: entrenador.ID,
					}
					if err := tx.Create(&cliente).Error; err != nil {
						return err
					}
				}
				resultado := tx.Exec(`
					UPDATE "rutinas" r
					SET "cliente_id" = ?
					WHERE r."cliente_id" IS NULL
					   OR NOT EXISTS (
							SELECT 1 FROM "clientes" c
							WHERE c."id" = r."cliente_id"
					   )
				`, cliente.ID)
				if resultado.Error != nil {
					return resultado.Error
				}
				registrosActualizados += resultado.RowsAffected
			}
		}
		return nil
	})
	if err != nil {
		return registrosActualizados, fmt.Errorf("preparación de datos antiguos: %w", err)
	}
	return registrosActualizados, nil
}
