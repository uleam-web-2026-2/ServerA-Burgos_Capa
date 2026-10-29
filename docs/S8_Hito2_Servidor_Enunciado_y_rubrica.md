# Proyecto autónomo · Hito 2 — Su API, demostrada en vivo

**Proyecto:** GymCoach  
**Equipo:** Capa Vargas Renato David (`Renato-capa`) y Burgos Macias Adrian Omar (`adrian03`)  
**Paralelo:** Servidor Web A  
**API local:** `http://localhost:8080`

Esta es la guía de demostración de GymCoach, completada con las rutas, usuarios y reglas que existen en el proyecto. Las capturas indicadas son espacios para añadir imágenes reales tomadas al ejecutar la API; no se han generado imágenes de ejemplo.

## Preparación

- Iniciar PostgreSQL y el servidor con `go run .`.
- En Bruno, seleccionar el entorno `Local` y configurar `cliente_id` con un cliente existente.
- Ejecutar `preparacion/Login Entrenador.bru` para obtener un token de entrenador.
- Ejecutar `0 - Preparar rutina.bru`. La petición crea una rutina nueva en estado `pendiente` y guarda su ID en `rutina_id`.
- Ejecutar las peticiones 1 a 6 en orden. La petición 2 obtiene un token nuevo de Ana; no reutilizar tokens de ensayos anteriores.
- Antes de cada demostración, volver a ejecutar la petición 0 para que la petición 5 use una rutina que todavía esté `pendiente`.
- Ejecutar `npx @redocly/cli lint docs/openapi.yaml --format=stylish` y `go test ./...`.

## 1 · Usuarios de ejemplo

| Usuario | Rol | Contraseña local de demostración |
|---|---|---|
| `ana` | Cliente | `ana-2026` |
| `luis` | Entrenador | `luis-2026` |

Son credenciales ficticias sembradas para pruebas locales. No usarlas en producción.

## 2 · Las seis peticiones

La petición 0 de preparación crea la rutina con el entrenador Luis. Su ID se guarda como `rutina_id`.

| # | Petición | Credencial | Respuesta esperada |
|---|---|---|---|
| 1 | `GET /rutinas/{{rutina_id}}` | Sin token | `401` |
| 2 | `POST /login` con Ana | Sin token | `200` y token Bearer |
| 3 | `GET /rutinas/{{rutina_id}}` | Token de Ana | `200` |
| 4 | `PATCH /rutinas/{{rutina_id}}` con `{"estado":"en_progreso"}` | Token de Luis (entrenador) | `403` |
| 5 | La misma transición a `en_progreso` | Token de Ana (cliente) | `200` |
| 6 | `PATCH /rutinas/{{rutina_id}}` con `{"estado":"pendiente"}` | Token de Ana (cliente) | `422`, código `transicion_prohibida` |

La petición 5 avanza la rutina de `pendiente` a `en_progreso`. La petición 6 intenta retrocederla, transición que el negocio no permite.

## 3 · Matriz ruta × rol

| Ruta | Sin sesión | Cliente | Entrenador |
|---|---:|---:|---:|
| `POST /login` | Permitido | Permitido | Permitido |
| `GET /clientes` | `401` | `403` | Permitido |
| `POST /clientes` | `401` | `403` | Permitido |
| `GET /rutinas` | `401` | `403` | Permitido |
| `POST /rutinas` | `401` | `403` | Permitido |
| `GET /rutinas/{id}` | `401` | Permitido | Permitido |
| `PUT /rutinas/{id}` | `401` | `403` | Permitido |
| `PATCH /rutinas/{id}` | `401` | Permitido | `403` |
| `DELETE /rutinas/{id}` | `401` | `403` | Permitido |

La autorización implementada comprueba el rol, pero todavía no relaciona la cuenta de sesión con un cliente o entrenador específico ni filtra registros por propietario.

## 4 · Resumen de cumplimiento y evidencias

- [x] Los modelos de entrenador, cliente y rutina están implementados con sus relaciones; el servidor crea sus tablas al iniciar.
- [x] `docs/openapi.yaml` declara el esquema Bearer, documenta las rutas implementadas y las respuestas de autenticación y validación.
- [x] `redocly.yaml` tiene habilitada la regla `security-defined`.
- [x] La colección `docs/hito2_bruno/` contiene la preparación y las peticiones 1 a 6.
- [x] `go test ./...` finalizó correctamente en la revisión del proyecto.
- [x] Hay una prueba de integración del router y la autenticación sin token en `internal/gymcoach/sin_token_test.go`.
- [x] La regla de transición está en dominio; el manejador HTTP delega las operaciones de rutina al servicio.
- [x] `docs/decisiones.md` explica la estructura por capas.
- [ ] Añadir las capturas reales de las seis respuestas y de `go test ./...` en `docs/hito2_capturas/`.
- [ ] Confirmar que los cambios finales están committeados y subidos al repositorio de la organización antes del cierre.

### Evidencias de funcionamiento

#### Captura 1 · Ruta protegida sin token — `401`

**Aquí poner la captura real del funcionamiento:** respuesta de la petición 1 en Bruno.

![Aquí poner la captura real de GET /rutinas/{rutina_id} sin token, con respuesta 401](hito2_capturas/01-sin-token.png)

#### Captura 2 · Inicio de sesión de Ana — `200`

**Aquí poner la captura real del funcionamiento:** respuesta de la petición 2 con el token recibido. Ocultar el valor completo del token en la imagen.

![Aquí poner la captura real del inicio de sesión de Ana con respuesta 200](hito2_capturas/02-login-cliente.png)

#### Captura 3 · Consulta autenticada — `200`

**Aquí poner la captura real del funcionamiento:** respuesta de la petición 3 con el token de Ana.

![Aquí poner la captura real de GET /rutinas/{rutina_id} con token, con respuesta 200](hito2_capturas/03-rutina-con-token.png)

#### Captura 4 · Rol entrenador no autorizado para cambiar estado — `403`

**Aquí poner la captura real del funcionamiento:** respuesta de la petición 4 con el token de Luis.

![Aquí poner la captura real del PATCH con token de entrenador, con respuesta 403](hito2_capturas/04-entrenador-prohibido.png)

#### Captura 5 · Transición permitida — `200`

**Aquí poner la captura real del funcionamiento:** respuesta de la petición 5 con el token de Ana y estado `en_progreso`.

![Aquí poner la captura real de la transición pendiente a en_progreso, con respuesta 200](hito2_capturas/05-transicion-valida.png)

#### Captura 6 · Transición prohibida — `422`

**Aquí poner la captura real del funcionamiento:** respuesta de la petición 6 con el código `transicion_prohibida`.

![Aquí poner la captura real del intento de volver a pendiente, con respuesta 422](hito2_capturas/06-transicion-prohibida.png)

#### Captura 7 · Pruebas automatizadas

**Aquí poner la captura real del funcionamiento:** terminal mostrando la ejecución completa de `go test ./...` en verde.

![Aquí poner la captura real de go test ./... en verde](hito2_capturas/07-go-test.png)

## Frases preparadas para la demostración

1. **El contrato:** «Este es nuestro contrato. Redocly lo valida sin errores ni avisos. Aquí está declarado BearerAuth y aquí se documentan las respuestas 401, 403 y 422».
2. **Quién es:** «Sin token, `GET /rutinas/{rutina_id}` responde 401. Iniciamos sesión como Ana y recibimos un token. Con ese token, la consulta responde 200».
3. **Qué puede hacer:** «Luis es entrenador y no puede cambiar el estado: recibe 403. Ana es cliente y sí puede iniciar la rutina: recibe 200».
4. **La regla del negocio:** «Pasar de `en_progreso` a `pendiente` no está permitido en GymCoach: responde 422. Las pruebas terminan en verde».
5. **La estructura:** «La petición HTTP se atiende en `internal/gymcoach/`; el servicio de aplicación coordina el caso de uso, el repositorio accede a PostgreSQL y la regla de transición está en `internal/gymcoach/dominio/estado.go`».

## Respuestas breves para las preguntas

- **¿En qué línea se exige el token?** En `internal/gymcoach/manejadores.go`, `Rutas` crea `conSesion := r.With(tokens.Autenticar)` y registra allí las rutas protegidas. `POST /login` se registra aparte como pública en `main.go`.
- **¿Qué lleva el token?** Un JWT firmado con HS256 que contiene el nombre, el rol, el ID de usuario como `sub`, la fecha de emisión y la expiración. Si se altera el rol, la firma deja de coincidir y la petición recibe `401`.
- **¿Por qué solo el cliente puede hacer PATCH?** Porque el cliente es quien registra su progreso; el entrenador administra el contenido del plan. La implementación actual valida el rol, pero no verifica propiedad individual.
- **¿Qué prueba avisa si se quita la protección?** `TestSinTokenLasRutasProtegidasResponden401`, en `internal/gymcoach/sin_token_test.go`.
- **¿Por qué 422 y no 400 o 403?** El JSON y el estado son comprensibles, la persona tiene el rol permitido, pero esa transición viola la regla de negocio. La validación ocurre en el dominio y el manejador traduce el error a HTTP 422.
- **¿Por qué estas capas?** Separa transporte HTTP, casos de uso, reglas de negocio y persistencia. Si cambia PostgreSQL, el cambio principal queda en el repositorio y en la composición de dependencias de `main.go`; la regla pura del dominio no depende de la base ni de HTTP.
