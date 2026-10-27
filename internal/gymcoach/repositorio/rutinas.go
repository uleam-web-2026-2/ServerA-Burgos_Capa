package repositorio

import (
	"context"
	"errors"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/dominio"
	"gorm.io/gorm"
)

type Rutinas struct {
	DB *gorm.DB
}

func (r *Rutinas) ListarRutinas(ctx context.Context, estado string) ([]dominio.Rutina, error) {
	var rutinas []dominio.Rutina
	query := r.DB.WithContext(ctx).Preload("Cliente").Order("id ASC")
	if estado != "" {
		query = query.Where("estado = ?", estado)
	}
	err := query.Find(&rutinas).Error
	return rutinas, err
}

func (r *Rutinas) CrearRutina(ctx context.Context, rutina *dominio.Rutina) error {
	var cliente dominio.Cliente
	err := r.DB.WithContext(ctx).First(&cliente, rutina.ClienteID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dominio.ErrClienteInvalido
	}
	if err != nil {
		return err
	}
	return r.DB.WithContext(ctx).Create(rutina).Error
}

func (r *Rutinas) ObtenerRutina(ctx context.Context, id uint) (dominio.Rutina, error) {
	var rutina dominio.Rutina
	err := r.DB.WithContext(ctx).First(&rutina, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dominio.Rutina{}, dominio.ErrRutinaNoEncontrada
	}
	return rutina, err
}

func (r *Rutinas) ObtenerRutinaConCliente(ctx context.Context, id uint) (dominio.Rutina, error) {
	var rutina dominio.Rutina
	err := r.DB.WithContext(ctx).Preload("Cliente").First(&rutina, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dominio.Rutina{}, dominio.ErrRutinaNoEncontrada
	}
	return rutina, err
}

func (r *Rutinas) ActualizarRutina(ctx context.Context, id uint, cambios dominio.Rutina) (dominio.Rutina, error) {
	var cliente dominio.Cliente
	err := r.DB.WithContext(ctx).First(&cliente, cambios.ClienteID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dominio.Rutina{}, dominio.ErrClienteInvalido
	}
	if err != nil {
		return dominio.Rutina{}, err
	}

	resultado := r.DB.WithContext(ctx).Model(&dominio.Rutina{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"titulo":      cambios.Titulo,
			"descripcion": cambios.Descripcion,
			"dias_semana": cambios.DiasSemana,
			"cliente_id":  cambios.ClienteID,
		})
	if resultado.Error != nil {
		return dominio.Rutina{}, resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return dominio.Rutina{}, dominio.ErrRutinaNoEncontrada
	}
	return r.ObtenerRutinaConCliente(ctx, id)
}

func (r *Rutinas) ActualizarEstado(ctx context.Context, id uint, estadoActual, nuevoEstado string) (int64, error) {
	resultado := r.DB.WithContext(ctx).Model(&dominio.Rutina{}).
		Where("id = ? AND estado = ?", id, estadoActual).
		Update("estado", nuevoEstado)
	return resultado.RowsAffected, resultado.Error
}

func (r *Rutinas) EliminarRutina(ctx context.Context, id uint) (int64, error) {
	resultado := r.DB.WithContext(ctx).Delete(&dominio.Rutina{}, id)
	return resultado.RowsAffected, resultado.Error
}
