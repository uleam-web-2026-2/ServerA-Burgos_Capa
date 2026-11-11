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

El trabajo de la Persona 2 se registra aquí como pendiente de coordinación; no
se da por realizado por Persona 1. Ambas personas deben revisar y comprender el
`compose.yaml`, las decisiones D5–D9 y las pruebas integradas.

## Estado actual comprobado

- **Rama local:** `taller11`, creada desde `main`.
- **Rama remota:** aún no creada/publicada. Hay que subirla al completar los
  primeros cambios y seguir publicando los commits de cada punto.
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
- **Configuración local / línea base HTTP:** en la última comprobación, `.env`
  y `.env.example` indican `localhost:5433`, aunque `gymcoach-db` publica
  `5432`. La conexión de solo lectura se comprobó directamente en `5432`, sin
  cambiar ninguno de esos archivos. Antes de ejecutar `go run .` o el punto 2,
  confirmar/corregir localmente el puerto y volver a validar la conexión; las
  seis peticiones de línea base siguen pendientes.
- **Código Go:** no se modificó; el taller pide empaquetar la API existente,
  no cambiar su implementación.

## Punto 1 · Imagen · Responsable: Persona 1

| Tarea / evidencia | Estado | Observación |
| --- | --- | --- |
| Crear rama local `taller11` | Hecho localmente | Falta publicar la rama en `origin`. |
| Revisar que `main.go` está junto a `go.mod` | Hecho | El comando de compilación de la imagen puede usar `go build .`. |
| Crear Dockerfile multi-etapa | Hecho y construido | Creado en la raíz: `golang:1.27`, `CGO_ENABLED=0`, etapa final Distroless `nonroot`; ambas construcciones finalizaron correctamente. |
| Crear `.dockerignore` | Hecho y comprobado | Excluye `.env`, `.env.*`, `.git` y artefactos; conserva `.env.example`. `git check-ignore` confirmó `.env` y `.env.docker` excluidos. |
| `docker run --rm hello-world` | Hecho | Docker Desktop mostró `Hello from Docker!`. |
| Descargar imágenes base | Hecho | `golang:1.27` y `gcr.io/distroless/static-debian13:nonroot` descargadas. |
| Construir imagen final y etapa `compilacion` | Hecho | `docker build -t gymcoach .` y `docker build --target compilacion -t gymcoach:compilacion .` finalizaron correctamente. |
| Guardar `imagenes.txt` y `capas.txt` | Hecho | Archivos generados desde `docker image ls gymcoach` y `docker history gymcoach`. |
| Guardar `contexto.txt` y `usuario.txt` | Hecho | `/src` contiene `.env.example`, no contiene `.env` ni `.git`; la imagen final indica el usuario `65532`. |
| Revisión de Persona 2 y commit del punto 1 | Pendiente | Falta revisión cruzada, commit del punto y publicación de `taller11` en `origin`. |

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
- [ ] Persona 2 revisó el cambio; se creó el commit del punto 1 y se publicó la
  rama.

## Punto 2 · Contenedor atiende la API · Responsable: Persona 2

| Tarea / evidencia | Estado | Observación |
| --- | --- | --- |
| Crear `.env.docker` local | Pendiente | No subirlo; adaptar solo el host de PostgreSQL para acceso desde contenedor y ocultar secretos en evidencias. |
| Probar imagen sin variables (`1_sin_variables.txt`) | Pendiente | Debe registrar el error de configuración esperado, sin secretos. |
| Probar conexión con variables (`2_local.txt`) | Pendiente | Requiere Docker activo y PostgreSQL disponible; distinguir error de host, puerto, credenciales o base. |
| Arrancar API en contenedor con puerto publicado | Pendiente | El puerto de destino debe coincidir con `PUERTO`. |
| Ejecutar las seis peticiones de Bruno | Pendiente | Comparar códigos con la columna `Con go run`; usar una rutina nueva cuando se repite la transición. |
| Guardar `peticiones.md` y `logs_contenedor.txt` | Pendiente | No incluir tokens, contraseñas ni datos locales sensibles. |
| Comparar respuesta completa de petición 3 | Pendiente | Usar la respuesta real de línea base para analizar la diferencia horaria, si existe. |
| Escribir D1–D4 en `docs/decisiones.md` | Pendiente | Explicaciones basadas en evidencias, no en ejemplos del enunciado. |
| Revisión de Persona 1 y commit del punto 2 | Pendiente | La revisión verifica regresión, logs y ausencia de `.env.docker` versionado. |

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

**Estado:** implementación local validada técnicamente: pasaron tres arranques
desde volúmenes vacíos y se comprobó que los datos sobreviven a `down` sin
`-v`. Persona 2 aún debe revisar y participar en la regresión HTTP; el reto no
se marca como formalmente terminado hasta completar peticiones, evidencias,
revisión cruzada y commit.

| Tarea / evidencia | Estado | Observación |
| --- | --- | --- |
| Definir nombre del servicio, DNS interno y puerto PostgreSQL | Preparado; revisión de Persona 2 pendiente | `postgres:5432` en la red de Compose; no `localhost`, `host.docker.internal` ni el puerto publicado `5433`. |
| Diseñar disponibilidad de PostgreSQL antes de iniciar API | Hecho y probado 3/3 | Healthcheck `pg_isready` y dependencia `service_healthy`; los tres arranques desde volúmenes vacíos pasaron. |
| Diseñar variables/secretos sin valores sensibles en el YAML | Preparado; revisión pendiente | `.env.compose.example` es plantilla; los valores reales van en `.env.compose`, ignorado por Git y Docker. |
| Diseñar volumen persistente y aislamiento del PostgreSQL existente | Hecho y probado localmente | Volumen del proyecto `gymcoach-taller11`; se mantuvo el PostgreSQL host ajeno en `5432`, sin publicar el puerto DB de Compose; los conteos persistieron tras `down`/`up`. |
| Implementar y validar sintaxis de `compose.yaml` | Hecho localmente; revisión pendiente | `docker compose --env-file .env.compose config --quiet` pasó; faltan revisión de Persona 2 y commit compartido. |
| Probar desde DB vacía tres veces seguidas | Hecho localmente (3/3); revisión pendiente | En cada ciclo, PostgreSQL llegó a `healthy`, la API inició y se encontraron cuatro tablas migradas. Se usó `down -v` solo con el proyecto `gymcoach-taller11`. |
| Confirmar el arranque integrado y migraciones | Hecho localmente | API escuchó en `8080`; se crearon cuatro tablas públicas; PostgreSQL quedó `healthy`. Resultado registrado en `compose.txt`. |
| Smoke test HTTP de disponibilidad | Hecho; no sustituye la regresión | `GET /login` devolvió `405` como corresponde a una ruta que acepta `POST`; aún faltan las seis peticiones de Bruno. |
| Comprobar persistencia después de `docker compose down` | Hecho localmente | Después de `down` y `up`, se conservaron los conteos de entrenadores, clientes, rutinas y usuarios (1, 2, 2, 2). |
| Ejecutar seis peticiones con Compose | Pendiente | Completar columna `Con Compose` con códigos reales 401, 200, 200, 403, 200 y 422. |
| Guardar `compose.txt` y completar D5–D9 | Parcial | El archivo contiene la configuración, los tres arranques limpios, persistencia, estado y logs; faltan documentar/registrar las seis peticiones HTTP. |
| Revisión de ambas personas y commit del reto | Pendiente | Persona 2 revisa; ambas deben poder explicar red, salud, volumen, secretos y puertos. |

## Entrega final

- [ ] Docker Desktop está iniciado y las imágenes base están disponibles.
- [ ] La rama `taller11` está publicada en el remoto correcto.
- [ ] Los commits del punto 1, punto 2 y reto están publicados; cada uno declara
  qué herramienta de IA se usó y para qué, de acuerdo con el enunciado.
- [ ] Las evidencias de línea base, imagen, contenedor y Compose son reales,
  actuales y no contienen secretos.
- [ ] D1–D9 están en `docs/decisiones.md` y describen la solución implementada.
- [ ] `go test ./...` y `go build ./...` siguen pasando.
- [ ] `git status` no muestra `.env`, `.env.docker` ni otros secretos como
  archivos por agregar.
- [ ] La otra persona revisó los cambios de Persona 1 y Persona 1 revisó los de
  Persona 2 antes de cerrar cada punto.
- [ ] Todo quedó subido a la rama `taller11` dentro de la fecha de entrega.
