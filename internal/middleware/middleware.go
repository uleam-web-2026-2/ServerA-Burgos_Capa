// Paquete middleware: los dos middleware del curso.
//
// El orden importa: Registro va AFUERA y Recuperacion ADENTRO, para que
// el 500 que produce un pánico quede escrito en el log. Con el orden
// invertido, ese 500 no aparece registrado.
package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/respuesta"
)

// grabador envuelve al ResponseWriter para quedarse con el código de
// estado real que escribió el manejador. Sin esto, el registro tendría
// que suponer que siempre fue 200.
type grabador struct {
	http.ResponseWriter
	codigo int
}

func (g *grabador) WriteHeader(codigo int) {
	g.codigo = codigo
	g.ResponseWriter.WriteHeader(codigo)
}

// Registro escribe una línea por petición con el método, la ruta, el
// código de estado y cuánto tardó.
func Registro(siguiente http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		g := &grabador{ResponseWriter: w, codigo: http.StatusOK}
		siguiente.ServeHTTP(g, r)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, g.codigo, time.Since(inicio))
	})
}

// Recuperacion atrapa un pánico del manejador y responde 500 con la
// envoltura del curso, en vez de dejar caer la conexión.
func Recuperacion(siguiente http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if p := recover(); p != nil {
				log.Printf("pánico atrapado en %s %s: %v", r.Method, r.URL.Path, p)
				respuesta.Error(w, http.StatusInternalServerError,
					"error_interno", "Ocurrió un error inesperado")
			}
		}()
		siguiente.ServeHTTP(w, r)
	})
}
