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

La Persona 2 ya cuenta con evidencia de inicio, conectividad y regresión HTTP
en contenedor. La línea base actual con `go run` y la revisión cruzada de ambas
personas siguen pendientes; ambas deben revisar y comprender el `compose.yaml`,
las decisiones D5–D9 y las pruebas integradas.

## Estado actual comprobado

- **Rama local:** `taller11`, basada en la rama remota de entrega.
- **Rama remota:** `origin/taller11` ya está publicada; el commit inicial del
  punto 1 es `9456fcb`. Los commits de Persona 2 y de la regresión se agregan
  después de guardar y revisar estas evidencias.
- **Repositorio antes del taller:** no tenía `Dockerfile`, `.dockerignore` ni
  evidencias de `docs/taller11/`; sí contenía el enunciado y la guía de pareja.
- **Go:** `go test ./... -count=1` y `go build ./...` terminaron correctamente
  antes de iniciar la contenerización.
- **Docker Desktop / Engine:** inicialmente no estaba activo; se inició con
  Docker Desktop CLI. `docker info` ya responde y `docker run --rm hello-world`
  terminó con el mensaje esperado.
- **Docker Compose:** versión `v5.5.0` disponible.
- **Reto Compose:** `compose.yaml` se construyó e inició con una base aislada.
  PostgreSQL quedó `healthy`; las seis peticiones HTTP pasaron con la API en
  `18080`. Esta validación usó el proyecto aislado
  `gymcoach-taller11-validation`; no tocó bases ni contenedores ajenos.
  `.env.compose` contiene valores aleatorios solo para la prueba y está
  ignorado por Git.
- **Imágenes base:** se descargaron `golang:1.27` y
  `gcr.io/distroless/static-debian13:nonroot`.
- **Otros contenedores:** la validación usó contenedores con nombres únicos.
  No se detuvieron, borraron ni modificaron contenedores ajenos.
- **Configuración local:** `.env` existe y está excluido por el `.gitignore`
  existente. No copiar sus valores secretos a este documento, a `compose.yaml`,
  a los commits ni a las evidencias.
- **PostgreSQL GymCoach:** el contenedor existente `gymcoach-db` publica el
  puerto `5432`; estaba detenido y se inició para validar acceso. Una consulta
  estrictamente de solo lectura confirmó que las credenciales de `.env`
  funcionan al conectarse a `gymcoach:5432` con el usuario configurado. No se
  consultaron filas ni se ejecutaron migraciones contra ese contenedor.
- **Puerto de host:** `8080` estaba ocupado por un proceso PostgreSQL de
  Windows; la prueba con ese puerto recibió un reset. Se usaron `18080` y
  `18081`, sin detener ni cambiar el proceso existente. Ver `windows.md`.
- **Base de prueba:** se usó PostgreSQL temporal y aislado en el puerto host
  `15432` para probar `host.docker.internal`; no se usó la base de datos real
  del equipo. La comparación actual de seis peticiones con `go run` sigue
  pendiente.
- **Código Go:** no se modificó; el taller pide empaquetar la API existente,
  no cambiar su implementación.

## Punto 1 · Imagen · Responsable: Persona 1

| Tarea / evidencia | Estado | Observación |
| --- | --- | --- |
| Crear y publicar rama `taller11` | Hecho | `origin/taller11` existe desde el commit `9456fcb`. |
| Revisar que `main.go` está junto a `go.mod` | Hecho | El comando de compilación de la imagen puede usar `go build .`. |
| Crear Dockerfile multi-etapa | Hecho y construido | Creado en la raíz: `golang:1.27`, `CGO_ENABLED=0`, etapa final Distroless `nonroot`; ambas construcciones finalizaron correctamente. |
| Crear `.dockerignore` | Hecho y comprobado | Excluye `.env`, `.env.*`, `.git` y artefactos; conserva `.env.example`. `git check-ignore` confirmó `.env` y `.env.docker` excluidos. |
| `docker run --rm hello-world` | Hecho | Docker Desktop mostró `Hello from Docker!`. |
| Descargar imágenes base | Hecho | `golang:1.27` y `gcr.io/distroless/static-debian13:nonroot` descargadas. |
| Construir imagen final y etapa `compilacion` | Hecho | `docker build -t gymcoach .` y `docker build --target compilacion -t gymcoach:compilacion .` finalizaron correctamente. |
| Guardar `imagenes.txt` y `capas.txt` | Hecho | Archivos generados desde `docker image ls gymcoach` y `docker history gymcoach`. |
| Guardar `contexto.txt` y `usuario.txt` | Hecho | `/src` contiene `.env.example`, no contiene `.env` ni `.git`; la imagen final indica el usuario `65532`. |
| Revisión de Persona 2 y commit del punto 1 | Parcial | El commit `9456fcb` está en `origin/taller11`, pero no incluye una declaración de uso de IA; confirmar ese requisito con su autor. La revisión cruzada también sigue pendiente. |

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
| Crear `.env.docker` local | Hecho localmente | Contiene valores aleatorios de prueba, está ignorado por Git y se elimina al cerrar la validación. |
| Probar imagen sin variables (`1_sin_variables.txt`) | Hecho | Sin variables falla por falta de `CLAVE_FIRMA`; con solo la clave de firma falla por falta de `DATABASE_URL`. |
| Probar conexión con variables (`2_local.txt`) | Hecho con base efímera | Se conectó a un PostgreSQL temporal mediante `host.docker.internal:15432`; no se usó ni migró la base real del equipo. |
| Arrancar API en contenedor con puerto publicado | Hecho | API activa en `18081`, configurado en `.env.docker`; `GET /login` devolvió el `405` esperado. |
| Ejecutar las seis peticiones | Hecho en contenedor | Pasaron `401`, `200`, `200`, `403`, `200` y `422`; el estado final y el código de error también se comprobaron. |
| Guardar `peticiones.md` y `logs_contenedor.txt` | Hecho | Evidencias reales guardadas sin tokens ni contraseñas. |
| Comparar respuesta completa de petición 3 | Parcial | La respuesta del contenedor está registrada; no hay fechas en el cuerpo. Falta repetir con `go run` para una comparación directa. |
| Escribir D1–D4 en `docs/decisiones.md` | Hecho | Decisiones basadas en Dockerfile, contexto, UID y respuesta observada. |
| Revisión de Persona 1 y commit del punto 2 | Parcial | El commit `a0a2c32` guarda las evidencias de Persona 2; falta publicarlo y la revisión cruzada sigue pendiente. |

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

**Estado:** pasaron tres arranques desde volúmenes vacíos y se comprobó que los
datos sobreviven a `down` sin `-v`. La regresión HTTP pasó en una validación
aislada. La revisión cruzada sigue pendiente.

| Tarea / evidencia | Estado | Observación |
| --- | --- | --- |
| Definir nombre del servicio, DNS interno y puerto PostgreSQL | Hecho; revisión de pareja pendiente | `postgres:5432` en la red de Compose; la API no usa `localhost`, `host.docker.internal` ni el puerto publicado de la base. |
| Diseñar disponibilidad de PostgreSQL antes de iniciar API | Hecho y probado 3/3 | Healthcheck `pg_isready` y dependencia `service_healthy`; los tres arranques desde volúmenes vacíos pasaron. |
| Diseñar variables/secretos sin valores sensibles en el YAML | Hecho; revisión pendiente | `.env.compose.example` es plantilla; `.env.compose` está ignorado por Git y Docker. |
| Diseñar volumen persistente y aislamiento del PostgreSQL existente | Hecho y probado localmente | Volumen del proyecto `gymcoach-taller11`; se mantuvo el PostgreSQL host ajeno en `5432`, sin publicar el puerto DB de Compose; los conteos persistieron tras `down`/`up`. |
| Implementar y validar sintaxis de `compose.yaml` | Hecho | `docker compose --env-file .env.compose config --quiet` pasó; se construyó y arrancó el stack de validación aislado. |
| Probar desde DB vacía tres veces seguidas | Hecho localmente (3/3); revisión pendiente | En cada ciclo, PostgreSQL llegó a `healthy`, la API inició y se encontraron cuatro tablas migradas. Se usó `down -v` solo con el proyecto `gymcoach-taller11`. |
| Confirmar el arranque integrado y migraciones | Hecho | API escuchó en `18080`; PostgreSQL quedó `healthy`. Resultado registrado en `compose.txt`. |
| Smoke test HTTP de disponibilidad | Hecho | `GET /login` devolvió `405`, porque la ruta acepta `POST`. |
| Comprobar persistencia después de `docker compose down` | Hecho localmente | Después de `down` y `up`, se conservaron los conteos de entrenadores, clientes, rutinas y usuarios (1, 2, 2, 2). |
| Ejecutar seis peticiones con Compose | Hecho | `401`, `200`, `200`, `403`, `200`, `422`; estado `en_progreso` y error `transicion_prohibida` comprobados. |
| Guardar `compose.txt` y completar D5–D9 | Hecho | Se conservaron los resultados de tres arranques, persistencia, estado y regresión HTTP. |
| Revisión de ambas personas y commit del reto | Parcial | La evidencia integrada y D5–D9 están listas para el commit del reto; faltan la revisión cruzada y repetir la línea base con `go run`. |

## Entrega final

- [ ] Docker Desktop está iniciado y las imágenes base están disponibles.
- [x] La rama `taller11` está publicada en el remoto correcto.
- [ ] Los commits del punto 1, punto 2 y reto están publicados; cada uno declara
  qué herramienta de IA se usó y para qué, de acuerdo con el enunciado.
- [ ] Las evidencias de imagen, contenedor y Compose son reales y no contienen
  secretos; falta repetir y guardar la línea base actual de `go run`.
- [ ] D1–D9 están en `docs/decisiones.md` y describen la solución implementada.
- [x] `go test ./...`, `go vet ./...` y `go build ./...` pasan en la rama de entrega.
- [ ] `git status` no muestra `.env`, `.env.docker` ni otros secretos como
  archivos por agregar.
- [ ] La otra persona revisó los cambios de Persona 1 y Persona 1 revisó los de
  Persona 2 antes de cerrar cada punto.
- [ ] Todo quedó subido a la rama `taller11` dentro de la fecha de entrega.
