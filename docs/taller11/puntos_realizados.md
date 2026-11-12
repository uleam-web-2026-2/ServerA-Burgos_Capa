# Taller 11 · Control de puntos realizados y pendientes

Este archivo registra el avance real de la pareja en la rama `taller11`. No se
marcan como terminadas tareas cuya ejecución o evidencia todavía no se haya
comprobado. La guía completa de responsabilidades y criterios está en
[README.md](./README.md); el enunciado oficial es
[S11_06_Taller.docx](./S11_06_Taller.docx).

## Equipo y reparto

| Rol | Integrante | Alcance |
| --- | --- | --- |
| Persona 1 | __________________ | Dockerfile, `.dockerignore`, construcción y evidencias de la imagen; colaboración en el diseño/revisión de Compose. |
| Persona 2 | __________________ | Variables para el contenedor, peticiones HTTP, logs y documentación D1–D4; colaboración en la configuración de la API dentro de Compose. |

Las evidencias de Persona 2 y de Compose ya están en la rama; siguen pendientes
la línea base con `go run` y la confirmación de revisión cruzada por ambas
personas. No se atribuye esa revisión a un integrante sin su confirmación.

## Estado actual comprobado

- **Rama local:** `taller11`, creada desde `main`.
- **Rama remota:** `taller11` está publicada en `origin`; el PR #2 se fusionó
  en `main` (commit `814e194`). Los commits por punto y reto están en el
  historial de la rama.
- **Repositorio antes del taller:** no tenía `Dockerfile`, `.dockerignore` ni
  evidencias de `docs/taller11/`; sí contenía el enunciado y la guía de pareja.
- **Go:** `go test ./... -count=1` y `go build ./...` terminaron correctamente
  antes de iniciar la contenerización.
- **Docker Desktop / Engine:** inicialmente no estaba activo; se inició con
  Docker Desktop CLI. `docker info` ya responde y `docker run --rm hello-world`
  terminó con el mensaje esperado.
- **Docker Compose:** versión `v5.5.0` disponible.
- **Reto Compose:** `compose.yaml` validó su configuración y se construyó e
  inició con una base aislada nueva. PostgreSQL quedó `healthy`, la API inició
  en `8080`, ejecutó las migraciones (4 tablas públicas) y respondió al smoke
  test HTTP. Esta base es un volumen nuevo propio del Compose, independiente de
  `gymcoach-db`. La configuración local `.env.compose` se generó con valores
  locales aleatorios y está ignorada por Git; no se registran sus valores.
- **Imágenes base:** se descargaron `golang:1.27` y
  `gcr.io/distroless/static-debian13:nonroot`.
- **Otros contenedores:** Docker ya tenía servicios en ejecución; no se
  detuvo, borró ni modificó ninguno durante el punto 1.
- **Configuración local:** `.env` existe y está excluido por el `.gitignore`
  existente. No copiar sus valores secretos a este documento, a `compose.yaml`,
  a los commits ni a las evidencias.
- **PostgreSQL GymCoach:** el contenedor existente `gymcoach-db` publica el
  puerto `5432`; estaba detenido y se inició para validar acceso. Una consulta
  estrictamente de solo lectura confirmó que las credenciales de `.env`
  funcionan al conectarse a `gymcoach:5432` con el usuario configurado. No se
  consultaron filas ni se ejecutaron migraciones contra ese contenedor.
- **Línea base HTTP con `go run`:** pendiente. El 12 de noviembre de 2026,
  Docker Desktop reportó estado `stopped` y no se pudo verificar una base
  temporal aislada. No se ejecutó la API contra la base local habitual, porque
  el arranque aplica migraciones y datos de ejemplo. Ver
  [peticiones.md](./peticiones.md).
- **Código Go:** no se modificó; el taller pide empaquetar la API existente,
  no cambiar su implementación.

## Punto 1 · Imagen · Responsable: Persona 1

| Tarea / evidencia | Estado | Observación |
| --- | --- | --- |
| Crear y publicar rama `taller11` | Hecho | La rama remota existe en `origin`. |
| Revisar que `main.go` está junto a `go.mod` | Hecho | El comando de compilación de la imagen puede usar `go build .`. |
| Crear Dockerfile multi-etapa | Hecho y construido | Creado en la raíz: `golang:1.27`, `CGO_ENABLED=0`, etapa final Distroless `nonroot`; ambas construcciones finalizaron correctamente. |
| Crear `.dockerignore` | Hecho y comprobado | Excluye `.env`, `.env.*`, `.git` y artefactos; conserva `.env.example`. `git check-ignore` confirmó `.env` y `.env.docker` excluidos. |
| `docker run --rm hello-world` | Hecho | Docker Desktop mostró `Hello from Docker!`. |
| Descargar imágenes base | Hecho | `golang:1.27` y `gcr.io/distroless/static-debian13:nonroot` descargadas. |
| Construir imagen final y etapa `compilacion` | Hecho | `docker build -t gymcoach .` y `docker build --target compilacion -t gymcoach:compilacion .` finalizaron correctamente. |
| Guardar `imagenes.txt` y `capas.txt` | Hecho | Archivos generados desde `docker image ls gymcoach` y `docker history gymcoach`. |
| Guardar `contexto.txt` y `usuario.txt` | Hecho | `/src` contiene `.env.example`, no contiene `.env` ni `.git`; la imagen final indica el usuario `65532`. |
| Revisión de Persona 2 y commit del punto 1 | Parcial | El commit `9456fcb` está publicado; la revisión de Persona 2 no está acreditada. |

### Resultados medidos del punto 1

- Imagen final `gymcoach:latest`: **40,4 MB de uso en disco**, **12,2 MB de
  contenido**.
- Etapa `gymcoach:compilacion`: **1,63 GB de uso en disco**, **388 MB de
  contenido**.
- `docker build` reportó un contexto de aproximadamente **624,56 kB**; el
  listado `/src` confirma que el contexto no incluyó `.env` ni `.git`.
- La plantilla `.env.example` sí quedó en el contexto; la imagen final solo
  copia el binario desde la etapa de compilación.
- Evidencia conservada: `imagenes.txt`, `capas.txt`, `contexto.txt` y
  `usuario.txt` en esta carpeta.

### Criterio para marcar el punto 1 como terminado

- [x] `Dockerfile` y `.dockerignore` creados; pendiente la revisión de Persona 2.
- [x] Las construcciones final y de `compilacion` terminan correctamente.
- [x] Imagen final de 40,4 MB y etapa de compilación de 1,63 GB; tamaños
  registrados desde Docker.
- [x] `imagenes.txt`, `capas.txt`, `contexto.txt` y `usuario.txt` contienen la
  salida real.
- [x] `usuario.txt` indica `65532`; el contexto no contiene `.env` ni `.git`.
- [x] El commit del punto 1 (`9456fcb`) está publicado en `taller11`.
- [ ] Persona 2 confirmó la revisión del cambio.

## Punto 2 · Contenedor atiende la API · Responsable: Persona 2

| Tarea / evidencia | Estado | Observación |
| --- | --- | --- |
| Crear `.env.docker` local | Hecho para la prueba | Se usó localmente y no se versionó; la evidencia no contiene sus valores. |
| Probar imagen sin variables (`1_sin_variables.txt`) | Hecho | Registra los errores esperados por configuración faltante. |
| Probar conexión con variables (`2_local.txt`) | Hecho | Registra el inicio del contenedor y la conexión a PostgreSQL temporal. |
| Arrancar API en contenedor con puerto publicado | Hecho | La evidencia registra la API en `18081`, con el puerto interno alineado con `PUERTO`. |
| Ejecutar las seis peticiones de Bruno | Hecho en contenedor | Se registran `401, 200, 200, 403, 200, 422`; falta compararlos con `go run`. |
| Guardar `peticiones.md` y `logs_contenedor.txt` | Hecho | Ambos archivos están versionados y los logs no incluyen tokens ni contraseñas. |
| Comparar respuesta completa de petición 3 | Parcial | La respuesta del contenedor está guardada; falta la respuesta de `go run` para comparar. |
| Escribir D1–D4 en `docs/decisiones.md` | Hecho | Las decisiones están documentadas en la sección del contenedor. |
| Revisión de Persona 1 y commit del punto 2 | Parcial | El commit `a0a2c32` está publicado; falta acreditar la revisión cruzada. |

## Reto · API y PostgreSQL con Compose · Trabajo compartido

### Avance preparado por Persona 1 para revisión de Persona 2

- Se creó `compose.yaml` con servicios `api` y `postgres`; la API usa el DNS
  interno `postgres:5432` y la base no publica puertos en el host.
- Se añadió healthcheck `pg_isready`; Compose espera `service_healthy` antes
  de iniciar la API, que luego ejecuta sus migraciones y datos de ejemplo.
- Se añadió el volumen nombrado `postgres_data` dentro del proyecto Compose
  `gymcoach-taller11`; es independiente del contenedor PostgreSQL que ya ocupa
  `gymcoach-db` del host.
- Se creó `.env.compose.example`. La configuración local `.env.compose` no debe
  versionarse; `.gitignore` y `.dockerignore` permiten la plantilla y excluyen
  la configuración local.
- El puerto de la API se configura desde `PUERTO` (por defecto `8080`). No se
  exponen secretos en el YAML.

**Estado:** hay evidencia de tres arranques desde volúmenes vacíos, persistencia,
regresión HTTP y decisiones D5–D9. El commit `eb11f78` está publicado. La
revisión cruzada sigue pendiente de confirmación por las personas.

| Tarea / evidencia | Estado | Observación |
| --- | --- | --- |
| Definir nombre del servicio, DNS interno y puerto PostgreSQL | Hecho; revisión pendiente | Compose usa `postgres:5432`, no `localhost`, `host.docker.internal` ni un puerto publicado del host. |
| Diseñar disponibilidad de PostgreSQL antes de iniciar API | Hecho y probado 3/3 | Healthcheck `pg_isready` y dependencia `service_healthy`; evidencia en `compose.txt`. |
| Diseñar variables/secretos sin valores sensibles en el YAML | Hecho | `.env.compose.example` es plantilla; los valores locales no están versionados. |
| Diseñar volumen persistente y aislamiento del PostgreSQL existente | Hecho y probado | Los conteos persistieron tras `down`/`up`; la base no publica puertos al host. |
| Implementar y validar sintaxis de `compose.yaml` | Hecho | La evidencia registra que `docker compose ... config --quiet` pasó. |
| Probar desde DB vacía tres veces seguidas | Hecho (3/3) | PostgreSQL llegó a `healthy`, la API inició y se encontraron cuatro tablas en cada ciclo. |
| Confirmar el arranque integrado y migraciones | Hecho | API arriba, PostgreSQL saludable y cuatro tablas; evidencia en `compose.txt`. |
| Smoke test HTTP de disponibilidad | Hecho | `GET /login` devolvió `405`, esperado para una ruta que acepta `POST`. |
| Comprobar persistencia después de `docker compose down` | Hecho | Se conservaron los conteos `1, 2, 2, 2` de las cuatro tablas registradas. |
| Ejecutar seis peticiones con Compose | Hecho | `peticiones.md` registra `401, 200, 200, 403, 200, 422`. |
| Guardar `compose.txt` y completar D5–D9 | Hecho | La evidencia incluye arranques, persistencia, regresión HTTP y decisiones. |
| Revisión de ambas personas y commit del reto | Parcial | El commit `eb11f78` está publicado; la revisión cruzada y su confirmación siguen pendientes. |

## Entrega final

- [ ] Docker Desktop está iniciado y las imágenes base están disponibles.
- [x] La rama `taller11` está publicada en el remoto correcto.
- [x] Hay commits publicados para el punto 1, punto 2 y reto; cada uno declara
  qué herramienta de IA se usó y para qué, de acuerdo con el enunciado.
- [ ] Las evidencias de línea base, imagen, contenedor y Compose están completas:
  falta la línea base con `go run`.
- [x] D1–D9 están en `docs/decisiones.md` y describen la solución implementada.
- [ ] `go test ./...` y `go build ./...` siguen pasando.
- [ ] `git status` no muestra `.env`, `.env.docker` ni otros secretos como
  archivos por agregar.
- [ ] Ambas personas confirmaron la revisión cruzada de los cambios.
- [ ] Todo quedó subido a la rama `taller11` dentro de la fecha de entrega.
