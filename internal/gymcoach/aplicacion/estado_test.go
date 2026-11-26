package aplicacion

import (
	"context"
	"errors"
	"testing"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/dominio"
)

type repositorioEstadoPrueba struct {
	rutina                dominio.Rutina
	actualiza             int
	filasAfectadas        int64
	errorActualizar       error
	forzarConflictoEstado bool
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
	if r.errorActualizar != nil {
		return 0, r.errorActualizar
	}
	if r.forzarConflictoEstado {
		return 0, nil
	}
	if r.rutina.Estado != actual {
		return 0, nil
	}
	r.rutina.Estado = nuevo
	if r.filasAfectadas != 0 {
		return r.filasAfectadas, nil
	}
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

type avisadorEstadoPrueba struct {
	destinatario string
	avisos       []Aviso
}

func (a *avisadorEstadoPrueba) Enviar(quien string, aviso Aviso) {
	a.destinatario = quien
	a.avisos = append(a.avisos, aviso)
}

func TestCambiarEstadoAvisaAlClienteDuenio(t *testing.T) {
	repo := &repositorioEstadoPrueba{rutina: dominio.Rutina{
		ID:      42,
		Estado:  "pendiente",
		Cliente: dominio.Cliente{Nombre: "Ana"},
	}}
	servicio := NuevoServicioRutinas(repo)
	avisador := &avisadorEstadoPrueba{}
	servicio.AvisarCon(avisador)

	_, err := servicio.Cambiar(context.Background(), 42, "completada")
	var prohibida *TransicionProhibidaError
	if !errors.As(err, &prohibida) {
		t.Fatalf("se esperaba transición prohibida; llegó %v", err)
	}
	if repo.actualiza != 0 {
		t.Fatal("una transición prohibida no debe llegar al repositorio de escritura")
	}
	if len(avisador.avisos) != 0 {
		t.Fatal("una transición prohibida no debe enviar avisos")
	}

	rutina, err := servicio.Cambiar(context.Background(), 42, "en_progreso")
	if err != nil {
		t.Fatalf("la transición válida falló: %v", err)
	}
	if rutina.Estado != "en_progreso" || repo.actualiza != 1 {
		t.Fatalf("resultado inesperado: estado=%q escrituras=%d", rutina.Estado, repo.actualiza)
	}
	if avisador.destinatario != "Ana" || len(avisador.avisos) != 1 {
		t.Fatalf("destinatario/avisos inesperados: %q, %d", avisador.destinatario, len(avisador.avisos))
	}
	aviso := avisador.avisos[0]
	if aviso.Tipo != "estado_cambiado" || aviso.ID != 42 || aviso.Estado != "en_progreso" {
		t.Fatalf("contenido de aviso inesperado: %+v", aviso)
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

func TestCambiarEstadoNoAvisaSiNoSeActualizo(t *testing.T) {
	repo := &repositorioEstadoPrueba{
		rutina: dominio.Rutina{
			Estado:  "pendiente",
			Cliente: dominio.Cliente{Nombre: "Ana"},
		},
		forzarConflictoEstado: true,
	}
	servicio := NuevoServicioRutinas(repo)
	avisador := &avisadorEstadoPrueba{}
	servicio.AvisarCon(avisador)

	_, err := servicio.Cambiar(context.Background(), 1, "en_progreso")
	if !errors.Is(err, dominio.ErrEstadoModificado) {
		t.Fatalf("se esperaba ErrEstadoModificado; llegó %v", err)
	}
	if len(avisador.avisos) != 0 {
		t.Fatal("un estado que no se actualizó no debe enviar avisos")
	}
}

func TestCambiarEstadoNoAvisaSiFallaLaEscritura(t *testing.T) {
	repo := &repositorioEstadoPrueba{
		rutina: dominio.Rutina{
			Estado:  "pendiente",
			Cliente: dominio.Cliente{Nombre: "Ana"},
		},
		errorActualizar: errors.New("fallo de base de datos"),
	}
	servicio := NuevoServicioRutinas(repo)
	avisador := &avisadorEstadoPrueba{}
	servicio.AvisarCon(avisador)

	_, err := servicio.Cambiar(context.Background(), 1, "en_progreso")
	if err == nil {
		t.Fatal("se esperaba el error de escritura")
	}
	if len(avisador.avisos) != 0 {
		t.Fatal("un fallo de escritura no debe enviar avisos")
	}
}
