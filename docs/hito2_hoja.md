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

La petición 0 de preparación crea la rutina en estado `pendiente`. Como la petición 5 la cambia a `en_progreso`, la petición 6 prueba que el negocio no permite retroceder a `pendiente`. Antes de la demostración, inicia sesión como Luis con la petición existente `POST Login Entrenador` y guarda el token en `token_entrenador`.

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
