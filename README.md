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

## Configuración y ejecución

1. Cree una base de datos llamada `gymcoach`.
2. Copie `.env.example` a `.env` y ajuste `DATABASE_URL` si su configuración local es distinta.
3. Inicie el servidor:

```bash
go run .
```

El servidor escucha en `http://localhost:8080`, crea las tablas `Entrenador`, `Cliente` y `Rutina` y agrega datos ficticios de ejemplo si todavía no hay entrenadores.

Con Docker puede iniciar PostgreSQL así:

```bash
docker run --name gymcoach-db -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=gymcoach -p 5432:5432 -d postgres:16
```

No reutilice las credenciales de ejemplo fuera de una máquina local. El archivo `.env` está excluido del repositorio.

## API disponible

Todas las respuestas JSON correctas usan la envoltura `ok` y `datos`; los errores usan `ok` y `error`.

| Método y ruta | Función |
| --- | --- |
| `GET /clientes` | Lista clientes con entrenador y rutinas; admite `?estado=pendiente`, `?estado=en_progreso` o `?estado=completada`. |
| `POST /clientes` | Crea un cliente asociado a un entrenador existente. |
| `POST /rutinas` | Crea una rutina en estado `pendiente`. |
| `GET /rutinas/{id}` | Consulta una rutina y su cliente. |
| `PUT /rutinas/{id}` | Reemplaza los datos de una rutina y valida su transición de estado. |
| `DELETE /rutinas/{id}` | Elimina lógicamente una rutina. |

La secuencia permitida es `pendiente → en_progreso → completada`. El contrato esperado de roles se documenta en [docs/hito1_ficha_del_negocio.md](docs/hito1_ficha_del_negocio.md), pero autenticación y autorización todavía no están implementadas: no publique la API en Internet ni la use con datos personales reales.

## Pruebas

```bash
go test ./...
```

El script `pruebas.sh` contiene solicitudes manuales para probar la API contra `localhost:8080`. El contrato OpenAPI está en [docs/openapi.yaml](docs/openapi.yaml).

## Hito 1

- [Ficha del negocio](docs/hito1_ficha_del_negocio.md)
- [Addendum técnico](docs/hito1_addendum.md)
- [Presentación en PowerPoint](docs/hito1_presentacion.pptx)
- [Presentación en PDF](docs/hito1_presentacion.pdf)
- [Guion de diapositivas](docs/hito1_presentacion.md)
- [Decisiones de arquitectura](docs/decisiones.md)

## Estructura

```text
main.go                          → configuración de servidor, PostgreSQL y rutas
internal/config/                 → carga de variables de entorno
internal/gymcoach/               → entidades, transiciones, manejadores y pruebas
internal/respuesta/              → formato JSON compartido para respuestas
docs/hito1_ficha_del_negocio.md  → ficha del Hito 1
docs/hito1_addendum.md           → addendum técnico del Hito 1
docs/openapi.yaml                → contrato de la API
```
