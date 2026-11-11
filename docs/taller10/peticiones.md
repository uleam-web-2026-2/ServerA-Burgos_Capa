# Taller 10 · Peticiones de verificación

Base: `http://localhost:8080`. Credenciales ficticias del entorno local.

«Antes» se ejecutó con el commit base `a756a7c`. Los resultados «Después» son
del registro de verificación previo con sqlc y una versión anterior de los
cambios DTO; no representan una repetición de las seis peticiones contra los
últimos ajustes de DTO y OpenAPI.

| # | Petición                                                          | Resultado esperado           | Antes                        | Después                     |
| - | ------------------------------------------------------------------ | ---------------------------- | ---------------------------- | ---------------------------- |
| 1 | `GET /rutinas/{rutina_id}` sin token                             | 401                          | 401                          | 401                          |
| 2 | `POST /login` con Ana                                            | 200 y token                  | 200; token recibido          | 200; token recibido          |
| 3 | `GET /rutinas/{rutina_id}` con token de Ana                      | 200                          | 200                          | 200                          |
| 4 | `PATCH /rutinas/{rutina_id}` a `en_progreso` con token de Luis | 403                          | 403                          | 403                          |
| 5 | `PATCH /rutinas/{rutina_id}` a `en_progreso` con token de Ana  | 200                          | 200; estado `en_progreso`  | 200; estado `en_progreso`  |
| 6 | `PATCH /rutinas/{rutina_id}` a `pendiente` con token de Ana    | 422 `transicion_prohibida` | 422 `transicion_prohibida` | 422 `transicion_prohibida` |

## Entorno y alcance de la verificación

El PostgreSQL local del equipo usa el puerto `8080`, por lo que las pruebas se
ejecutaron con un clúster PostgreSQL 16 temporal y aislado en `5432`. La API
base corrió en `18080`; la API con sqlc y los cambios DTO locales corrió en
`18081`. Cada secuencia creó una rutina nueva y comprobó estado HTTP y campos
relevantes de las respuestas. No se imprimieron ni guardaron tokens.

El esquema real se inspeccionó con `\d rutinas`, `\d clientes` y
`\d entrenadores`; `db/schema.sql` refleja sus columnas, índices y relaciones.
La Fase 4 debe repetirse con PostgreSQL disponible contra la revisión exacta que
se entregará, después de integrar y publicar los cambios finales.

## Validaciones de los últimos cambios

- `go test ./...` finalizó correctamente después de los últimos cambios.
- `npx @redocly/cli@1.34.24 lint docs/openapi.yaml --format=stylish` validó el contrato.
- Se regeneró `docs/api.html` con
  `npx @redocly/cli@1.34.24 build-docs docs/openapi.yaml -o docs/api.html`.
- Las pruebas automatizadas verifican la conversión del DTO específico de lista
  de clientes, la omisión de campos internos y el rechazo de `estado` en el
  cuerpo de actualización `PUT`.

La repetición manual de las seis peticiones contra PostgreSQL queda pendiente
para confirmar el comportamiento de la API desplegada con esta revisión final.
