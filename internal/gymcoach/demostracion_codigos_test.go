package gymcoach

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/sesion"
)

func TestDemostracionCodigosHTTP(t *testing.T) {
	casos := []struct {
		nombre      string
		metodo      string
		ruta        string
		cuerpo      string
		rol         string
		codigoEsper int
	}{
		// 401: Sin token
		{"Sin token a ruta protegida", "GET", "/clientes", "", "", http.StatusUnauthorized},
		// 403: Rol cliente intentando acceder a ruta de entrenador
		{"Cliente intentando ver /clientes (Solo Entrenador)", "GET", "/clientes", "", sesion.RolCliente, http.StatusForbidden},
		// 403: Entrenador intentando cambiar estado (Solo Cliente)
		{"Entrenador intentando hacer PATCH /rutinas/1 (Solo Cliente)", "PATCH", "/rutinas/1", `{"estado":"en_progreso"}`, sesion.RolEntrenador, http.StatusForbidden},
		// 400: ID no numérico
		{"Ruta con ID no numerico /rutinas/abc", "GET", "/rutinas/abc", "", sesion.RolEntrenador, http.StatusBadRequest},
		// 400: JSON roto o mal formado
		{"JSON roto en POST /rutinas", "POST", "/rutinas", `{"titulo":"Fuerza"`, sesion.RolEntrenador, http.StatusBadRequest},
		// 400: Campo desconocido en el cuerpo
		{"Campo desconocido en POST /rutinas", "POST", "/rutinas", `{"titulo":"Fuerza","estado":"completada"}`, sesion.RolEntrenador, http.StatusBadRequest},
		// 422: Filtro de estado no válido
		{"Filtro de estado inexistente ?estado=inventado", "GET", "/rutinas?estado=inventado", "", sesion.RolEntrenador, http.StatusUnprocessableEntity},
		// 422: Datos requeridos faltantes
		{"Campos obligatorios ausentes en POST /rutinas", "POST", "/rutinas", `{}`, sesion.RolEntrenador, http.StatusUnprocessableEntity},
		// 422: Días de semana fuera de rango (1 a 7)
		{"Dias fuera de rango (dias_semana: 9)", "POST", "/rutinas", `{"titulo":"Fuerza","descripcion":"Rutina","dias_semana":9,"cliente_id":1}`, sesion.RolEntrenador, http.StatusUnprocessableEntity},
		// 422: Estado no permitido en PATCH
		{"Estado inventado en PATCH /rutinas/1", "PATCH", "/rutinas/1", `{"estado":"saltando"}`, sesion.RolCliente, http.StatusUnprocessableEntity},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			servidor := enrutadorCon(c.rol)
			req := httptest.NewRequest(c.metodo, c.ruta, strings.NewReader(c.cuerpo))
			rec := httptest.NewRecorder()

			servidor.ServeHTTP(rec, req)

			fmt.Printf("\n======================================================\n[PRUEBA] %s\n-> Peticion: %s %s (Rol: %q)\n-> Codigo HTTP obtenido: %d %s\n-> Respuesta JSON: %s======================================================\n",
				c.nombre, c.metodo, c.ruta, c.rol, rec.Code, http.StatusText(rec.Code), rec.Body.String())

			if rec.Code != c.codigoEsper {
				t.Fatalf("se esperaba codigo %d y llego %d", c.codigoEsper, rec.Code)
			}
		})
	}
}
