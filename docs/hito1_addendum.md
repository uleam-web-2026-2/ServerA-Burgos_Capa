# Hito 1 · Addendum técnico

**Pareja:** Capa Vargas Renato David · Burgos Macias Adrian Omar  
**Paralelo:** Aplicación para el Servidor Web A

## A. Estructura del proyecto

```text
.
├── main.go
├── go.mod
├── .env.example
├── internal/
│   ├── config/       → carga de configuración
│   ├── respuesta/    → envoltura JSON común
│   └── gymcoach/     → modelos, reglas, rutas, semilla y pruebas
├── docs/
│   ├── ficha_negocio.md
│   ├── hito1_ficha_del_negocio.md
│   ├── hito1_addendum.md
│   ├── openapi.yaml
│   └── hito1_presentacion.pptx
└── pruebas.sh        → solicitudes manuales a localhost:8080
```

`main.go` carga la configuración, abre PostgreSQL, crea primero la tabla `Entrenador`, adapta los datos existentes de versiones anteriores y luego migra las tablas `Cliente` y `Rutina`. Si encuentra campos nuevos vacíos, asigna valores temporales identificables y conserva los registros existentes; revise los datos con `.invalid` y complete la información real. Después carga datos de ejemplo solo si no existen entrenadores y registra las rutas de `internal/gymcoach`. El dominio `.invalid` está reservado y no representa una dirección entregable.

## B. Configuración y secretos

| Variable | Para qué sirve | Ejemplo (sin datos reales) |
|----------|----------------|----------------------------|
| `DATABASE_URL` | Conexión a PostgreSQL | `postgres://postgres:postgres@localhost:5432/gymcoach?sslmode=disable` |
| `PUERTO` | Puerto HTTP | `8080` |
| `TIEMPO_ESPERA_SEGUNDOS` | Tiempo máximo de espera del servidor | `5` |
| `CLAVE_FIRMA` | Reserva para firma de tokens cuando se implemente autenticación | `replace-with-a-random-local-secret` |

**Archivo de ejemplo:** `.env.example`. `.env` está excluido del repositorio. Los valores de ejemplo son solo para desarrollo local.

## C. Pruebas

| Prueba | Qué caso cubre |
|--------|----------------|
| `TestCrearRutinaValidaDatosAntesDeAccederABase` | JSON roto, campo desconocido, título vacío, días fuera de rango y cliente ausente. |
| `TestIDNoNumericoResponde400` | Rechazo de ID de ruta no numérico. |
| `TestTransicionesDeRutina` | Transiciones permitidas y prohibidas de la rutina. |

**Comando:** `go test ./...`

**Resultado conocido en esta máquina:** la ejecución no terminó porque Windows Application Control bloqueó `C:\Program Files\Go\pkg\tool\windows_amd64\link.exe`. `go vet ./...` sí se ejecutó. Se debe volver a correr `go test ./...` en un entorno donde el enlazador de Go esté permitido y adjuntar la salida real; no se presenta una prueba bloqueada como aprobada.

**Captura de `go test ./...`:** pendiente de una corrida exitosa.

## D. Boceto de la pantalla principal

![Boceto del panel del entrenador](./hito1_boceto.svg)

## E. Diagrama de secuencia del caso de uso principal

```mermaid
sequenceDiagram
    actor Entrenador
    participant Panel
    participant API
    participant BD as PostgreSQL
    Entrenador->>Panel: Selecciona cliente y crea una rutina
    Panel->>API: POST /rutinas
    API->>BD: Comprueba el cliente y guarda la rutina
    BD-->>API: Rutina creada en estado pendiente
    API-->>Panel: 201 con la rutina
    Panel-->>Entrenador: Confirma que el plan fue asignado
```

## F. Capturas de respuestas

Para generar las capturas desde la raíz del proyecto, abre Docker Desktop. En la primera terminal de PowerShell inicia PostgreSQL (en ejecuciones posteriores usa `docker start gymcoach-db`):

```powershell
if (-not (docker ps -a --format "{{.Names}}" | Select-String -Quiet "^gymcoach-db$")) {
  docker run --name gymcoach-db -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=gymcoach -p 5432:5432 -d postgres:16
} else {
  docker start gymcoach-db
}
if (-not (Test-Path .env)) { Copy-Item .env.example .env }
go run .
```

Deja `go run .` activo y abre una segunda terminal PowerShell en la raíz del proyecto. Obtén el ID de un cliente de ejemplo y solicita una rutina válida:

```powershell
$clientes = Invoke-RestMethod http://localhost:8080/clientes
$clienteId = $clientes.datos[0].id
$body = @{ titulo = "Rutina de prueba"; descripcion = "Captura del Hito 1"; dias_semana = 3; cliente_id = $clienteId } | ConvertTo-Json -Compress
curl.exe -i -X POST http://localhost:8080/rutinas -H "Content-Type: application/json" --data-raw $body
```

La respuesta debe incluir `HTTP/1.1 201 Created` y un JSON con `"ok":true` y la rutina en `"datos"`. Para capturar el error de validación con días fuera del rango 1–7, ejecuta:

```powershell
$body = @{ titulo = "Rutina de prueba"; descripcion = "Días fuera de rango"; dias_semana = 8; cliente_id = $clienteId } | ConvertTo-Json -Compress
curl.exe -i -X POST http://localhost:8080/rutinas -H "Content-Type: application/json" --data-raw $body
```

La respuesta debe incluir `HTTP/1.1 422 Unprocessable Entity` y un error con código `dias_semana_invalidos`. Para cada captura, deja visible en la terminal el comando, el código HTTP y el cuerpo JSON; usa `Win + Shift + S` para tomar el recorte y guarda las imágenes como `docs/captura_post_rutinas_201.png` y `docs/captura_post_rutinas_422.png`. Adjunta esas capturas al addendum y a la presentación. No escribas ni pegues una respuesta esperada como si fuera una captura real.

En este equipo Docker Desktop no estaba iniciado y Windows Application Control bloqueó `link.exe` al intentar compilar la aplicación. Por eso no se adjuntan respuestas reales todavía. Si `go run .` vuelve a mostrar ese bloqueo, solicita al administrador del equipo que permita el enlazador de Go; si falla la conexión a PostgreSQL, confirma primero que Docker Desktop haya iniciado el contenedor y que el puerto 5432 esté disponible.

## Presentación

La plantilla oficial de cinco diapositivas está completada en `docs/hito1_presentacion.pptx` y exportada a `docs/hito1_presentacion.pdf`. Las capturas reales de endpoints y de pruebas siguen pendientes de ejecutar en una máquina que permita correr el enlazador de Go y conectar PostgreSQL.
