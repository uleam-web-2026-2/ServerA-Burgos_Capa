# Taller 13 · reparto y seguimiento

Este reparto resume [el enunciado oficial](./Semana_13_-_Taller.docx) y lo
adapta a GymCoach. La entrega se prepara en la rama `taller13`, creada desde
`origin/main` después de integrar el Taller 12.

## Decisión de adaptación a GymCoach

Se conservan las reglas existentes del Hito 2: el entrenador crea y asigna una
rutina mediante `cliente_id`, y solo el cliente puede cambiar su estado con
`PATCH /rutinas/{id}`. El enunciado usa una reserva y permite que quien
administra cambie su estado; copiar esa regla aquí rompería el 403 que ya se
prueba para el entrenador.

La rutina ya identifica a su cliente dueño con `Rutina.ClienteID` y la relación
`Rutina.Cliente`. El servicio de aplicación usará el nombre de ese cliente como
destinatario lógico del aviso; las cuentas de prueba de clientes deben usar el
mismo nombre que el perfil. No se agrega otro campo de dueño ni se cambia el
contrato de creación de rutinas. El adaptador de transporte debe usar la misma
clave al registrar la conexión autenticada.

## Persona 1 · preparación y punto 1

### Fase 0 · preparación compartida

- [x] Crear `taller13` desde `origin/main`, con los cambios integrados del
  Taller 12.
- [x] Comprobar `go test ./...` antes de cambiar el código.
- [x] Agregar `github.com/coder/websocket` para que Persona 2 pueda implementar
  el transporte.
- [x] `go build ./...` y `go test ./...` pasan después de los cambios de Persona 1.
- [ ] Persona 2 completa `peticiones.md` con las seis respuestas reales de
  línea base; el intento local que no respondió está en
  [`pendientes.md`](./pendientes.md).

### Punto 1 · servicio y dueño de la rutina

- [x] Confirmar que `ClienteID` ya asocia cada rutina con el cliente dueño.
- [x] Definir el aviso de cambio de estado y la interfaz `Avisador` en la capa
  de aplicación, sin introducir dependencias de WebSocket.
- [x] Emitir el aviso solo después de que la transición válida quede guardada;
  una transición prohibida o un fallo de escritura no debe avisar.
- [x] Mantener las reglas de permisos y los códigos HTTP del Hito 2.
- [x] Agregar dos cuentas de prueba con rol cliente cuyos nombres coinciden
  con los perfiles sembrados.
- [x] Probar el destinatario, el contenido del aviso y que no se envíe en una
  transición prohibida.
- [x] Guardar la salida real de las pruebas en
  [`prueba_servicio.txt`](./prueba_servicio.txt).
- [x] Publicar el código y la evidencia del punto 1 con este commit.

- [x] Preparar la página local de demostración en `pagina/`; Persona 2 la usará
  para la evidencia del punto 3.

## Persona 2 · punto 2 y reto

Además, Persona 2 toma la línea base HTTP y el punto 3 según el acuerdo de
reparto registrado en el `README.md`.

### Punto 2 · ruta WebSocket autenticada

- [x] Implementar el adaptador de transporte de avisos (`internal/gymcoach/avisos.go`)
  y registrar `/ws` en `main.go`, fuera del middleware de sesión por cabecera.
- [x] Validar el token de sesión de la query antes de aceptar la conexión; rechazar
  sin token con 401 y origen no autorizado con 403 (`http://localhost:5173`).
- [x] Mantener canales por usuario con mutex, aceptar varias conexiones por
  usuario y enviar cada aviso solo a la clave destinataria recibida del
  servicio.
- [x] Probar desconexión, concurrencia y entrega selectiva; documentar el resultado
  de `go test` en [`prueba_websocket.txt`](./prueba_websocket.txt).
- [x] Guardar salidas reales de `curl` en `ws_sin_token.txt`,
  `ws_otro_origen.txt` y `ws_permitido.txt` con el servidor conectado a PostgreSQL.

### Punto 3 · navegador y regresión reasignados

- [x] Abrir dos ventanas con `"Daniela Vera"` y `"Luis Mera"`; comprobar que el
  dueño ve el aviso de su rutina y que el otro cliente no lo recibe.
- [x] Guardar una captura real en `aviso_navegador.png`.
- [x] Guardar las líneas relevantes del servidor en `logs_avisos.txt` y borrar los
  tokens antes de versionar.
- [x] Completar la línea base HTTP y repetir las seis peticiones después de
  integrar WebSocket; registrar las respuestas reales en [`peticiones.md`](./peticiones.md).
- [x] Escribir D1–D3 en `docs/decisiones.md`.

### Reto · avisar cuando se crea una rutina

- [x] Diseñar el aviso de rutina nueva sin hacer que el servicio de aplicación
  conozca WebSocket ni roles (`AvisadorCreacion` en `internal/gymcoach/aplicacion/rutinas.go`).
- [x] En GymCoach, avisar a las conexiones del rol `entrenador` cuando se crea una
  rutina; no avisar a los clientes ni a la persona que la creó. Mantener sin
  cambios los avisos de cambio de estado al cliente dueño.
- [x] Probar con entrenadores y cliente conectados; registrar la salida en
  [`prueba_reto.txt`](./prueba_reto.txt).
- [x] Repetir las seis peticiones en [`peticiones.md`](./peticiones.md).
- [x] Agregar D4 a `docs/decisiones.md`.

## Secuencia de coordinación

1. Persona 1 termina preparación y punto 1, y publica este commit.
2. Persona 2 actualiza su rama con `git pull`, completa la línea base, el punto
   2 y el punto 3 reasignado; publica los commits correspondientes.
3. Persona 2 completa el reto, sus evidencias y D4, y publica el commit final.

No se deben fabricar evidencias. La captura, las respuestas HTTP y las salidas
de pruebas se guardan únicamente después de ejecutarlas. No se versionan `.env`,
contraseñas reales ni tokens.
