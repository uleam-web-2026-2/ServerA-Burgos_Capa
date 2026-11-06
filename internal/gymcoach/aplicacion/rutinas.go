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

func (s *ServicioRutinas) Listar(ctx context.Context, estado string) ([]RutinaSalida, error) {
	rutinas, err := s.repositorio.ListarRutinas(ctx, estado)
	if err != nil {
		return nil, err
	}
	respuesta := make([]RutinaSalida, 0, len(rutinas))
	for _, rutina := range rutinas {
		respuesta = append(respuesta, RutinaSalidaDesdeDominio(rutina))
	}
	return respuesta, nil
}

func (s *ServicioRutinas) Crear(ctx context.Context, entrada any) (RutinaSalida, error) {
	var rutina dominio.Rutina
	switch valor := entrada.(type) {
	case *RutinaEntrada:
		if valor == nil {
			return RutinaSalida{}, errors.New("rutina vacía")
		}
		rutina = valor.AModelo()
	case *dominio.Rutina:
		if valor == nil {
			return RutinaSalida{}, errors.New("rutina vacía")
		}
		rutina = *valor
	case RutinaEntrada:
		rutina = valor.AModelo()
	case RutinaActualizacion:
		rutina = valor.AModelo()
	case dominio.Rutina:
		rutina = valor
	default:
		return RutinaSalida{}, errors.New("tipo de rutina no soportado")
	}
	if rutina.Estado == "" {
		rutina.Estado = "pendiente"
	}
	if err := s.repositorio.CrearRutina(ctx, &rutina); err != nil {
		return RutinaSalida{}, err
	}
	return RutinaSalidaDesdeDominio(rutina), nil
}

func (s *ServicioRutinas) Obtener(ctx context.Context, id uint) (RutinaSalida, error) {
	rutina, err := s.repositorio.ObtenerRutinaConCliente(ctx, id)
	if err != nil {
		return RutinaSalida{}, err
	}
	return RutinaSalidaDesdeDominio(rutina), nil
}

func (s *ServicioRutinas) Actualizar(ctx context.Context, id uint, entrada any) (RutinaSalida, error) {
	var cambios dominio.Rutina
	switch valor := entrada.(type) {
	case *RutinaEntrada:
		if valor == nil {
			return RutinaSalida{}, errors.New("rutina vacía")
		}
		cambios = valor.AModelo()
	case *dominio.Rutina:
		if valor == nil {
			return RutinaSalida{}, errors.New("rutina vacía")
		}
		cambios = *valor
	case RutinaEntrada:
		cambios = valor.AModelo()
	case RutinaActualizacion:
		cambios = valor.AModelo()
	case dominio.Rutina:
		cambios = valor
	default:
		return RutinaSalida{}, errors.New("tipo de rutina no soportado")
	}

	rutina, err := s.repositorio.ActualizarRutina(ctx, id, cambios)
	if err != nil {
		return RutinaSalida{}, err
	}
	return RutinaSalidaDesdeDominio(rutina), nil
}

func (s *ServicioRutinas) Eliminar(ctx context.Context, id uint) (int64, error) {
	return s.repositorio.EliminarRutina(ctx, id)
}

func (s *ServicioRutinas) Cambiar(ctx context.Context, id uint, nuevoEstado string) (RutinaSalida, error) {
	if !dominio.EstadoValido(nuevoEstado) {
		return RutinaSalida{}, ErrEstadoInvalido
	}

	rutina, err := s.repositorio.ObtenerRutina(ctx, id)
	if err != nil {
		return RutinaSalida{}, err
	}
	if !dominio.PuedeTransicionar(rutina.Estado, nuevoEstado) {
		return RutinaSalida{}, &TransicionProhibidaError{De: rutina.Estado, A: nuevoEstado}
	}

	filas, err := s.repositorio.ActualizarEstado(ctx, id, rutina.Estado, nuevoEstado)
	if err != nil {
		return RutinaSalida{}, err
	}
	if filas == 0 {
		return RutinaSalida{}, dominio.ErrEstadoModificado
	}

	actualizada, err := s.repositorio.ObtenerRutinaConCliente(ctx, id)
	if err != nil {
		return RutinaSalida{}, err
	}
	return RutinaSalidaDesdeDominio(actualizada), nil
}
