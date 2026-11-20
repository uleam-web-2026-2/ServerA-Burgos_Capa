package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func TestCORSOrigenPermitidoYPreflight(t *testing.T) {
	r := chi.NewRouter()
	r.Use(configurarCORS())
	r.Post("/login", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("origen permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		if res.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
			t.Fatalf("falta la cabecera CORS permitida: %q", res.Header().Get("Access-Control-Allow-Origin"))
		}
	})

	t.Run("origen no permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.Header.Set("Origin", "http://malicioso.example")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		if res.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Fatalf("un origen no permitido recibió CORS: %q", res.Header().Get("Access-Control-Allow-Origin"))
		}
	})

	t.Run("preflight permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/login", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		req.Header.Set("Access-Control-Request-Method", http.MethodPatch)
		req.Header.Set("Access-Control-Request-Headers", "Authorization, Content-Type")
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		if res.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
			t.Fatalf("el preflight permitido perdió su origen: %q", res.Header().Get("Access-Control-Allow-Origin"))
		}
		if !strings.Contains(res.Header().Get("Access-Control-Allow-Methods"), http.MethodPatch) {
			t.Fatalf("el preflight no permitió PATCH: %q", res.Header().Get("Access-Control-Allow-Methods"))
		}
		allowedHeaders := strings.ToLower(res.Header().Get("Access-Control-Allow-Headers"))
		if !strings.Contains(allowedHeaders, "authorization") || !strings.Contains(allowedHeaders, "content-type") {
			t.Fatalf("el preflight no permitió los headers necesarios: %q", res.Header().Get("Access-Control-Allow-Headers"))
		}
	})
}

func TestLimitadorLoginBloqueaSextoIntento(t *testing.T) {
	r := chi.NewRouter()
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.With(limitadorLogin()).Post("/login", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	for intento := 1; intento <= 6; intento++ {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "192.0.2.10:4321"
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		if intento <= 5 && res.Code != http.StatusUnauthorized {
			t.Fatalf("intento %d: esperado 401, obtenido %d", intento, res.Code)
		}
		if intento == 6 && res.Code != http.StatusTooManyRequests {
			t.Fatalf("intento 6: esperado 429, obtenido %d", res.Code)
		}
		if intento == 6 {
			retryAfter, err := strconv.Atoi(res.Header().Get("Retry-After"))
			if err != nil || retryAfter <= 0 {
				t.Fatalf("el 429 debe indicar Retry-After positivo; obtenido %q", res.Header().Get("Retry-After"))
			}
			if res.Header().Get("X-RateLimit-Limit") != "5" || res.Header().Get("X-RateLimit-Remaining") != "0" {
				t.Fatalf("cabeceras inesperadas en 429: limit=%q remaining=%q", res.Header().Get("X-RateLimit-Limit"), res.Header().Get("X-RateLimit-Remaining"))
			}
		}
	}
}
