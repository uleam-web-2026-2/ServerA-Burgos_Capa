package gymcoach

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/aplicacion"
	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/gymcoach/dominio"
	"github.com/uleam-web-2026-2/ServerA-Burgos_Capa/internal/sesion"
)

const origenPermitidoPrueba = "http://localhost:5173"

func emitirTokenPrueba(t *testing.T, id uint, nombre, rol string) string {
	t.Helper()
	token, err := tokensDePrueba.Emitir(sesion.Usuario{ID: id, Nombre: nombre, Rol: rol})
	if err != nil {
		t.Fatalf("no se pudo emitir el token de prueba: %v", err)
	}
	return token
}

func servidorWebSocketPrueba(t *testing.T, timeout time.Duration) (*Avisos, *httptest.Server) {
	t.Helper()
	avisos := NuevosAvisos(tokensDePrueba, origenPermitidoPrueba)
	r := chi.NewRouter()
	r.Group(func(httpRutas chi.Router) {
		if timeout > 0 {
			httpRutas.Use(middleware.Timeout(timeout))
		}
		(&Manejador{DB: nil}).Rutas(httpRutas, tokensDePrueba)
	})
	r.Get("/ws", avisos.Conectar)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)
	return avisos, srv
}

func conectarWSPrueba(t *testing.T, baseURL, token, origen string) *websocket.Conn {
	t.Helper()
	wsURL := "ws" + strings.TrimPrefix(baseURL, "http") + "/ws?token=" + token
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
		HTTPHeader: http.Header{"Origin": []string{origen}},
	})
	if err != nil {
		t.Fatalf("no se pudo abrir WebSocket para %q: %v", origen, err)
	}
	t.Cleanup(func() { _ = conn.CloseNow() })
	return conn
}

func leerAvisosAsync(conn *websocket.Conn) <-chan aplicacion.Aviso {
	ch := make(chan aplicacion.Aviso, 8)
	go func() {
		defer close(ch)
		for {
			var aviso aplicacion.Aviso
			if err := wsjson.Read(context.Background(), conn, &aviso); err != nil {
				return
			}
			ch <- aviso
		}
	}()
	return ch
}

func esperarAvisoCanal(t *testing.T, ch <-chan aplicacion.Aviso, plazo time.Duration) aplicacion.Aviso {
	t.Helper()
	select {
	case aviso, ok := <-ch:
		if !ok {
			t.Fatal("la conexión WebSocket se cerró antes de recibir el aviso")
		}
		return aviso
	case <-time.After(plazo):
		t.Fatalf("no llegó ningún aviso en %s", plazo)
		return aplicacion.Aviso{}
	}
}

func confirmarSinAvisoCanal(t *testing.T, ch <-chan aplicacion.Aviso, espera time.Duration) {
	t.Helper()
	select {
	case aviso, ok := <-ch:
		if ok {
			t.Fatalf("no debía llegar aviso, pero llegó: %+v", aviso)
		}
	case <-time.After(espera):
	}
}

func esperarAviso(t *testing.T, conn *websocket.Conn, plazo time.Duration) aplicacion.Aviso {
	t.Helper()
	return esperarAvisoCanal(t, leerAvisosAsync(conn), plazo)
}

func confirmarSinAviso(t *testing.T, conn *websocket.Conn, espera time.Duration) {
	t.Helper()
	confirmarSinAvisoCanal(t, leerAvisosAsync(conn), espera)
}

func TestWebSocketSinTokenResponde401(t *testing.T) {
	_, srv := servidorWebSocketPrueba(t, 0)

	for _, url := range []string{srv.URL + "/ws", srv.URL + "/ws?token=no-valido"} {
		res, err := http.Get(url)
		if err != nil {
			t.Fatalf("falló la petición a %s: %v", url, err)
		}
		var cuerpo struct {
			OK    bool `json:"ok"`
			Error struct {
				Codigo string `json:"codigo"`
			} `json:"error"`
		}
		_ = json.NewDecoder(res.Body).Decode(&cuerpo)
		_ = res.Body.Close()

		if res.StatusCode != http.StatusUnauthorized {
			t.Fatalf("%s: esperado 401, obtenido %d", url, res.StatusCode)
		}
		if cuerpo.OK || cuerpo.Error.Codigo != "sin_sesion" {
			t.Fatalf("%s: se esperaba código sin_sesion fuera del guardia; llegó %+v", url, cuerpo)
		}
	}
}

func TestWebSocketOtroOrigenResponde403(t *testing.T) {
	_, srv := servidorWebSocketPrueba(t, 0)
	token := emitirTokenPrueba(t, 3, "Daniela Vera", sesion.RolCliente)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws?token=" + token
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, res, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
		HTTPHeader: http.Header{"Origin": []string{"http://malo.example"}},
	})
	if err == nil {
		_ = conn.CloseNow()
		t.Fatal("se esperaba rechazo 403 para origen no autorizado")
	}
	if res == nil || res.StatusCode != http.StatusForbidden {
		t.Fatalf("esperado 403, obtenido %v (err=%v)", res, err)
	}
}

func TestWebSocketAvisoLlegaSoloAlDuenio(t *testing.T) {
	avisos, srv := servidorWebSocketPrueba(t, 0)
	tokenDaniela := emitirTokenPrueba(t, 3, "Daniela Vera", sesion.RolCliente)
	tokenLuisMera := emitirTokenPrueba(t, 4, "Luis Mera", sesion.RolCliente)

	connDaniela := conectarWSPrueba(t, srv.URL, tokenDaniela, origenPermitidoPrueba)
	connLuisMera := conectarWSPrueba(t, srv.URL, tokenLuisMera, origenPermitidoPrueba)

	avisos.Enviar("Daniela Vera", aplicacion.Aviso{
		Tipo:   "estado_cambiado",
		ID:     15,
		Estado: "en_progreso",
	})

	aviso := esperarAviso(t, connDaniela, time.Second)
	if aviso.Tipo != "estado_cambiado" || aviso.ID != 15 || aviso.Estado != "en_progreso" {
		t.Fatalf("aviso inesperado para Daniela Vera: %+v", aviso)
	}
	confirmarSinAviso(t, connLuisMera, 300*time.Millisecond)
}

func TestWebSocketVariasConexionesDesconexionYTimeout(t *testing.T) {
	avisos, srv := servidorWebSocketPrueba(t, 60*time.Millisecond)
	tokenDaniela := emitirTokenPrueba(t, 3, "Daniela Vera", sesion.RolCliente)

	conn1 := conectarWSPrueba(t, srv.URL, tokenDaniela, origenPermitidoPrueba)
	conn2 := conectarWSPrueba(t, srv.URL, tokenDaniela, origenPermitidoPrueba)
	ch1 := leerAvisosAsync(conn1)
	ch2 := leerAvisosAsync(conn2)

	// Supera el middleware.Timeout para comprobar que CloseRead(context.Background()) mantiene el socket abierto.
	time.Sleep(120 * time.Millisecond)

	avisos.Enviar("Daniela Vera", aplicacion.Aviso{Tipo: "estado_cambiado", ID: 1, Estado: "en_progreso"})
	if got := esperarAvisoCanal(t, ch1, time.Second); got.ID != 1 {
		t.Fatalf("conn1 recibió aviso inesperado: %+v", got)
	}
	if got := esperarAvisoCanal(t, ch2, time.Second); got.ID != 1 {
		t.Fatalf("conn2 recibió aviso inesperado: %+v", got)
	}

	_ = conn1.Close(websocket.StatusNormalClosure, "cierre")
	time.Sleep(40 * time.Millisecond)

	avisos.Enviar("Daniela Vera", aplicacion.Aviso{Tipo: "estado_cambiado", ID: 2, Estado: "completada"})
	if got := esperarAvisoCanal(t, ch2, time.Second); got.ID != 2 {
		t.Fatalf("conn2 tras desconexión de conn1 recibió aviso inesperado: %+v", got)
	}
}

func TestWebSocketConcurrencia(t *testing.T) {
	avisos, srv := servidorWebSocketPrueba(t, 0)
	tokenDaniela := emitirTokenPrueba(t, 3, "Daniela Vera", sesion.RolCliente)
	tokenLuisMera := emitirTokenPrueba(t, 4, "Luis Mera", sesion.RolCliente)

	connDaniela := conectarWSPrueba(t, srv.URL, tokenDaniela, origenPermitidoPrueba)
	connLuisMera := conectarWSPrueba(t, srv.URL, tokenLuisMera, origenPermitidoPrueba)
	chDaniela := leerAvisosAsync(connDaniela)
	chLuisMera := leerAvisosAsync(connLuisMera)

	const total = 6
	var wg sync.WaitGroup
	for i := 1; i <= total; i++ {
		wg.Add(1)
		go func(id uint) {
			defer wg.Done()
			avisos.Enviar("Daniela Vera", aplicacion.Aviso{Tipo: "estado_cambiado", ID: id, Estado: "en_progreso"})
		}(uint(i))
	}
	wg.Wait()

	recibidos := make(map[uint]bool, total)
	for i := 0; i < total; i++ {
		av := esperarAvisoCanal(t, chDaniela, time.Second)
		recibidos[av.ID] = true
	}
	if len(recibidos) != total {
		t.Fatalf("se esperaban %d avisos distintos; llegaron %d", total, len(recibidos))
	}
	confirmarSinAvisoCanal(t, chLuisMera, 200*time.Millisecond)
}

type repoRetoPrueba struct {
	rutina dominio.Rutina
}

func (*repoRetoPrueba) ListarRutinas(context.Context, string) ([]dominio.Rutina, error) {
	return nil, nil
}

func (r *repoRetoPrueba) CrearRutina(_ context.Context, rutina *dominio.Rutina) error {
	rutina.ID = 91
	r.rutina = *rutina
	r.rutina.Cliente = dominio.Cliente{Nombre: "Daniela Vera"}
	return nil
}

func (r *repoRetoPrueba) ObtenerRutina(context.Context, uint) (dominio.Rutina, error) {
	return r.rutina, nil
}

func (r *repoRetoPrueba) ObtenerRutinaConCliente(context.Context, uint) (dominio.Rutina, error) {
	return r.rutina, nil
}

func (*repoRetoPrueba) ActualizarRutina(context.Context, uint, dominio.Rutina) (dominio.Rutina, error) {
	return dominio.Rutina{}, nil
}

func (*repoRetoPrueba) EliminarRutina(context.Context, uint) (int64, error) {
	return 0, nil
}

func (r *repoRetoPrueba) ActualizarEstado(_ context.Context, _ uint, actual, nuevo string) (int64, error) {
	if r.rutina.Estado != actual {
		return 0, nil
	}
	r.rutina.Estado = nuevo
	return 1, nil
}

func TestRetoAvisaAEntrenadoresAlCrearSinAvisarAClienteNiCreador(t *testing.T) {
	avisos, srv := servidorWebSocketPrueba(t, 0)
	servicio := aplicacion.NuevoServicioRutinas(&repoRetoPrueba{})
	servicio.AvisarCon(avisos)
	servicio.AvisarCreacionCon(avisos)

	tokenEntrenador1 := emitirTokenPrueba(t, 2, "luis", sesion.RolEntrenador)
	tokenEntrenador2 := emitirTokenPrueba(t, 5, "carla", sesion.RolEntrenador)
	tokenEntrenador3 := emitirTokenPrueba(t, 6, "marco", sesion.RolEntrenador)
	tokenCliente := emitirTokenPrueba(t, 3, "Daniela Vera", sesion.RolCliente)

	connLuis := conectarWSPrueba(t, srv.URL, tokenEntrenador1, origenPermitidoPrueba)
	connCarla := conectarWSPrueba(t, srv.URL, tokenEntrenador2, origenPermitidoPrueba)
	connMarco := conectarWSPrueba(t, srv.URL, tokenEntrenador3, origenPermitidoPrueba)
	connDaniela := conectarWSPrueba(t, srv.URL, tokenCliente, origenPermitidoPrueba)

	chLuis := leerAvisosAsync(connLuis)
	chCarla := leerAvisosAsync(connCarla)
	chMarco := leerAvisosAsync(connMarco)
	chDaniela := leerAvisosAsync(connDaniela)

	_, err := servicio.Crear(context.Background(), aplicacion.RutinaEntrada{
		Titulo:      "Fuerza nueva",
		Descripcion: "Plan semanal",
		DiasSemana:  3,
		ClienteID:   1,
	}, "luis")
	if err != nil {
		t.Fatalf("no se esperaba error al crear rutina: %v", err)
	}

	// Dos entrenadores distintos al creador reciben rutina_nueva; ni el creador ni el cliente la reciben.
	for nombre, ch := range map[string]<-chan aplicacion.Aviso{"carla": chCarla, "marco": chMarco} {
		aviso := esperarAvisoCanal(t, ch, time.Second)
		if aviso.Tipo != "rutina_nueva" || aviso.ID != 91 || aviso.Estado != "pendiente" {
			t.Fatalf("aviso inesperado para entrenador %s: %+v", nombre, aviso)
		}
	}
	confirmarSinAvisoCanal(t, chLuis, 250*time.Millisecond)
	confirmarSinAvisoCanal(t, chDaniela, 250*time.Millisecond)

	// El aviso del punto 1 sigue llegando solo al cliente dueño al cambiar estado.
	if _, err := servicio.Cambiar(context.Background(), 91, "en_progreso"); err != nil {
		t.Fatalf("error al cambiar estado: %v", err)
	}
	avisoDuenio := esperarAvisoCanal(t, chDaniela, time.Second)
	if avisoDuenio.Tipo != "estado_cambiado" || avisoDuenio.ID != 91 || avisoDuenio.Estado != "en_progreso" {
		t.Fatalf("aviso de estado inesperado para Daniela Vera: %+v", avisoDuenio)
	}
	confirmarSinAvisoCanal(t, chCarla, 250*time.Millisecond)
}
