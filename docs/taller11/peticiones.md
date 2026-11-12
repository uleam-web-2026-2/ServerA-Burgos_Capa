# Validación de las seis peticiones del Hito 2

Las pruebas se ejecutaron contra la rama `taller11` el 11 de noviembre de 2026,
usando una base PostgreSQL temporal y aislada. No se usaron las bases locales
existentes ni se guardaron contraseñas o tokens.

El puerto de host `8080` estaba ocupado por un proceso PostgreSQL de Windows.
Para evitar interferencias, la API independiente se publicó en `18081` y la API
de Compose en `18080`; ambos puertos internos coincidieron con `PUERTO`.

| # | Petición | Esperado | Con `go run` | Contenedor independiente | Compose |
|---|---|---:|---|---:|---:|
| 1 | `GET /rutinas/{id}` sin token | 401 | No medido en esta verificación | 401 | 401 |
| 2 | `POST /login` con Ana | 200 | No medido en esta verificación | 200; token recibido, no registrado | 200; token recibido, no registrado |
| 3 | `GET /rutinas/{id}` con token de Ana | 200 | No medido en esta verificación | 200 | 200 |
| 4 | `PATCH /rutinas/{id}` a `en_progreso` con token de Luis | 403 | No medido en esta verificación | 403 | 403 |
| 5 | La misma transición con token de Ana | 200 | No medido en esta verificación | 200; estado `en_progreso` | 200; estado `en_progreso` |
| 6 | `PATCH /rutinas/{id}` a `pendiente` con token de Ana | 422 `transicion_prohibida` | No medido en esta verificación | 422 `transicion_prohibida` | 422 `transicion_prohibida` |

Las peticiones 1–6 se ejecutaron en orden con una rutina nueva en cada entorno.
Las peticiones auxiliares de preparación (login de Luis, listado de clientes y
creación de la rutina) no forman parte de la tabla. La respuesta de login no se
imprimió ni guardó; los tokens se mantuvieron solo en memoria durante la prueba.

## Respuesta completa de la petición 3

Respuesta observada en el contenedor independiente:

```json
{"ok":true,"datos":{"id":4,"titulo":"Validacion contenedor Taller 11","descripcion":"Prueba aislada con PostgreSQL externo","dias_semana":3,"estado":"pendiente","cliente_id":1,"cliente":{"id":1,"nombre":"Daniela Vera","correo":"daniela@example.test","objetivo":"Mejorar fuerza general","entrenador_id":1}}}
```

La respuesta no contiene campos de fecha/hora, por lo que no se observó una
diferencia de zona horaria en la petición 3. La petición 5 confirmó el estado
`en_progreso`; la petición 6 confirmó el error `transicion_prohibida`.

## Entornos probados

- **Contenedor independiente:** imagen `gymcoach-taller11-validation-api:latest`,
  ejecutada con el archivo local ignorado `.env.docker`, conectada a un
  PostgreSQL temporal por `host.docker.internal:15432`, publicada en
  `127.0.0.1:18081`.
- **Compose:** servicios `api` y `postgres` en el proyecto aislado
  `gymcoach-taller11-validation`, con PostgreSQL saludable y API publicada en
  `127.0.0.1:18080`.
- Las evidencias de inicio sin configuración, arranque del contenedor, logs y
  estado de Compose están en esta carpeta. No se registró una repetición actual
  de la línea base con `go run`; esa comparación debe completarse si la rúbrica
  exige una columna medida para ese entorno.
