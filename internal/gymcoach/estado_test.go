package gymcoach

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/sesion"
)

// El cliente cambia el estado: lo que se valida ANTES de tocar la base.
func TestCambiarEstadoValidaAntesDeAccederABase(t *testing.T) {
	casos := []struct {
		nombre, ruta, cuerpo string
		codigo               int
	}{
		{"id no numérico", "/rutinas/abc", `{"estado":"en_progreso"}`, http.StatusBadRequest},
		{"JSON roto", "/rutinas/1", `{"estado":`, http.StatusBadRequest},
		{"campo desconocido", "/rutinas/1", `{"estado":"en_progreso","otro":1}`, http.StatusBadRequest},
		{"estado inventado", "/rutinas/1", `{"estado":"inventado"}`, http.StatusUnprocessableEntity},
		{"estado vacío", "/rutinas/1", `{}`, http.StatusUnprocessableEntity},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPatch, c.ruta, strings.NewReader(c.cuerpo))
			rec := httptest.NewRecorder()
			enrutadorCon(sesion.RolCliente).ServeHTTP(rec, req)
			if rec.Code != c.codigo {
				t.Fatalf("se esperaba %d y llegó %d con cuerpo %s", c.codigo, rec.Code, rec.Body.String())
			}
		})
	}
}

// El entrenador no cambia el estado: el guardia de rol responde 403.
func TestEntrenadorNoCambiaElEstado(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/rutinas/1", strings.NewReader(`{"estado":"en_progreso"}`))
	rec := httptest.NewRecorder()
	enrutadorCon(sesion.RolEntrenador).ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("se esperaba 403 y llegó %d con cuerpo %s", rec.Code, rec.Body.String())
	}
}

// El cliente no entra a las rutas del entrenador.
func TestClienteNoUsaLasRutasDelEntrenador(t *testing.T) {
	rutas := []struct{ metodo, ruta string }{
		{"GET", "/clientes"},
		{"POST", "/clientes"},
		{"GET", "/rutinas"},
		{"POST", "/rutinas"},
		{"PUT", "/rutinas/1"},
		{"DELETE", "/rutinas/1"},
	}
	for _, p := range rutas {
		t.Run(p.metodo+" "+p.ruta, func(t *testing.T) {
			req := httptest.NewRequest(p.metodo, p.ruta, strings.NewReader(`{}`))
			rec := httptest.NewRecorder()
			enrutadorCon(sesion.RolCliente).ServeHTTP(rec, req)
			if rec.Code != http.StatusForbidden {
				t.Fatalf("se esperaba 403 y llegó %d", rec.Code)
			}
		})
	}
}

func TestListarRutinasRechazaFiltroInvalidoAntesDeAccederABase(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rutinas?estado=terminada", nil)
	rec := httptest.NewRecorder()
	enrutadorCon(RolEntrenador).ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("se esperaba %d y llegó %d con cuerpo %s",
			http.StatusUnprocessableEntity, rec.Code, rec.Body.String())
	}
}
