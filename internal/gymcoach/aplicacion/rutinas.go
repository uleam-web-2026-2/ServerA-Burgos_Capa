package aplicacion

import (
	"context"
	"errors"
	"fmt"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/dominio"
)

var ErrEstadoInvalido = errors.New("estado no permitido")

type TransicionProhibidaError struct {
	De string
	A  string
}

func (e *TransicionProhibidaError) Error() string {
	return fmt.Sprintf("de %s no se puede pasar a %s", e.De, e.A)
}

type RepositorioRutinas interface {
	ListarRutinas(context.Context, string) ([]dominio.Rutina, error)
	CrearRutina(context.Context, *dominio.Rutina) error
	ObtenerRutina(context.Context, uint) (dominio.Rutina, error)
	ObtenerRutinaConCliente(context.Context, uint) (dominio.Rutina, error)
	ActualizarRutina(context.Context, uint, dominio.Rutina) (dominio.Rutina, error)
	EliminarRutina(context.Context, uint) (int64, error)
	ActualizarEstado(context.Context, uint, string, string) (int64, error)
}

type ServicioRutinas struct {
	repositorio RepositorioRutinas
}

func NuevoServicioRutinas(repositorio RepositorioRutinas) *ServicioRutinas {
	return &ServicioRutinas{repositorio: repositorio}
}

func (s *ServicioRutinas) Listar(ctx context.Context, estado string) ([]dominio.Rutina, error) {
	return s.repositorio.ListarRutinas(ctx, estado)
}

func (s *ServicioRutinas) Crear(ctx context.Context, rutina *dominio.Rutina) error {
	return s.repositorio.CrearRutina(ctx, rutina)
}

func (s *ServicioRutinas) Obtener(ctx context.Context, id uint) (dominio.Rutina, error) {
	return s.repositorio.ObtenerRutinaConCliente(ctx, id)
}

func (s *ServicioRutinas) Actualizar(ctx context.Context, id uint, rutina dominio.Rutina) (dominio.Rutina, error) {
	return s.repositorio.ActualizarRutina(ctx, id, rutina)
}

func (s *ServicioRutinas) Eliminar(ctx context.Context, id uint) (int64, error) {
	return s.repositorio.EliminarRutina(ctx, id)
}

func (s *ServicioRutinas) Cambiar(ctx context.Context, id uint, nuevoEstado string) (dominio.Rutina, error) {
	if !dominio.EstadoValido(nuevoEstado) {
		return dominio.Rutina{}, ErrEstadoInvalido
	}

	rutina, err := s.repositorio.ObtenerRutina(ctx, id)
	if err != nil {
		return dominio.Rutina{}, err
	}
	if !dominio.PuedeTransicionar(rutina.Estado, nuevoEstado) {
		return dominio.Rutina{}, &TransicionProhibidaError{De: rutina.Estado, A: nuevoEstado}
	}

	filas, err := s.repositorio.ActualizarEstado(ctx, id, rutina.Estado, nuevoEstado)
	if err != nil {
		return dominio.Rutina{}, err
	}
	if filas == 0 {
		return dominio.Rutina{}, dominio.ErrEstadoModificado
	}

	return s.repositorio.ObtenerRutinaConCliente(ctx, id)
}
