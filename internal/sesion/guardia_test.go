package sesion

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTokensEmitirYLeer(t *testing.T) {
	tokens := NuevosTokens("secreto-de-prueba-muy-seguro-123456")
	usuario := Usuario{ID: 10, Nombre: "carlos", Rol: RolEntrenador}

	tokenStr, err := tokens.Emitir(usuario)
	if err != nil {
		t.Fatalf("error al emitir token: %v", err)
	}

	sesion, err := tokens.Leer(tokenStr)
	if err != nil {
		t.Fatalf("error al leer token válido: %v", err)
	}

	if sesion.UsuarioID != usuario.ID || sesion.Nombre != usuario.Nombre || sesion.Rol != usuario.Rol {
		t.Fatalf("los datos de sesión no coinciden: obtenido %+v, esperado %+v", sesion, usuario)
	}
}

func TestTokensConFirmaInvalida(t *testing.T) {
	tokensA := NuevosTokens("secreto-A-12345678901234567890")
	tokensB := NuevosTokens("secreto-B-12345678901234567890")

	tokenStr, err := tokensA.Emitir(Usuario{ID: 1, Nombre: "ana", Rol: RolCliente})
	if err != nil {
		t.Fatalf("error al emitir token: %v", err)
	}

	_, err = tokensB.Leer(tokenStr)
	if err == nil {
		t.Fatal("se esperaba error al validar token firmado con otra clave")
	}
}

func TestGuardiaAutenticar(t *testing.T) {
	tokens := NuevosTokens("secreto-guardia-1234567890123456")
	handlerFinal := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, ok := De(r)
		if !ok {
			t.Error("se esperaba sesión en contexto pero no se encontró")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(s.Nombre))
	})

	middleware := tokens.Autenticar(handlerFinal)

	casos := []struct {
		nombre     string
		authHeader string
		codigo     int
	}{
		{"sin cabecera", "", http.StatusUnauthorized},
		{"sin prefijo Bearer", "Token abcdef", http.StatusUnauthorized},
		{"token basura", "Bearer token-invalido-o-roto", http.StatusUnauthorized},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/ruta-protegida", nil)
			if c.authHeader != "" {
				req.Header.Set("Authorization", c.authHeader)
			}
			rec := httptest.NewRecorder()
			middleware.ServeHTTP(rec, req)

			if rec.Code != c.codigo {
				t.Fatalf("se esperaba código %d y llegó %d con cuerpo %s", c.codigo, rec.Code, rec.Body.String())
			}
		})
	}

	// Caso exitoso con token real
	tokenValido, err := tokens.Emitir(Usuario{ID: 1, Nombre: "ana", Rol: RolCliente})
	if err != nil {
		t.Fatalf("error emitiendo token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/ruta-protegida", nil)
	req.Header.Set("Authorization", "Bearer "+tokenValido)
	rec := httptest.NewRecorder()
	middleware.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("se esperaba 200 con token válido y llegó %d", rec.Code)
	}
	if rec.Body.String() != "ana" {
		t.Fatalf("se esperaba cuerpo 'ana' y llegó '%s'", rec.Body.String())
	}
}

func TestGuardiaSoloRol(t *testing.T) {
	handlerFinal := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	guardiaEntrenador := SoloRol(RolEntrenador)(handlerFinal)

	// Petición con rol entrenador -> 200 OK
	reqValido := httptest.NewRequest(http.MethodGet, "/solo-entrenador", nil)
	ctxValido := context.WithValue(reqValido.Context(), clave{}, Sesion{UsuarioID: 1, Nombre: "luis", Rol: RolEntrenador})
	recValido := httptest.NewRecorder()
	guardiaEntrenador.ServeHTTP(recValido, reqValido.WithContext(ctxValido))

	if recValido.Code != http.StatusOK {
		t.Fatalf("se esperaba 200 para entrenador y llegó %d", recValido.Code)
	}

	// Petición con rol cliente -> 403 Forbidden
	reqCliente := httptest.NewRequest(http.MethodGet, "/solo-entrenador", nil)
	ctxCliente := context.WithValue(reqCliente.Context(), clave{}, Sesion{UsuarioID: 2, Nombre: "ana", Rol: RolCliente})
	recCliente := httptest.NewRecorder()
	guardiaEntrenador.ServeHTTP(recCliente, reqCliente.WithContext(ctxCliente))

	if recCliente.Code != http.StatusForbidden {
		t.Fatalf("se esperaba 403 para cliente en ruta de entrenador y llegó %d", recCliente.Code)
	}

	// Petición sin contexto de sesión -> 403 Forbidden
	reqSinSesion := httptest.NewRequest(http.MethodGet, "/solo-entrenador", nil)
	recSinSesion := httptest.NewRecorder()
	guardiaEntrenador.ServeHTTP(recSinSesion, reqSinSesion)

	if recSinSesion.Code != http.StatusForbidden {
		t.Fatalf("se esperaba 403 sin sesión y llegó %d", recSinSesion.Code)
	}
}

func TestLoginJSONInvalidoResponde400(t *testing.T) {
	manejador := &Manejador{DB: nil, Tokens: NuevosTokens("clave")}
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{usuario: "invalido"`))
	rec := httptest.NewRecorder()

	manejador.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba 400 para JSON roto en login y llegó %d", rec.Code)
	}
}
