# Taller 10 · Peticiones de verificación

Base: `http://localhost:8080`. Credenciales ficticias del entorno local.

«Antes» se ejecutó con el commit base `a756a7c`. «Después» se ejecutó en el
árbol de trabajo actual con la integración sqlc y los cambios DTO locales de la
Persona 2, todavía no publicados en `origin/taller10`.

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
La Fase 4 oficial se debe repetir después de que la Persona 2 publique sus
cambios y la Persona 1 actualice la rama compartida.

Verificación final de cierre con la refactorización DTO de la Persona 2:

- `go test ./...` se ejecutó correctamente y no presentó fallos en la
  integración del repositorio sqlc con la capa de aplicación.
- Las seis peticiones de la Fase 4 mantuvieron los resultados esperados en el
  entorno de verificación: códigos HTTP, transiciones de estado y mensajes de
  error coherentes con la especificación.
- La refactorización de DTO quedó validada: los handlers decodifican entradas y
  responden con DTOs, sin exponer los modelos de dominio ni los campos internos
  de GORM al contrato JSON de la API.
- La integración entre sqlc y la capa HTTP/aplicación quedó confirmada como
  estable y compatible con la lógica de negocio actual.
