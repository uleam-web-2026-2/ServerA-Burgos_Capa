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

func (s *ServicioRutinas) Listar(ctx context.Context, estado string) ([]RutinaDTO, error) {
	rutinas, err := s.repositorio.ListarRutinas(ctx, estado)
	if err != nil {
		return nil, err
	}
	resultado := make([]RutinaDTO, 0, len(rutinas))
	for _, rutina := range rutinas {
		resultado = append(resultado, rutinaDTO(rutina))
	}
	return resultado, nil
}

func (s *ServicioRutinas) Crear(ctx context.Context, entrada RutinaEntradaDTO) (RutinaDTO, error) {
	rutina := dominio.Rutina{
		Titulo:      entrada.Titulo,
		Descripcion: entrada.Descripcion,
		DiasSemana:  entrada.DiasSemana,
		Estado:      "pendiente",
		ClienteID:   entrada.ClienteID,
	}
	if err := s.repositorio.CrearRutina(ctx, &rutina); err != nil {
		return RutinaDTO{}, err
	}
	return rutinaDTO(rutina), nil
}

func (s *ServicioRutinas) Obtener(ctx context.Context, id uint) (RutinaDTO, error) {
	rutina, err := s.repositorio.ObtenerRutinaConCliente(ctx, id)
	if err != nil {
		return RutinaDTO{}, err
	}
	return rutinaDTO(rutina), nil
}

func (s *ServicioRutinas) Actualizar(ctx context.Context, id uint, entrada RutinaEntradaDTO) (RutinaDTO, error) {
	rutina, err := s.repositorio.ActualizarRutina(ctx, id, dominio.Rutina{
		Titulo:      entrada.Titulo,
		Descripcion: entrada.Descripcion,
		DiasSemana:  entrada.DiasSemana,
		ClienteID:   entrada.ClienteID,
	})
	if err != nil {
		return RutinaDTO{}, err
	}
	return rutinaDTO(rutina), nil
}

func (s *ServicioRutinas) Eliminar(ctx context.Context, id uint) (int64, error) {
	return s.repositorio.EliminarRutina(ctx, id)
}

func (s *ServicioRutinas) Cambiar(ctx context.Context, id uint, nuevoEstado string) (RutinaDTO, error) {
	if !dominio.EstadoValido(nuevoEstado) {
		return RutinaDTO{}, ErrEstadoInvalido
	}

	rutina, err := s.repositorio.ObtenerRutina(ctx, id)
	if err != nil {
		return RutinaDTO{}, err
	}
	if !dominio.PuedeTransicionar(rutina.Estado, nuevoEstado) {
		return RutinaDTO{}, &TransicionProhibidaError{De: rutina.Estado, A: nuevoEstado}
	}

	filas, err := s.repositorio.ActualizarEstado(ctx, id, rutina.Estado, nuevoEstado)
	if err != nil {
		return RutinaDTO{}, err
	}
	if filas == 0 {
		return RutinaDTO{}, dominio.ErrEstadoModificado
	}

	rutina, err = s.repositorio.ObtenerRutinaConCliente(ctx, id)
	if err != nil {
		return RutinaDTO{}, err
	}
	return rutinaDTO(rutina), nil
}
