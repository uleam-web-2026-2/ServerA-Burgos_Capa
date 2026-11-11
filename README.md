# GymCoach

**Hilo Servidor · ULEAM · Período 2026-2**
API SaaS para que entrenadores personales administren clientes y rutinas de gimnasio.

## Integrantes

| Integrante | Usuario de GitHub | Paralelo |
| --- | --- | --- |
| Capa Vargas Renato David | Renato-capa | Servidor Web A |
| Burgos Macias Adrian Omar | adrian03 | Servidor Web A |

## Requisitos

- Go compatible con la versión indicada en `go.mod`.
- PostgreSQL 16.
- Node.js 20 o superior para validar y generar la documentación OpenAPI.
- Bash, `curl` y `jq` para ejecutar `pruebas.sh`.

## Configuración y ejecución

1. Cree una base de datos llamada `gymcoach`.
2. Copie `.env.example` a `.env`, reemplace `CLAVE_FIRMA` por un secreto aleatorio largo y ajuste `DATABASE_URL` si su configuración local es distinta.
3. Inicie el servidor:

```bash
go run .
```

El servidor escucha en `http://localhost:8080`, crea las tablas de entrenadores, clientes, rutinas y usuarios, y carga datos ficticios de ejemplo.

### Variables de entorno

La aplicación carga las variables de `.env` al iniciar. Configure:

| Variable | Propósito | Valor local de referencia |
| --- | --- | --- |
| `DATABASE_URL` | Conexión a PostgreSQL. | `postgres://usuario:contraseña@localhost:5432/gymcoach?sslmode=disable` |
| `CLAVE_FIRMA` | Secreto para firmar los tokens de sesión. | Use un valor aleatorio propio; no publique el valor real. |
| `PUERTO` | Puerto HTTP de la API. | `8080` si no se define. |
| `TIEMPO_ESPERA_SEGUNDOS` | Tiempo máximo de espera de las peticiones. | `5` si no se define. |

PostgreSQL debe estar iniciado y accesible en el puerto indicado por
`DATABASE_URL` antes de ejecutar `go run .`. Si la conexión falla, revise el
servicio de PostgreSQL, el nombre de la base, el usuario, la contraseña y el
puerto. No comparta ni suba `.env`; el archivo de ejemplo no debe contener
credenciales reales.

Con Docker puede iniciar PostgreSQL así:

```bash
docker run --name gymcoach-db -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=gymcoach -p 5432:5432 -d postgres:16
```

Si PostgreSQL ya ocupa el puerto local `5432` en Windows, publique Docker en `5433` con `-p 5433:5432` y cambie el puerto de `DATABASE_URL` en `.env` a `5433`.

No reutilice las credenciales de ejemplo fuera de una máquina local. El archivo `.env` está excluido del repositorio.

## API disponible

Todas las respuestas JSON correctas usan la envoltura `ok` y `datos`; los errores usan `ok` y `error`.

| Método y ruta | Función |
| --- | --- |
| `POST /login` | Público. Recibe `usuario` y `clave`; entrega un token Bearer. |
| `GET /clientes` | Solo entrenador. Lista clientes con entrenador y rutinas; admite filtro `?estado=`. |
| `POST /clientes` | Solo entrenador. Crea un cliente asociado a un entrenador existente. |
| `GET /rutinas` | Solo entrenador. Lista rutinas con su cliente; admite filtro `?estado=`. |
| `POST /rutinas` | Solo entrenador. Crea una rutina en estado `pendiente`. |
| `GET /rutinas/{id}` | Requiere sesión. Consulta una rutina y su cliente. |
| `PUT /rutinas/{id}` | Solo entrenador. Actualiza datos de la rutina; no cambia el estado. |
| `PATCH /rutinas/{id}` | Solo cliente. Cambia el estado siguiendo la máquina de estados; transición inválida responde 422. |
| `DELETE /rutinas/{id}` | Solo entrenador. Elimina lógicamente una rutina. |

La secuencia permitida es `pendiente → en_progreso → completada`. Para probar el login localmente, se crean los usuarios ficticios `ana` / `ana-2026` (cliente) y `luis` / `luis-2026` (entrenador). La autorización comprueba el rol, pero todavía no vincula cada usuario con sus propios registros; no publique la API ni la use con datos personales reales.

Las respuestas de rutinas y clientes usan DTOs de la API, no serializan los
modelos completos de GORM ni incluyen sus marcas de tiempo internas. La ruta
`GET /clientes` incluye los datos del entrenador y una lista de rutinas;
`PUT /rutinas/{id}` no acepta `estado` (use `PATCH` para cambiarlo).

## Transacciones de base de datos

La migración de datos antiguos se ejecuta dentro de una transacción GORM en
`internal/gymcoach/migracion.go`. La preparación de clientes y rutinas —incluidos
los valores predeterminados y las relaciones que falten— se confirma en conjunto
si todas las operaciones terminan correctamente. Si una operación falla, se
devuelve el error y PostgreSQL revierte los cambios de esa migración.

El cambio de estado de una rutina usa una estrategia distinta: sqlc ejecuta un
`UPDATE` atómico condicionado por el estado actual (`estado = estado_actual`) y
por la rutina no eliminada. Así se detecta si otra petición modificó el estado
entre la lectura y la escritura; si no se actualiza ninguna fila, el servicio
responde con conflicto en lugar de sobrescribir el cambio concurrente. Esta
operación no envuelve la lectura y la escritura en una transacción explícita.

## Pruebas

```bash
go test ./...
```

El script `pruebas.sh` ejecuta comprobaciones contra `localhost:8080`, inicia sesión con los usuarios ficticios y verifica los códigos HTTP y la transición de estado. El contrato OpenAPI está en [docs/openapi.yaml](docs/openapi.yaml).

Para validar y regenerar la documentación del contrato:

```bash
npx @redocly/cli@1.34.24 lint docs/openapi.yaml --format=stylish
npx @redocly/cli@1.34.24 build-docs docs/openapi.yaml -o docs/api.html
```

Se fija la versión de Redocly CLI para que la documentación generada conserve
el renderizador Redoc 2 usado por `docs/api.html`.

Las credenciales `ana` / `ana-2026` y `luis` / `luis-2026` son ficticias y solo sirven para pruebas locales.

## Hito 1

- [Ficha del negocio](docs/hito1_ficha_del_negocio.md)
- [Addendum técnico](docs/hito1_addendum.md)
- [Presentación en PowerPoint](docs/hito1_presentacion.pptx)
- [Presentación en PDF](docs/hito1_presentacion.pdf)
- [Guion de diapositivas](docs/hito1_presentacion.md)
- [Decisiones de arquitectura](docs/decisiones.md)

## Hito 2

- [Hoja de demostración](docs/hito2_hoja.md)
- [Colección Bruno](docs/hito2_bruno/)
- [Capturas de respaldo](docs/hito2_capturas/)
- [Contrato OpenAPI](docs/openapi.yaml), validado con Redocly CLI 1.34.24.

La guía de `docs/hito2_bruno/` explica cómo preparar los tokens, crear una rutina nueva y guardar su id en `rutina_id`. La secuencia de demostración usa Ana para probar el permiso por rol y el cambio de estado. No reemplaces las capturas por ejemplos: guarda respuestas reales de tu servidor y la salida real de `go test ./...`.

## Estructura

```text
main.go                          → configuración de servidor, PostgreSQL y rutas
internal/config/                 → carga de variables de entorno
internal/gymcoach/               → manejadores HTTP, modelos, migraciones y pruebas
internal/gymcoach/dominio/        → entidades y reglas de negocio sin HTTP
internal/gymcoach/aplicacion/     → casos de uso y DTOs de la API
internal/gymcoach/repositorio/    → persistencia GORM/sqlc de Rutinas
db/                               → esquema y consultas de sqlc
internal/respuesta/              → formato JSON compartido para respuestas
docs/hito1_ficha_del_negocio.md  → ficha del Hito 1
docs/hito1_addendum.md           → addendum técnico del Hito 1
docs/openapi.yaml                → contrato de la API
```
