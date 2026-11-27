package gymcoach

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/aplicacion"
	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/respuesta"
	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/sesion"
)

const bufferAvisos = 8

type suscripcionAviso struct {
	rol string
	ch  chan aplicacion.Aviso
}

// Avisos mantiene las conexiones WebSocket activas por usuario y reparte
// avisos en tiempo real sin bloquear cuando una conexión es lenta.
type Avisos struct {
	tokens *sesion.Tokens
	origen string

	mu      sync.Mutex
	canales map[string][]*suscripcionAviso
}

func NuevosAvisos(tokens *sesion.Tokens, origen string) *Avisos {
	return &Avisos{
		tokens:  tokens,
		origen:  origen,
		canales: make(map[string][]*suscripcionAviso),
	}
}

func (a *Avisos) suscribir(nombre, rol string) *suscripcionAviso {
	sub := &suscripcionAviso{
		rol: rol,
		ch:  make(chan aplicacion.Aviso, bufferAvisos),
	}
	a.mu.Lock()
	a.canales[nombre] = append(a.canales[nombre], sub)
	a.mu.Unlock()
	return sub
}

func (a *Avisos) desuscribir(nombre string, objetivo *suscripcionAviso) {
	a.mu.Lock()
	defer a.mu.Unlock()

	actuales := a.canales[nombre]
	for i, sub := range actuales {
		if sub != objetivo {
			continue
		}
		actuales[i] = actuales[len(actuales)-1]
		actuales[len(actuales)-1] = nil
		restantes := actuales[:len(actuales)-1]
		if len(restantes) == 0 {
			delete(a.canales, nombre)
		} else {
			a.canales[nombre] = restantes
		}
		close(objetivo.ch)
		return
	}
}

func enviarSinBloquear(sub *suscripcionAviso, aviso aplicacion.Aviso) {
	select {
	case sub.ch <- aviso:
	default:
	}
}

// Enviar entrega un aviso únicamente a las conexiones activas del usuario indicado.
func (a *Avisos) Enviar(quien string, aviso aplicacion.Aviso) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, sub := range a.canales[quien] {
		enviarSinBloquear(sub, aviso)
	}
}

// AvisarCreacion notifica a todas las conexiones activas del rol entrenador,
// excluyendo al usuario que originó la creación.
func (a *Avisos) AvisarCreacion(excluir string, aviso aplicacion.Aviso) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for nombre, subs := range a.canales {
		if nombre == excluir {
			continue
		}
		for _, sub := range subs {
			if sub.rol == sesion.RolEntrenador {
				enviarSinBloquear(sub, aviso)
			}
		}
	}
}

// Conectar autentica la sesión mediante el token de la query antes de aceptar
// el WebSocket y restringe el origen al frontend autorizado.
func (a *Avisos) Conectar(w http.ResponseWriter, r *http.Request) {
	s, err := a.tokens.Leer(r.URL.Query().Get("token"))
	if err != nil {
		if errors.Is(err, sesion.ErrVencido) {
			respuesta.Error(w, http.StatusUnauthorized, "token_vencido",
				"el token está vencido: inicie sesión otra vez")
			return
		}
		respuesta.Error(w, http.StatusUnauthorized, "sin_sesion",
			"falta el token o no es válido")
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{a.origen},
	})
	if err != nil {
		return
	}
	defer conn.CloseNow()

	sub := a.suscribir(s.Nombre, s.Rol)
	defer a.desuscribir(s.Nombre, sub)

	// Se usa context.Background() porque middleware.Timeout cancela r.Context().
	ctx := conn.CloseRead(context.Background())

	for {
		select {
		case <-ctx.Done():
			return
		case aviso, ok := <-sub.ch:
			if !ok {
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := wsjson.Write(writeCtx, conn, aviso)
			cancel()
			if err != nil {
				return
			}
		}
	}
}
