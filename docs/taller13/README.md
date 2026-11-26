# Taller 13 · Persona 1 y Persona 2

- [Enunciado oficial](./Semana_13_-_Taller.docx)
- [Reparto detallado y checklist](./reparto.md)
- [Tabla de regresión HTTP](./peticiones.md)
- [Pendientes y bloqueos](./pendientes.md)

## Persona 1 · trabajo realizado

- Se creó la rama `taller13` desde `origin/main`, que ya contiene el Taller 12.
- Se agregó `github.com/coder/websocket` como dependencia para el trabajo de
  transporte que realizará Persona 2.
- Se confirmó que `Rutina.ClienteID` ya identifica al cliente dueño. En
  `internal/gymcoach/aplicacion` se definió `Aviso` y la interfaz `Avisador`;
  el servicio notifica al dueño solo cuando el cambio de estado se guardó.
- Se agregaron pruebas para confirmar el destinatario y el contenido del aviso,
  y que no se envía en transiciones prohibidas ni si falla la escritura. La
  salida está en [prueba_servicio.txt](./prueba_servicio.txt).
- Se agregaron cuentas de demostración para dos clientes con nombres iguales
  a los perfiles sembrados.
- Se preparó la página local para iniciar sesión y abrir una conexión en
  [`pagina/`](./pagina/).
- `go test ./...` y `go build ./...` pasan.

Se conservan los permisos de GymCoach: el entrenador crea y asigna rutinas, y
solo el cliente cambia su estado. No se modificó el contrato existente del
Hito 2.

El trabajo de código y pruebas del punto 1 de Persona 1 se incluye en este
commit. La línea base y las pruebas restantes se reasignan a Persona 2 según
el acuerdo de la pareja. Aunque el enunciado propone que Persona 1 haga el
punto 3 después del punto 2, en este equipo también lo realizará Persona 2.

## Persona 2 · trabajo asignado

### Punto 2 · ruta WebSocket autenticada

- Implementar el adaptador de avisos y registrar `/ws` fuera del guardia que
  exige `Authorization`.
- Validar el token de la query antes de aceptar la conexión; responder 401 sin
  token y 403 si el origen no está permitido. Usar el mismo origen de CORS:
  `http://localhost:5173`.
- Mantener canales separados por usuario, protegidos con mutex, y entregar el
  aviso solo a su destinatario. Probar concurrencia, varias conexiones,
  desconexión y entrega selectiva.
- Guardar las pruebas reales en `prueba_websocket.txt` y las respuestas de
  `curl` en `ws_sin_token.txt`, `ws_otro_origen.txt` y `ws_permitido.txt`.
- Validar y publicar el commit del punto 2 para continuar con las pruebas del
  punto 3 reasignado y el reto.

### Línea base y punto 3 · trabajo reasignado a Persona 2

- Revisar por qué la API local no respondió y confirmar que PostgreSQL esté
  disponible. Completar la línea base con las seis respuestas HTTP reales en
  [peticiones.md](./peticiones.md); no registrar los intentos fallidos como
  respuestas de la API.
- Después de integrar `/ws`, usar la página preparada en `pagina/` con dos
  clientes conectados. Verificar que el aviso llegue solo al dueño de la
  rutina y guardar la captura en `aviso_navegador.png`.
- Guardar las líneas reales del servidor en `logs_avisos.txt` y eliminar los
  tokens antes de versionar.
- Repetir las seis peticiones tras WebSocket, registrar los códigos reales y
  escribir D1–D3 en `docs/decisiones.md`.
- Hacer el commit de esta parte después de completar las pruebas y evidencias.

### Reto · aviso de rutina nueva

- Notificar a las conexiones del rol entrenador al crear una rutina, sin
  acoplar la capa de aplicación a WebSocket ni a roles.
- No notificar a los clientes ni a quien creó la rutina; conservar los avisos
  de cambio de estado dirigidos al cliente dueño.
- Probar entrega con dos entrenadores y un cliente; guardar
  `prueba_reto.txt`, repetir las seis peticiones, escribir D4 y hacer el commit
  del reto.

No se deben fabricar evidencias ni versionar `.env`, contraseñas reales o
tokens.
