# Taller 10 · Peticiones de verificación

Base: `http://localhost:8080`. Credenciales ficticias del entorno local.

La columna «Antes» corresponde a la línea de base previa a la integración sqlc.
La columna «Después» se completará en la Fase 4, cuando estén disponibles los
cambios de DTO de la Persona B.

| # | Petición | Resultado esperado | Antes | Después |
|---|---|---|---|---|
| 1 | `GET /rutinas/{rutina_id}` sin token | 401 | No ejecutada: la conexión a `localhost:8080` se reinicia | Pendiente (Fase 4) |
| 2 | `POST /login` con Ana | 200 y token | No ejecutada: la conexión a `localhost:8080` se reinicia | Pendiente (Fase 4) |
| 3 | `GET /rutinas/{rutina_id}` con token de Ana | 200 | No ejecutada: la conexión a `localhost:8080` se reinicia | Pendiente (Fase 4) |
| 4 | `PATCH /rutinas/{rutina_id}` a `en_progreso` con token de Luis | 403 | No ejecutada: la conexión a `localhost:8080` se reinicia | Pendiente (Fase 4) |
| 5 | `PATCH /rutinas/{rutina_id}` a `en_progreso` con token de Ana | 200 | No ejecutada: la conexión a `localhost:8080` se reinicia | Pendiente (Fase 4) |
| 6 | `PATCH /rutinas/{rutina_id}` a `pendiente` con token de Ana | 422 `transicion_prohibida` | No ejecutada: la conexión a `localhost:8080` se reinicia | Pendiente (Fase 4) |

## Bloqueo del entorno

En la verificación, `localhost:8080` estaba ocupado por el servicio local
PostgreSQL 16 y reiniciaba las conexiones HTTP. La configuración local de la
aplicación apunta a `localhost:5432`, donde no había un servicio escuchando.
La autenticación al PostgreSQL de `localhost:8080` con las credenciales
configuradas fue rechazada. Por eso las seis peticiones HTTP no se completaron
y no se registran códigos como si fueran resultados. Tampoco se pudo verificar
el esquema con `\d rutinas`. Repetir la secuencia cuando la API y la base de
datos estén disponibles con sus credenciales correctas.
