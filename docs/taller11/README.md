# Taller 11 · GymCoach en contenedores

Guía de trabajo en pareja para trasladar el enunciado
[S11_06_Taller.docx](./S11_06_Taller.docx) al proyecto GymCoach. El objetivo es
empaquetar la API existente, demostrar su comportamiento en Docker y resolver
el reto de ejecutarla junto con PostgreSQL mediante Docker Compose.

El avance y las tareas que aún no se pueden marcar como terminadas se registran
en [puntos_realizados.md](./puntos_realizados.md).

El enunciado fija la entrega en la rama `taller11` del repositorio de la pareja,
antes del jueves 12 de noviembre a las 23:59. Confirmen la fecha y hora con el
docente si recibieron una actualización.

> **Alcance del taller:** no se debe cambiar código Go para hacer que la API
> arranque en el contenedor. Primero revisen configuración, variables, rutas,
> puertos, imagen y red. La publicación de la imagen en Docker Hub, Kubernetes,
> CORS y seguridad avanzada están fuera de este taller.

## 1. Resultado esperado

Al finalizar, la rama `taller11` debe contener:

- Un `Dockerfile` de dos etapas: compilación con Go y ejecución mínima con
  `gcr.io/distroless/static-debian13:nonroot`.
- Un `.dockerignore` que excluya secretos, `.git` y artefactos locales, pero
  permita `.env.example`.
- Evidencia de que la imagen final es pequeña, el contexto no contiene `.env`
  ni `.git`, y el proceso de la imagen final usa el usuario no privilegiado
  `65532`.
- Evidencia de que el contenedor de GymCoach recibe la configuración al
  arrancar, se conecta a PostgreSQL y mantiene los resultados HTTP del Hito 2.
- Un `compose.yaml` que permita iniciar API y PostgreSQL juntos, con un comando,
  sin `host.docker.internal` ni el puerto publicado de PostgreSQL de la máquina.
- Evidencias y decisiones D1–D9 en la ubicación y el formato pedidos por el
  enunciado.

La revisión es sobre los archivos subidos al repositorio; no se debe depender de
que el revisor ejecute los comandos. Por eso las salidas solicitadas se guardan
en `docs/taller11/`.

## 2. Organización de la pareja

Completen los nombres antes de empezar. Cada persona debe ejecutar y entender
su parte; la revisión cruzada no reemplaza la participación individual.

| Responsable | Enfoque principal | Responsabilidades |
| --- | --- | --- |
| **Persona 1: __________________** | Imagen y compilación reproducible | Lidera el punto 1 (`Dockerfile`, `.dockerignore`, construcción y evidencias de tamaño/contexto/usuario); aporta el diseño de red y persistencia del Compose; revisa la configuración de Persona 2. |
| **Persona 2: __________________** | Configuración y comportamiento de la API | Lidera el punto 2 (`.env.docker`, ejecución del contenedor, seis peticiones y logs); implementa junto con Persona 1 la configuración API–PostgreSQL de Compose; revisa la imagen de Persona 1. |
| **Ambas personas** | Línea base, reto e integración | Preparan Docker y la rama, mantienen las decisiones, revisan evidencias, prueban desde una base vacía, verifican persistencia y confirman que ningún secreto quede versionado. |

### Regla de colaboración

La persona responsable prepara el cambio y su evidencia; la otra lo revisa
leyendo el archivo, comprobando el criterio de aceptación y dejando cualquier
corrección antes del commit. En Compose, Persona 1 puede liderar red, puertos y
volumen; Persona 2, variables, disponibilidad de PostgreSQL y arranque de la
API. Ambas deben poder explicar el archivo completo.

El enunciado pide **un commit por cada punto y otro para el reto** y considera
los commits de cada integrante. Organicen los autores para que ambos tengan
commits identificables en `taller11` (por ejemplo: Persona 1 firma el punto 1,
Persona 2 el punto 2 y Persona 2 lidera el commit del reto, con revisión y
participación real de Persona 1). No atribuyan a una persona un cambio que no
haya realizado. Cada commit debe incluir en su cuerpo la declaración de uso de
IA y qué parte se apoyó con ella.

## 3. Preparación compartida · Fase 0

Completen esta fase antes de trabajar en paralelo:

- [ ] Confirmar que `go test ./...` y `go build ./...` pasan.
- [ ] Iniciar la API con `go run .` y confirmar que las seis peticiones de
  Bruno del Hito 2 funcionan con PostgreSQL.
- [ ] Tener Docker Desktop iniciado y comprobar `docker run --rm hello-world`.
- [ ] Descargar las imágenes base del enunciado: `golang:1.27` y
  `gcr.io/distroless/static-debian13:nonroot`.
- [ ] Partir de la rama principal actualizada, crear `taller11` y subirla al
  remoto.
- [ ] Hacer una construcción de prueba con el Dockerfile del kit, adaptando el
  destino de `go build` solo si `main.go` no está junto a `go.mod`.
- [ ] Crear `docs/taller11/peticiones.md` desde la tabla del taller. Registrar
  la columna **Con go run** con respuestas reales y añadir debajo la respuesta
  completa de la petición 3 para compararla en el paso 2.5.

Las seis peticiones de referencia son: sin token (401), login de Ana (200),
consulta autenticada (200), cambio por Luis/entrenador (403), transición válida
por Ana (200) y transición prohibida de vuelta a `pendiente` (422). La petición
de preparación crea una rutina nueva; la transición válida consume ese estado,
así que creen otra rutina al repetir la secuencia.

En Windows, resuelvan primero Docker Desktop/WSL 2 y los puertos ocupados. Si
una solución depende de Windows o difiere del macOS usado para probar el
enunciado, registren el fallo, la causa y la solución en
`docs/taller11/windows.md`. No cambien código Go para ocultar un problema de
Docker.

## 4. Punto 1 · Construir y revisar la imagen

**Lidera:** Persona 1. **Revisa:** Persona 2.

### Trabajo

1. Crear el `Dockerfile` en la raíz y explicar antes de avanzar qué hace cada
   instrucción. Debe conservar la estructura del kit:
   - etapa de compilación con `golang:1.27`;
   - copia de `go.mod` y `go.sum`, descarga de módulos y copia del proyecto;
   - binario compilado con `CGO_ENABLED=0`;
   - etapa final `gcr.io/distroless/static-debian13:nonroot`, con solo el
     ejecutable y `ENTRYPOINT` al binario;
   - ruta de compilación correspondiente a la ubicación real del `main.go`
     (en GymCoach, está en la raíz junto a `go.mod`).
2. Crear `.dockerignore` en la raíz. Debe excluir `.env`, `.env.*`, `.git`, los
   ejecutables y otros artefactos locales, conservando explícitamente
   `.env.example`. No incluyan secretos en la imagen ni en el contexto de build.
3. Construir la imagen de producción y la etapa de compilación por separado.
   Adaptar el nombre de ejemplo `canchas` a una etiqueta minúscula, por ejemplo
   `gymcoach`.

   Desde PowerShell, en la raíz del proyecto:

   ```powershell
   docker build -t gymcoach .
   docker build --target compilacion -t gymcoach:compilacion .
   docker image ls gymcoach
   docker history gymcoach
   ```

4. Guardar las evidencias del enunciado en `docs/taller11/`:
   - `imagenes.txt`: `docker image ls` de ambas etiquetas.
   - `capas.txt`: `docker history` de la imagen final.
   - `contexto.txt`: listado de `/src` dentro de la etapa de compilación.
   - `usuario.txt`: usuario configurado en la imagen final.
5. Revisar que la imagen final mida decenas de megabytes, la etapa de
   compilación sea sustancialmente mayor, `contexto.txt` no muestre `.env` ni
   `.git`, y `usuario.txt` muestre `65532`.

### Criterio de aceptación del punto 1

- [ ] `Dockerfile` y `.dockerignore` están en la raíz y coinciden con el kit y
  con la estructura real del proyecto.
- [ ] La construcción de ambas etapas termina correctamente.
- [ ] Las cuatro evidencias están guardadas y corresponden a la construcción
  revisada.
- [ ] El archivo `.env` no se envió en el contexto ni aparece en la imagen.
- [ ] El proceso final no corre como `root`.
- [ ] Persona 2 revisó el cambio; se creó el commit del punto 1.

> D1–D4 pertenecen a la sección del contenedor en `docs/decisiones.md`; se
> completan al cerrar el punto 2, cuando ya están disponibles todas las
> evidencias necesarias.

## 5. Punto 2 · Ejecutar el contenedor y comprobar la API

**Lidera:** Persona 2. **Revisa:** Persona 1.

### Trabajo

1. Arrancar la imagen sin variables y guardar la salida en
   `docs/taller11/1_sin_variables.txt`. Es esperado que la configuración
   detenga el proceso al faltar `DATABASE_URL`.
2. Crear `.env.docker` a partir de la configuración local, pero adaptada al
   contenedor. No subirla. El `.gitignore` actual ya contiene `.env.*` y
   permite `.env.example`; confirmar la exclusión con
   `git check-ignore -v .env.docker`. Si el patrón se modifica, preservar la
   excepción de `.env.example`.
3. Para esta fase, la API se conecta a PostgreSQL que corre fuera del contenedor:
   - no usar `localhost` como host de la base dentro del contenedor;
   - en Docker Desktop para Windows/macOS, usar `host.docker.internal`;
   - conservar el puerto que PostgreSQL publica en la máquina;
   - antes de probar, confirmar que el puerto corresponde a la base
     `gymcoach` y que las credenciales son válidas. El contenedor local
     `gymcoach-db` publica `5432`; no asumir que un puerto alternativo apunta a
     la misma base;
   - en Docker Engine Linux sin Docker Desktop, usar la opción
     `--add-host=host.docker.internal:host-gateway` si se necesita;
   - usar valores sin comillas ni comentarios al final de las líneas para el
     `--env-file` de Docker.
4. Guardar la ejecución con `.env.docker` en `2_local.txt`, distinguiendo un
   fallo de configuración de uno de conectividad, credenciales o nombre de base.
5. Arrancar el contenedor con `--env-file`, publicar el puerto de la API y
   revisar `docker logs`. El puerto del contenedor debe coincidir con `PUERTO`.
   Si el puerto de la máquina está ocupado, mapear otro puerto del host al
   puerto real del contenedor y actualizar Bruno solo para esa prueba.

   Con `PUERTO=8080`, el mapeo normal es `-p 8080:8080`. La prueba manual
   se inicia con `docker run -d --name gymcoach-api -p 8080:8080 --env-file .env.docker gymcoach`.
   Si se elige otro puerto del host, por ejemplo `-p 8081:8080`, Bruno debe
   apuntar a `localhost:8081`; el puerto derecho sigue siendo el del servidor
   dentro del contenedor.
6. Ejecutar contra el contenedor las mismas seis peticiones, en el mismo orden,
   con una rutina recién creada. Completar en `peticiones.md` la columna
   **En el contenedor**; guardar `docker logs` en
   `docs/taller11/logs_contenedor.txt`.
7. Comparar la respuesta completa de la petición 3 con la línea base. Explicar
   D4 en `docs/decisiones.md`: si se observa otro formato horario, verificar si
   representa el mismo instante y relacionar el cambio con la zona horaria
   disponible en la imagen mínima. Si la API no devuelve fechas, indicar por
   qué no aplica al contrato actual.
8. Detener y eliminar solo el contenedor de la API del taller. No borrar el
   contenedor PostgreSQL de la semana 3 ni sus datos.

### Criterio de aceptación del punto 2

- [ ] `1_sin_variables.txt` y `2_local.txt` registran resultados reales y
  explicables.
- [ ] `.env.docker` está ignorado por Git y no aparece en `git status` como
  archivo versionable.
- [ ] `peticiones.md` contiene las seis respuestas de referencia y la respuesta
  completa de la petición 3.
- [ ] Las seis respuestas **En el contenedor** tienen los mismos códigos que
  **Con go run**: 401, 200, 200, 403, 200 y 422.
- [ ] `logs_contenedor.txt` contiene evidencia de las seis peticiones al
  servidor del contenedor, sin tokens ni contraseñas.
- [ ] D1–D4 están escritos en la sección `## Semana 11 · el contenedor` de
  `docs/decisiones.md`.
- [ ] Persona 1 revisó evidencias y decisiones; se creó el commit del punto 2.

## 6. Reto · API y PostgreSQL juntos con Compose

**Co-diseño:** Persona 1 lidera red, volumen y puertos; Persona 2 lidera
configuración de la API y orden de arranque. La persona designada para el commit
integra el archivo y evidencias después de la revisión de ambas.

### Diseño inicial implementado y comprobaciones obligatorias

Se preparó `compose.yaml` en la raíz y una plantilla local
`.env.compose.example`. El diseño propuesto para revisión de ambas personas es:

1. **Red interna:** la API encuentra PostgreSQL por el nombre del servicio de
   Compose y el puerto interno de PostgreSQL (`5432`). No usa
   `host.docker.internal` ni el puerto publicado del host (`5433`) para hablar
   con la base.
2. **Disponibilidad inicial:** con PostgreSQL recién creado y vacío, la API no
   debe intentar migrar/sembrar antes de que la base acepte conexiones. Definan
   cómo Compose espera la salud de la base antes de arrancar la API, usando
   opciones compatibles con la versión instalada de Compose.
3. **Configuración y secretos:** `compose.yaml` se versiona, pero no contiene
   contraseñas ni claves reales. Los valores secretos se reciben por variables
   locales ignoradas por Git; eviten imprimirlos en logs o evidencias.
4. **Persistencia:** un volumen nombrado conserva los datos después de
   `docker compose --env-file .env.compose down`;
   `docker compose --env-file .env.compose down -v` elimina deliberadamente
   los datos del volumen del proyecto.
5. **Puertos y aislamiento:** Bruno llega a la API por el puerto habitual de
   GymCoach. La base de Compose no debe chocar con el contenedor PostgreSQL de
   la semana 3. No publiquen PostgreSQL en el host si no es necesario; si no se
   publica, expliquen cómo consultarla con
   `docker compose --env-file .env.compose exec ... psql`.
6. **Prueba desde cero:** ejecutar
   `docker compose --env-file .env.compose down -v` solo para el proyecto
   Compose del taller, luego
   `docker compose --env-file .env.compose up -d --build`, y repetir hasta
   obtener tres arranques correctos seguidos desde bases vacías. Confirmar
   aparte que `docker compose --env-file .env.compose down` sin `-v` conserva
   los registros.
7. **Regresión:** correr las seis peticiones de Bruno con ambos servicios arriba
   y completar la columna **Con Compose** en `peticiones.md`.
8. **Evidencia:** guardar
   `docker compose --env-file .env.compose ps` y
   `docker compose --env-file .env.compose logs api` en
   `docs/taller11/compose.txt`.
9. **Decisiones:** completar D5–D9 bajo
   `## Semana 11 · reto: la API y la base, juntas` en `docs/decisiones.md`.

Para probar esta configuración desde PowerShell, copia la plantilla a
`.env.compose`, cambia la contraseña de PostgreSQL y la clave de firma por
valores locales, y conserva la contraseña de PostgreSQL en formato URL-safe
(letras, números, guiones o guiones bajos):

```powershell
if (-not (Test-Path .env.compose)) {
    Copy-Item .env.compose.example .env.compose
}
docker compose --env-file .env.compose config --quiet
docker compose --env-file .env.compose up -d --build
docker compose --env-file .env.compose ps
docker compose --env-file .env.compose logs api
```

`.env.compose` está excluido de Git y del contexto de construcción. El servicio
PostgreSQL de Compose no publica puertos en Windows y usa un volumen aislado
del proyecto `gymcoach-taller11`; no apunta al contenedor local `gymcoach-db`.
Para inspeccionar la base aislada del Compose se puede usar
`docker compose --env-file .env.compose exec postgres psql -U gymcoach -d gymcoach`.
`docker compose --env-file .env.compose down` conserva el volumen; añadir
`-v` lo elimina. Ejecuten ese comando solo desde esta carpeta y cuando acepten
perder los datos de prueba de este proyecto.

El diseño y la configuración inicial no completan por sí solos el reto: faltan
las seis peticiones HTTP, la revisión cruzada y el commit. Las pruebas desde
bases vacías y persistencia ya se ejecutaron localmente y están registradas en
`compose.txt`; Persona 2 debe revisar esos resultados.

### Criterio de aceptación del reto

- [x] `docker compose --env-file .env.compose up -d --build` levanta API y
  PostgreSQL con un solo comando.
- [x] La API usa el nombre DNS del servicio y el puerto interno de PostgreSQL;
  no depende de la máquina del desarrollador.
- [x] Una base nueva y vacía no provoca que la API se caiga al arrancar.
- [x] Tres ciclos desde cero terminan correctamente; la persistencia tras
  `down` sin `-v` está comprobada.
- [x] El volumen no interfiere con el PostgreSQL local existente: no se publica
  el puerto de la base y el volumen pertenece al proyecto Compose del taller.
- [x] `compose.yaml` no contiene secretos reales ni configuraciones locales
  privadas.
- [ ] Las seis respuestas **Con Compose** son 401, 200, 200, 403, 200 y 422.
- [ ] `peticiones.md` incluye la regresión y las seis respuestas; falta completar
  la evidencia HTTP junto a `compose.txt` y confirmar D5–D9 en la revisión.
- [ ] Ambas personas revisaron el reto; se creó su commit.

## 7. Estructura de entrega y evidencias

Al completar el taller, revisar que la raíz y los documentos queden así:

```text
Dockerfile
.dockerignore
compose.yaml
docs/
  decisiones.md                 # D1–D9, sin reemplazar decisiones anteriores
  taller11/
    README.md                   # esta guía de organización
    peticiones.md               # tabla comparativa y cuerpo de petición 3
    imagenes.txt
    capas.txt
    contexto.txt
    usuario.txt
    1_sin_variables.txt
    2_local.txt
    logs_contenedor.txt
    compose.txt
    windows.md                  # solo si hubo diferencias/errores en Windows
```

`.env` y `.env.docker` son archivos locales y nunca forman parte de la entrega.
Antes de cada commit, inspeccionen `git status`, ignoren los archivos locales,
revisen que los logs no contengan secretos y confirmen que solo agregan
evidencias auténticas. No inventen salidas para completar la tabla.

## 8. Plan de commits, revisión y cierre

El enunciado requiere un commit por punto y otro por el reto. Propuesta de
secuencia:

1. **Punto 1 — Persona 1:** Dockerfile, `.dockerignore` y evidencias de imagen.
   Persona 2 revisa tamaño, contexto y usuario.
2. **Punto 2 — Persona 2:** tabla de peticiones del contenedor, logs y D1–D4.
   Persona 1 verifica la regresión y que no se subió `.env.docker`.
3. **Reto — Persona 2 (o quien necesite completar sus commits):** `compose.yaml`,
   columna Compose, `compose.txt` y D5–D9. La otra persona revisa red, salud,
   volumen, secretos y los tres arranques desde cero.

Cada commit debe decir qué punto cierra y, en el cuerpo, declarar la herramienta
de IA utilizada y para qué se usó. No incluyan una afirmación que ninguno de los
dos pueda defender en la revisión o en el hito 3.

### Revisión final de la pareja

- [ ] Ambos entienden el Dockerfile, `.dockerignore` y el flujo de variables.
- [ ] Ambos pueden explicar por qué `localhost` cambia de significado dentro
  de un contenedor y por qué Compose usa un nombre de servicio.
- [ ] Las decisiones D1–D9 describen la configuración real, no una solución
  hipotética.
- [ ] `go test ./...` y `go build ./...` siguen pasando; el taller no modificó
  código Go.
- [ ] Los resultados guardados corresponden a la revisión final de `taller11`.
- [ ] `git status` no revela archivos secretos pendientes.
- [ ] Los tres commits requeridos están subidos a la rama `taller11` del
  repositorio de la pareja antes del cierre indicado en el enunciado.

## 9. Nota para Windows

El enunciado indica que el taller se comprobó en macOS, no en Windows. En
Windows, distingan errores de Docker Desktop/WSL 2, puertos ya ocupados,
formato de `.env.docker` y conectividad a PostgreSQL. Registren en
`windows.md` únicamente lo que les haya ocurrido:

| Qué falló | Mensaje observado | Causa comprobada | Solución aplicada | Evidencia |
| --- | --- | --- | --- | --- |
| Completar si aplica | | | | |

No escriban contraseñas, tokens ni el contenido de `.env` en este registro.
