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

## Pruebas

```bash
go test ./...
```

El script `pruebas.sh` ejecuta comprobaciones contra `localhost:8080`, inicia sesión con los usuarios ficticios y verifica los códigos HTTP y la transición de estado. El contrato OpenAPI está en [docs/openapi.yaml](docs/openapi.yaml).

Para validar y regenerar la documentación del contrato:

```bash
npx @redocly/cli lint docs/openapi.yaml --format=stylish
npx @redocly/cli build-docs docs/openapi.yaml -o docs/api.html
```

Las credenciales `ana` / `ana-2026` y `luis` / `luis-2026` son ficticias y solo sirven para pruebas locales.

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
