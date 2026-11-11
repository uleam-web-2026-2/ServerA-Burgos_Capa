package gymcoach

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/aplicacion"
)

func enrutadorDePrueba() http.Handler { return enrutadorCon(rolDePrueba) }

func TestCrearRutinaValidaDatosAntesDeAccederABase(t *testing.T) {
	casos := []struct {
		nombre string
		cuerpo string
		codigo int
	}{
		{"JSON roto", `{"titulo":"Fuerza"`, http.StatusBadRequest},
		{"campo desconocido", `{"titulo":"Fuerza","estado":"completada"}`, http.StatusBadRequest},
		{"título vacío", `{"titulo":" ","descripcion":"Rutina","dias_semana":3,"cliente_id":1}`, http.StatusUnprocessableEntity},
		{"días fuera de rango", `{"titulo":"Fuerza","descripcion":"Rutina","dias_semana":8,"cliente_id":1}`, http.StatusUnprocessableEntity},
		{"cliente ausente", `{"titulo":"Fuerza","descripcion":"Rutina","dias_semana":3}`, http.StatusUnprocessableEntity},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/rutinas", strings.NewReader(caso.cuerpo))
			rec := httptest.NewRecorder()
			enrutadorDePrueba().ServeHTTP(rec, req)
			if rec.Code != caso.codigo {
				t.Fatalf("se esperaba %d y llegó %d con cuerpo %s", caso.codigo, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestActualizarRutinaRechazaEstadoEnElCuerpoDelPut(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/rutinas/1", strings.NewReader(
		`{"titulo":"Fuerza","descripcion":"Rutina","dias_semana":3,"estado":"completada","cliente_id":1}`,
	))
	var entrada aplicacion.RutinaActualizacion
	if err := decodificar(req, &entrada); err == nil {
		t.Fatal("el DTO de actualización no debe aceptar el campo estado")
	}
}

func TestIDNoNumericoResponde400(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rutinas/abc", nil)
	rec := httptest.NewRecorder()
	enrutadorDePrueba().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba %d y llegó %d", http.StatusBadRequest, rec.Code)
	}
}

func TestTransicionesDeRutina(t *testing.T) {
	casos := []struct {
		de, a string
		ok    bool
	}{
		{"pendiente", "en_progreso", true},
		{"en_progreso", "completada", true},
		{"pendiente", "completada", false},
		{"completada", "pendiente", false},
		{"completada", "en_progreso", false},
		{"desconocido", "pendiente", false},
	}
	for _, caso := range casos {
		t.Run(caso.de+"_a_"+caso.a, func(t *testing.T) {
			if resultado := puedeTransicionar(caso.de, caso.a); resultado != caso.ok {
				t.Fatalf("puedeTransicionar(%q, %q) = %v; se esperaba %v", caso.de, caso.a, resultado, caso.ok)
			}
		})
	}
}
