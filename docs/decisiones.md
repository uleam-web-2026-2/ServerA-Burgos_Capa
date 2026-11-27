# Registro de decisiones de arquitectura

## D1 Modelo Entrenador-Cliente-Rutina

Separamos el entrenador, el cliente y la rutina en entidades relacionadas. Un entrenador administra varios clientes y cada cliente puede tener distintas rutinas a lo largo del tiempo. Esto evita guardar programas completos dentro del perfil del cliente y permite filtrar y cambiar el estado de cada plan de forma independiente.

## D2 Objetivo y datos de salud

El MVP guarda un objetivo general como texto, pero no guarda diagnósticos, lesiones, peso, medidas ni historial médico. Estos datos pueden ser sensibles; se debe validar consentimiento, controles de acceso y política de retención antes de ampliar el modelo.

## D3 Máquina de estados

La rutina empieza `pendiente`, pasa a `en_progreso` cuando el cliente inicia el plan y termina `completada` cuando finaliza el ciclo. Una rutina completada no se reinicia; el entrenador crea una nueva para conservar trazabilidad.

## D4 Permisos frente a implementación

El taller implementa tokens Bearer y autorización por rol en las rutas. La sesión todavía no está vinculada con un entrenador o cliente concreto: no se filtran registros por propietario. Por eso esta API sigue siendo de desarrollo y no debe recibir datos personales reales.

## Semana 7 · D1 · Qué nivel lleva cada ruta

| Ruta | Nivel | Por qué |
| --- | --- | --- |
| `POST /login` | Pública | Es la puerta para obtener el token que protege las demás rutas. |
| `GET /clientes` | Solo entrenador | El entrenador necesita revisar su cartera y las rutinas asociadas. |
| `POST /clientes` | Solo entrenador | El entrenador registra a las personas a quienes asigna planes. |
| `GET /rutinas` | Solo entrenador | El entrenador administra y revisa las rutinas que creó. |
| `POST /rutinas` | Solo entrenador | El entrenador diseña el plan; una rutina nueva nace pendiente. |
| `GET /rutinas/{id}` | Con sesión | Entrenador y cliente necesitan consultar el detalle; el filtro por propietario queda pendiente. |
| `PUT /rutinas/{id}` | Solo entrenador | El entrenador corrige el contenido del plan; el estado no se cambia por PUT. |
| `PATCH /rutinas/{id}` | Solo cliente | El cliente registra su avance al iniciar o terminar la rutina. |
| `DELETE /rutinas/{id}` | Solo entrenador | El entrenador administra y retira los planes que ya no corresponden. |

## Semana 7 · D2 · La ruta de un solo rol

**Decisión:** `PATCH /rutinas/{id}` solo lo puede usar el rol cliente.

**Por qué:** el cliente es quien sigue el plan y puede afirmar si lo inició o completó. El entrenador diseña y corrige la rutina, pero no debe registrar el progreso en nombre del cliente.

**Qué descartamos:** permitir que el entrenador cambie el estado. Se deja fuera del alcance del taller para mantener una sola fuente del progreso; el estado se cambia con `PATCH`, nunca con `PUT`.

## Semana 7 · D3 · La transición prohibida

**Transición:** de `en_progreso` a `pendiente`.

**Por qué el negocio no la permite:** cuando el cliente ya inició el plan, no puede devolverlo a `pendiente`, porque se perdería el avance registrado. Si necesita reiniciar el trabajo, el entrenador crea una rutina nueva.

**Qué responde y dónde se valida:** HTTP `422`, código `transicion_prohibida` y mensaje `de en_progreso no se puede pasar a pendiente`. La regla pura se comprueba en `internal/gymcoach/dominio/estado.go`, función `PuedeTransicionar`; no conoce HTTP ni la base de datos. El servicio de aplicación la usa antes de pedir al repositorio que guarde el cambio. El manejador HTTP solo valida la petición, llama al servicio y traduce sus errores a códigos HTTP.

## Semana 8 · estructura

La entidad principal que demostramos es **Rutina**. Para el cambio de estado elegimos una estructura por capas:

| Carpeta | Responsabilidad |
| --- | --- |
| `internal/gymcoach/dominio/` | Entidades y reglas del negocio. La máquina de estados no depende de HTTP ni de PostgreSQL. |
| `internal/gymcoach/aplicacion/` | Coordina las operaciones de rutina y, para el cambio de estado, valida la transición antes de solicitar el cambio. |
| `internal/gymcoach/repositorio/` | Acceso a PostgreSQL mediante GORM para listar, crear, consultar, actualizar y eliminar rutinas; persiste el cambio de estado sin decidir qué transición permite el negocio. |
| `internal/gymcoach/` | Adaptador HTTP: registra rutas, valida el JSON y convierte los errores de aplicación en respuestas HTTP. Las rutas de clientes mantienen sus manejadores actuales. |

Elegimos esta separación porque cada parte tiene un motivo distinto para cambiar. Si cambia la máquina de estados, se modifica dominio; si cambia PostgreSQL, el repositorio; si cambia el contrato HTTP, el manejador. En `main.go` se conectan las capas. Todas las rutas de Rutina usan el servicio y el repositorio; las rutas de Cliente no forman parte de la entidad principal demostrada.

Para enseñar el recorrido del `PATCH /rutinas/{id}`, abrir `internal/gymcoach/estado.go` (manejador) y `internal/gymcoach/dominio/estado.go` (regla). La escritura en base de datos queda en `internal/gymcoach/repositorio/rutinas.go`.

## Semana 11 · el contenedor

### D1 · Construcción en etapas

La imagen compila la API con `golang:1.27` y ejecuta únicamente el binario
estático en `gcr.io/distroless/static-debian13:nonroot`. Así las herramientas
de compilación no forman parte de la imagen de ejecución. La imagen final
`gymcoach:latest` midió 40,4 MB de uso en disco y 12,2 MB de contenido; la
etapa `gymcoach:compilacion` midió 1,63 GB de uso en disco y 388 MB de
contenido. Por eso la imagen final es decenas de veces menor que el builder.

### D2 · Contexto de construcción

`.dockerignore` excluye `.env`, `.env.*`, `.git` y artefactos locales, y permite
las plantillas públicas `.env.example` y `.env.compose.example`. La configuración
local `.env.docker` se verificó como ignorada; ningún secreto se registró en
los logs o evidencias. Dentro del contenedor, `/src` contiene `.env.example`
pero no `.env` ni `.git`. El servidor escucha en `:<PUERTO>` y con
`PUERTO=8080` no usa `localhost` para decidir el puerto de PostgreSQL:
`localhost` solo aparece dentro de `DATABASE_URL` al probar desde el host.
Para el contenedor separado se usa `host.docker.internal` y el puerto publicado
por el PostgreSQL del host.

### D3 · Usuario de ejecución

La imagen final usa el usuario `nonroot`; `docs/taller11/usuario.txt` registra
el UID `65532`. El proceso no necesita privilegios de root para atender la API.
La imagen final expone el puerto `8080` y publica un puerto del host con
`-p` como `-p 18081:8080`: el lado izquierdo es el puerto accesible desde Bruno
y Windows, el lado derecho es el puerto real del servidor dentro del contenedor.
El mismo principio se usa con `-p 18080:8080` en el entorno de validación
cuando `8080` ya estaba ocupado.

### D4 · Fechas en la respuesta

La respuesta completa de `GET /rutinas/{id}` se guardó en
`docs/taller11/peticiones.md`. No contiene campos de fecha ni hora, así que no
hay una representación temporal que comparar entre la imagen local y
Distroless; el posible efecto de la zona horaria no aplica al contrato actual.
La línea base actual con `go run` sigue pendiente de repetición.

## Semana 11 · reto: la API y la base, juntas

Estas decisiones registran la configuración del Taller 11. El arranque
desde tres bases vacías y la persistencia tras `down` se validaron localmente.
La regresión HTTP de las seis peticiones pasó en la validación aislada descrita
en `docs/taller11/peticiones.md`; falta la revisión cruzada.

### D5 · Red entre API y PostgreSQL

La API se conecta a PostgreSQL mediante el nombre de servicio `postgres` y el
puerto interno `5432`, dentro de la red privada creada por Compose. No usa
`localhost`, `host.docker.internal` ni el puerto publicado `5433` de la máquina.
PostgreSQL no publica un puerto en el host, por lo que no compite con los otros
contenedores locales.

### D6 · Orden de arranque

PostgreSQL define un healthcheck con `pg_isready`. La API depende de que el
servicio esté `healthy` antes de iniciar; así sus migraciones y datos de ejemplo
no intentan ejecutarse mientras el servidor de base de datos aún está
inicializando. Se validó en tres arranques con bases vacías.

### D7 · Configuración y secretos locales

`compose.yaml` no guarda contraseñas ni claves reales. Los valores requeridos
se leen desde `.env.compose`, archivo local ignorado por Git y por el contexto
de Docker; `.env.compose.example` contiene solo valores de reemplazo. La
contraseña configurada debe ser URL-safe porque se incorpora a `DATABASE_URL`.

### D8 · Persistencia aislada

En el Taller 11, el servicio PostgreSQL montaba el volumen nombrado
`postgres_data`, creado dentro del proyecto `gymcoach-taller11`. Para el reto
del Taller 12 se cambió el nombre del proyecto a `gymcoach-taller12` y así se
mantiene aislado el volumen anterior.
`docker compose --env-file .env.compose down` conserva el volumen;
`docker compose --env-file .env.compose down -v` lo elimina intencionalmente.
Este volumen no es el del contenedor PostgreSQL existente fuera de Compose.

### D9 · Puertos y acceso a la base

La API publica el puerto definido por `PUERTO` (por defecto `8080`) y Bruno
debe apuntar a ese puerto. En la validación de Windows, `8080` estaba ocupado
por un proceso PostgreSQL local, por lo que se usó `PUERTO=18080`; la API
respondió desde el host sin cambiar el puerto interno de PostgreSQL. PostgreSQL
queda solo en la red de Compose; para consultarlo se usa
`docker compose --env-file .env.compose exec postgres psql` con el usuario y la
base configurados localmente. La persistencia, los tres arranques limpios y las
seis peticiones HTTP pasaron.

## Semana 12 · la seguridad de la API

### D1 · Lo que encontramos y cerramos

La búsqueda encontró 40 coincidencias. La auditoría agrupó mensajes,
direcciones y consultas parametrizadas; no encontró concatenación de texto de
usuario dentro de SQL. El filtro de `estado` usa `estado = ?`, los IDs se
validan antes del repositorio y SQLC usa consultas fijas con argumentos. La
prueba guardada en `inyeccion.txt` devolvió `200` al filtro válido y `422
estado_invalido` al intento de inyección, sin ampliar resultados. No se cambió
la consulta porque ya trataba el filtro como dato.

### D2 · El origen que puede leer nuestras respuestas

Se configuró CORS para el origen local `http://localhost:5173`, justo después
de `chi.NewRouter()` y antes de logger, guardia y rutas. Se permiten los
métodos reales `GET`, `POST`, `PUT`, `PATCH`, `DELETE` y los headers
`Authorization` y `Content-Type`. No se usa `*` ni `AllowCredentials`; la
prueba unitaria confirma origen permitido, origen rechazado y preflight de
`PATCH`. El origen no permitido no recibe `Access-Control-Allow-Origin`;
`curl` ejecuta la petición de todas formas porque no aplica la política del
navegador. Cuando publiquen el front, deben reemplazar el origen de desarrollo
por el origen exacto de producción; no deben abrirlo a `*`. Falta guardar las
salidas reales de `curl`.

### D3 · El límite del login

El límite se aplica exclusivamente a `POST /login`: 5 intentos por IP dentro
de una ventana de un minuto. Las demás rutas no quedan limitadas. Se usa
`middleware.ClientIPFromRemoteAddr`, que toma la dirección TCP remota, y la
función `porIP` lee esa IP del contexto y la normaliza con
`httprate.CanonicalizeIP`. Las pruebas unitarias
verificaron que el sexto intento recibe `429`, `Retry-After` positivo y las
cabeceras de límite. Cinco intentos por minuto reducen la fuerza bruta sin
aplicar el tope general del kit al resto de rutas. Falta guardar en `rafaga.txt`
la respuesta HTTP real y el valor observado de `Retry-After`.

### D4 · Detrás de un proxy

`ClientIPFromRemoteAddr` usa la dirección remota del socket. Si la API se
publica detrás de un proxy, ese middleware podría ver siempre la IP del proxy,
agrupando a muchos clientes en un mismo límite. Una cabecera
`X-Forwarded-For` no debe confiarse sin configurar primero los proxies
confiables y validar la cadena de IP; aceptar cualquier cabecera permitiría
que un cliente cambie su clave y evite el límite.

## Semana 12 · reto: mínimo privilegio en la base

### D5 · Con qué usuario se conecta la API

La configuración de Compose separa al administrador PostgreSQL de
`gymcoach_api`, un rol propio usado por la API. El administrador conserva la
capacidad de crear tablas y datos; `gymcoach_api` recibe `USAGE` sobre el
esquema; `SELECT` sobre `entrenadores` y `usuarios`; `SELECT/INSERT` sobre
`clientes`; y `SELECT/INSERT/UPDATE` sobre `rutinas`. Tiene solo `USAGE` sobre
las secuencias de clientes y rutinas, necesarias para insertar. La API no
recibe `DELETE`: GORM elimina rutinas mediante soft delete, que requiere
`UPDATE`. El diseño no concede `CREATE` en el esquema ni privilegios de
superusuario. La prueba de borrado y la lista de roles siguen pendientes hasta
disponer de Docker Engine activo. El volumen Compose usa el proyecto
`gymcoach-taller12`, aislado del volumen del Taller 11.

### D6 · Quién crea y cambia las tablas

Las migraciones y la siembra se ejecutan en un servicio corto
`migracion` con la conexión administrativa, con `MIGRACION_SOLO=1`.
Después se ejecuta `permisos-api` para conceder permisos específicos por tabla.
La API normal usa
`EJECUTAR_MIGRACIONES=0`, por lo que no debe llamar `AutoMigrate` ni sembrar
datos. Un cambio de modelo requiere reconstruir y ejecutar primero el servicio
administrativo; el usuario `gymcoach_api` no debe hacerlo por sí mismo.

### D7 · Lo que esto no frena

El diseño del usuario limitado no concede permisos para borrar tablas ni crear
roles. Si alguien consigue una inyección, aún podría leer, modificar o crear
registros según los permisos concedidos y consumir secuencias. Las pruebas para
confirmar esos límites siguen pendientes. El aislamiento no reemplaza
validación, autorización por propietario ni tratamiento seguro de datos.

## Semana 13 · los avisos en tiempo real

### D1 · Por qué WebSocket y no SSE ni sondeo

- **Qué avisan y en qué sentido:** GymCoach avisa en tiempo real desde el servidor hacia el navegador del cliente dueño cuando cambia el estado de una rutina (`Aviso` con `tipo: "estado_cambiado"`, `id` y `estado` desde `ServicioRutinas.Cambiar` en `internal/gymcoach/aplicacion/rutinas.go`) y hacia las conexiones de entrenadores cuando se crea una rutina nueva (`tipo: "rutina_nueva"` desde `ServicioRutinas.Crear`).
- **Por qué no basta sondeo ni SSE:** El sondeo periódico (`GET /rutinas/{id}`) retrasa la notificación hasta el siguiente intervalo y genera tráfico repetido sin cambios. Aunque estos avisos son unidireccionales (servidor → cliente), `EventSource` (SSE) está pensado para navegador y no admite fácilmente el mismo transporte en clientes futuros móviles o de escritorio ni reutilizar el canal para mensajes de ida y vuelta; WebSocket cuenta con soporte estándar multiplataforma.
- **Qué cuesta WebSocket:** Mantener el mapa en memoria `Avisos.canales` protegido con `sync.Mutex` en `internal/gymcoach/avisos.go`, limpiar canales al desconectarse, aislar `/ws` fuera de `middleware.Timeout` en `main.go` usando `conn.CloseRead(context.Background())` para no cortar conexiones inactivas, y validar explícitamente el token en query y `OriginPatterns`.

### D2 · Quién recibe cada aviso

- **Llave del mapa y origen del dueño:** La llave de `Avisos.canales` (`map[string][]*suscripcionAviso` en `internal/gymcoach/avisos.go`) es `s.Nombre` (tomado del JWT en `Conectar`). El dueño de la rutina sale de `Rutina.ClienteID` y la relación `Rutina.Cliente.Nombre` obtenida en `ServicioRutinas.Cambiar` (`internal/gymcoach/aplicacion/rutinas.go`) después de confirmar que `ActualizarEstado` modificó la fila.
- **Cómo coinciden:** Se mantuvieron las reglas del Hito 2 (el entrenador asigna `cliente_id` al crear en `POST /rutinas`) sin añadir un campo extra ni romper el `403` del entrenador en `PATCH /rutinas/{id}`, y en `internal/sesion/usuarios.go` (`Sembrar`) se agregaron cuentas `RolCliente` (`"Daniela Vera"` y `"Luis Mera"`) con el mismo nombre exacto que los perfiles sembrados en `internal/gymcoach/semilla.go`.
- **Conexión lenta y dos copias de la API:** Cada suscripción usa un canal con buffer de 8 (`bufferAvisos = 8`) y envío no bloqueante (`select { case sub.ch <- aviso: default: }` en `enviarSinBloquear`): si una conexión está lenta, ese aviso se descarta sin bloquear `Enviar`, al `Mutex` ni a otras conexiones. Si hubiera dos copias de la API detrás de un balanceador, el mapa en memoria no se comparte entre procesos y los avisos quedarían solo en la instancia donde ocurrió la petición HTTP (compartirlos requiere un bus pub/sub externo, fuera de este taller).

### D3 · El token en la dirección

- **Por qué `/ws` va fuera del guardia y el token se revisa antes de `Accept`:** La API `WebSocket` del navegador (`docs/taller13/pagina/index.html`) no permite enviar la cabecera `Authorization: Bearer ...` durante el handshake; si `/ws` pasara por `tokens.Autenticar` (`conSesion`), devolvería siempre `401 sin_token`. Por eso `main.go` registra `r.Get("/ws", avisos.Conectar)` fuera de `conSesion`, y `Conectar` valida `?token=` con `a.tokens.Leer` antes de llamar a `websocket.Accept`, respondiendo `401` (`sin_sesion` o `token_vencido`) sin abrir el socket.
- **Para qué sirve el origen si ya hay token:** `OriginPatterns: []string{origenFront}` (`http://localhost:5173`, compartido con CORS en `main.go`) impide que un sitio ajeno en el navegador abra un WebSocket hacia la API local (protección contra CSWSH), rechazando con `403`.
- **Dónde queda el token y qué pasa si vence:** Al viajar en la query (`/ws?token=...`), `middleware.Logger` de `chi` imprime la ruta completa al finalizar `Conectar` (cuando la conexión se cierra), por lo que cualquier log guardado en `docs/taller13/logs_avisos.txt` debe enmascarar el token antes de versionarse. Si el token vence (`Duracion = 1h` en `internal/sesion/tokens.go`) mientras el socket sigue abierto, la conexión continúa viva porque `Leer` solo corre en el handshake; nuevas peticiones HTTP o una reconexión a `/ws` fallarán con `401`.

## Semana 13 · reto: lo nuevo, al negocio

### D4 · Cómo sabe a quién avisar

- **De dónde sale el rol de cada conexión:** De `s.Rol` contenido en el JWT validado por `a.tokens.Leer` durante `Conectar` (`internal/gymcoach/avisos.go`), almacenado en cada `suscripcionAviso` dentro de `Avisos.canales`.
- **Dónde se decide qué rol se entera y por qué:** Se decide en `Avisos.AvisarCreacion` (`internal/gymcoach/avisos.go`, capa de transporte), filtrando `sub.rol == sesion.RolEntrenador` y `nombre != excluir`. El servicio `ServicioRutinas.Crear` (`internal/gymcoach/aplicacion/rutinas.go`) solo invoca `AvisadorCreacion.AvisarCreacion(creadoPor, Aviso{Tipo: "rutina_nueva", ID: rutina.ID, Estado: rutina.Estado})` después de persistir la rutina, manteniendo la capa de aplicación desacoplada de WebSocket y de los roles del sistema.
- **Si cambian el rol con la conexión abierta:** Como la API no reconsulta `usuarios` mientras el WebSocket está abierto ni reevalúa el rol en cada evento, la conexión conserva el rol del token con el que abrió el handshake hasta cerrarse o reconectarse con un token nuevo.

