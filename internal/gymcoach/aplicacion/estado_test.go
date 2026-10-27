package aplicacion

import (
	"context"
	"errors"
	"testing"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/dominio"
)

type repositorioEstadoPrueba struct {
	rutina    dominio.Rutina
	actualiza int
}

func (*repositorioEstadoPrueba) ListarRutinas(context.Context, string) ([]dominio.Rutina, error) {
	return nil, nil
}

func (*repositorioEstadoPrueba) CrearRutina(context.Context, *dominio.Rutina) error {
	return nil
}

func (r *repositorioEstadoPrueba) ObtenerRutina(context.Context, uint) (dominio.Rutina, error) {
	return r.rutina, nil
}

func (r *repositorioEstadoPrueba) ActualizarEstado(_ context.Context, _ uint, actual, nuevo string) (int64, error) {
	r.actualiza++
	if r.rutina.Estado != actual {
		return 0, nil
	}
	r.rutina.Estado = nuevo
	return 1, nil
}

func (r *repositorioEstadoPrueba) ObtenerRutinaConCliente(context.Context, uint) (dominio.Rutina, error) {
	return r.rutina, nil
}

func (*repositorioEstadoPrueba) ActualizarRutina(context.Context, uint, dominio.Rutina) (dominio.Rutina, error) {
	return dominio.Rutina{}, nil
}

func (*repositorioEstadoPrueba) EliminarRutina(context.Context, uint) (int64, error) {
	return 0, nil
}

func TestCambiarEstadoAplicaReglaDeNegocio(t *testing.T) {
	repo := &repositorioEstadoPrueba{rutina: dominio.Rutina{Estado: "pendiente"}}
	servicio := NuevoServicioRutinas(repo)

	_, err := servicio.Cambiar(context.Background(), 1, "completada")
	var prohibida *TransicionProhibidaError
	if !errors.As(err, &prohibida) {
		t.Fatalf("se esperaba transición prohibida; llegó %v", err)
	}
	if repo.actualiza != 0 {
		t.Fatal("una transición prohibida no debe llegar al repositorio de escritura")
	}

	rutina, err := servicio.Cambiar(context.Background(), 1, "en_progreso")
	if err != nil {
		t.Fatalf("la transición válida falló: %v", err)
	}
	if rutina.Estado != "en_progreso" || repo.actualiza != 1 {
		t.Fatalf("resultado inesperado: estado=%q escrituras=%d", rutina.Estado, repo.actualiza)
	}
}

func TestCambiarEstadoRechazaEstadoInvalidoAntesDeConsultar(t *testing.T) {
	repo := &repositorioEstadoPrueba{}
	servicio := NuevoServicioRutinas(repo)

	_, err := servicio.Cambiar(context.Background(), 1, "pausada")
	if !errors.Is(err, ErrEstadoInvalido) {
		t.Fatalf("se esperaba ErrEstadoInvalido; llegó %v", err)
	}
	if repo.actualiza != 0 {
		t.Fatal("el estado inválido no debe generar escrituras")
	}
}
