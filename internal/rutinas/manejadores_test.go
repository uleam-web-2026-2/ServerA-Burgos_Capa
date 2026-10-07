package rutinas

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

// enrutador arma el servidor de prueba SIN base de datos (DB: nil).
// Sirve porque todas las pruebas de este archivo responden ANTES
// de la línea "desde aquí hace falta la base de datos".
func enrutador() http.Handler {
	r := chi.NewRouter()
	(&Manejador{DB: nil}).Rutas(r) // el mismo Rutas de main.go
	return r
}

// TestCrear: una función, varias filas. Cada fila es un caso.
func TestCrear(t *testing.T) {
	casos := []struct {
		nombre string
		cuerpo string
		codigo int
	}{
		{
			nombre: "JSON roto responde 400",
			cuerpo: `{"asunto": "Proyector"`,
			codigo: http.StatusBadRequest,
		},
		{
			nombre: "estado inventado responde 422",
			cuerpo: `{"asunto": "Proyector", "estado": "urgente"}`,
			codigo: http.StatusUnprocessableEntity,
		},
		{
			nombre: "asunto vacío responde 422",
			cuerpo: `{"asunto": "", "estado": "pendiente"}`,
			codigo: http.StatusUnprocessableEntity,
		},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			peticion := httptest.NewRequest(http.MethodPost, "/tickets", strings.NewReader(caso.cuerpo))
			grabadora := httptest.NewRecorder()
			enrutador().ServeHTTP(grabadora, peticion)
			if grabadora.Code != caso.codigo {
				t.Fatalf("se esperaba %d y llegó %d con cuerpo %s",
					caso.codigo, grabadora.Code, grabadora.Body.String())
			}
		})
	}
}

// El id de la ruta se valida antes de consultar.
func TestVerUnoConIDQueNoEsNumeroResponde400(t *testing.T) {
	peticion := httptest.NewRequest(http.MethodGet, "/tickets/abc", nil)
	grabadora := httptest.NewRecorder()
	enrutador().ServeHTTP(grabadora, peticion)
	if grabadora.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba %d y llegó %d", http.StatusBadRequest, grabadora.Code)
	}
}
