# Taller 12: qué se implementó

En este taller revisamos cómo responde la API ante consultas, solicitudes desde
el navegador y varios intentos de inicio de sesión. También preparamos una
configuración de PostgreSQL para que la API use solo los permisos que necesita.

## Revisión de consultas

Auditamos las coincidencias relacionadas con consultas y SQL. El filtro por
estado ya usa parámetros (`estado = ?`) y las consultas SQLC separan el SQL de
los valores recibidos. No encontramos concatenación directa de datos del
usuario en las consultas revisadas, así que no fue necesario cambiar esa parte
del código.

## CORS

Configuramos CORS para aceptar el origen de desarrollo
`http://localhost:5173`, sin abrir el acceso a cualquier sitio. Se permiten los
métodos que usa la API y los encabezados `Authorization` y `Content-Type`.
También añadimos pruebas unitarias para el origen permitido, uno rechazado y
una solicitud previa (`preflight`).

## Límite de intentos en el inicio de sesión

El límite se aplica solamente a `POST /login`: cinco intentos por dirección IP
dentro de un minuto. La dirección se obtiene de la conexión y se normaliza; no
se confía en una cabecera que el cliente pueda modificar. Las pruebas unitarias
comprueban que el sexto intento recibe `429`.

## Acceso de la API a PostgreSQL

En Compose separamos las tareas administrativas del usuario que atiende las
solicitudes. Un servicio ejecuta las migraciones y la siembra; la API utiliza
el rol `gymcoach_api`, sin ejecutar migraciones al arrancar. Se le asignan
permisos por tabla y uso de las secuencias necesarias. No se le concede
privilegio de superusuario ni `CREATE` en el esquema. Las rutinas usan borrado
lógico, por lo que la API requiere `UPDATE`, no `DELETE`.

La configuración de base de datos está en [compose.yaml](../../compose.yaml),
los controles HTTP en [main.go](../../main.go) y las pruebas correspondientes
en [main_test.go](../../main_test.go).

## Comprobaciones y estado

Pasaron `go test ./... -count=1`, `go build ./...`,
`docker compose --env-file .env.compose.example config --quiet` y
`git diff --check`.

La configuración de Compose se validó, pero no se pudo iniciar Docker Engine
para probar el contenedor y los permisos de PostgreSQL en ejecución. Por eso,
este documento describe lo implementado y comprobado en código; no afirma que
las pruebas HTTP o de base de datos se hayan ejecutado.
