# Hito 2 · Hoja de demostración

## 1 · Usuarios de ejemplo

| Usuario | Rol | Contraseña de ejemplo |
| --- | --- | --- |
| ana | cliente | ana-2026 |
| luis | entrenador | luis-2026 |

Son credenciales ficticias para la demostración local.

## 2 · Las seis peticiones

| # | Petición | Credencial | Debe responder |
| --- | --- | --- | --- |
| 1 | `GET /rutinas/{rutina_id}` | Ninguna | 401 |
| 2 | `POST /login` con ana | Ninguna | 200 y token |
| 3 | `GET /rutinas/{rutina_id}` | Token de Ana | 200 |
| 4 | `PATCH /rutinas/{rutina_id}` → `en_progreso` | Token de Luis (entrenador) | 403 |
| 5 | La misma transición a `en_progreso` | Token de Ana (cliente) | 200 |
| 6 | `PATCH /rutinas/{rutina_id}` → `pendiente` | Token de Ana (cliente) | 422 `transicion_prohibida` |

La petición 0 de preparación crea la rutina en estado `pendiente`. Como la petición 5 la cambia a `en_progreso`, la petición 6 prueba que el negocio no permite devolverla a `pendiente` y perder el avance. Antes de la demostración, inicia sesión como Luis con `preparacion/Login Entrenador.bru`; el token se guarda en `token_entrenador`.

## 3 · Matriz ruta × rol

| Ruta | Sin sesión | Cliente | Entrenador |
| --- | --- | --- | --- |
| `POST /login` | Sí | Sí | Sí |
| `GET /clientes` | No (401) | No (403) | Sí |
| `POST /clientes` | No (401) | No (403) | Sí |
| `GET /rutinas` | No (401) | No (403) | Sí |
| `POST /rutinas` | No (401) | No (403) | Sí |
| `GET /rutinas/{id}` | No (401) | Sí | Sí |
| `PUT /rutinas/{id}` | No (401) | No (403) | Sí |
| `PATCH /rutinas/{id}` | No (401) | Sí | No (403) |
| `DELETE /rutinas/{id}` | No (401) | No (403) | Sí |

## 4 · Recuadro de cumplimiento

- [ ] Los structs del modelo están en el código y sus tablas se crean al iniciar.
- [ ] `docs/openapi.yaml` incluye las rutas del mapa y sus respuestas 401, 403 y 422.
- [ ] `npx @redocly/cli lint docs/openapi.yaml --format=stylish` termina sin errores ni avisos.
- [ ] Las peticiones 0–6 están en `docs/hito2_bruno/`, en orden.
- [ ] Cada petición responde lo indicado en esta hoja.
- [ ] `go test ./...` termina en verde e incluye una prueba de integración.
- [ ] `docs/hito2_capturas/` contiene capturas reales de las respuestas y de las pruebas.
- [ ] La regla del cambio de estado está en dominio; el manejador no consulta la base directamente.
- [ ] `docs/decisiones.md` explica la estructura por capas.
- [ ] El repositorio no contiene contraseñas reales ni el secreto de firma del token.

## 5 · Evidencias

Guarda las capturas reales en `docs/hito2_capturas/` con estos nombres. Al agregarlas, aparecerán en los espacios siguientes:

### Evidencia 1 · GET sin token (401)

Archivo: `01-sin-token.png`

![GET sin token: respuesta 401](hito2_capturas/01-sin-token.png)

### Evidencia 2 · Inicio de sesión de Ana (200)

Archivo: `02-login-cliente.png`

![Inicio de sesión de Ana: respuesta 200 con token](hito2_capturas/02-login-cliente.png)

### Evidencia 3 · GET con token de Ana (200)

Archivo: `03-rutina-con-token.png`

![Consulta de rutina autenticada: respuesta 200](hito2_capturas/03-rutina-con-token.png)

### Evidencia 4 · Entrenador intenta cambiar el estado (403)

Archivo: `04-entrenador-prohibido.png`

![Cambio de estado por entrenador: respuesta 403](hito2_capturas/04-entrenador-prohibido.png)

### Evidencia 5 · Ana cambia el estado válidamente (200)

Archivo: `05-transicion-valida.png`

![Transición válida por Ana: respuesta 200](hito2_capturas/05-transicion-valida.png)

### Evidencia 6 · Transición prohibida (422)

Archivo: `06-transicion-prohibida.png`

![Transición prohibida: respuesta 422](hito2_capturas/06-transicion-prohibida.png)

### Evidencia 7 · Pruebas automatizadas

Archivo: `07-go-test.png`

![Ejecución de go test ./... en verde](hito2_capturas/07-go-test.png)
