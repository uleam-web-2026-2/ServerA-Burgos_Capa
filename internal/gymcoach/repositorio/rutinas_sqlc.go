package repositorio

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/dominio"
	consultassqlc "github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/repositorio/sqlc"
)

type RutinasSqlc struct {
	*Rutinas
	consultas *consultassqlc.Queries
}

func NuevoRutinasSqlc(base *Rutinas, consultas *consultassqlc.Queries) *RutinasSqlc {
	return &RutinasSqlc{Rutinas: base, consultas: consultas}
}

func (r *RutinasSqlc) ObtenerRutina(ctx context.Context, id uint) (dominio.Rutina, error) {
	fila, err := r.consultas.ObtenerRutina(ctx, int64(id))
	if errors.Is(err, sql.ErrNoRows) {
		return dominio.Rutina{}, dominio.ErrRutinaNoEncontrada
	}
	if err != nil {
		return dominio.Rutina{}, err
	}

	return dominio.Rutina{
		BaseModel: dominio.BaseModel{
			ID:        uint(fila.ID),
			CreatedAt: fila.CreatedAt.Time,
			UpdatedAt: fila.UpdatedAt.Time,
		},
		Titulo:      fila.Titulo,
		Descripcion: fila.Descripcion,
		DiasSemana:  int(fila.DiasSemana),
		Estado:      fila.Estado,
		ClienteID:   uint(fila.ClienteID),
	}, nil
}

func (r *RutinasSqlc) ActualizarEstado(ctx context.Context, id uint, estadoActual, nuevoEstado string) (int64, error) {
	return r.consultas.ActualizarEstado(ctx, consultassqlc.ActualizarEstadoParams{
		NuevoEstado:  nuevoEstado,
		ID:           int64(id),
		EstadoActual: estadoActual,
	})
}
