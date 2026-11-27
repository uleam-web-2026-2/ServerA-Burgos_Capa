# Taller 13 · pendientes y bloqueos

## Código y pruebas automatizadas (completados)

El código y las pruebas automatizadas del **Punto 1**, **Punto 2** (`/ws`,
`Avisos`, `401`/`403`/`101`, concurrencia, entrega selectiva, desconexión y
supervivencia al timeout HTTP) y **Reto** (`rutina_nueva` a `entrenador` sin
acoplar `aplicacion`, excluyendo al creador y a clientes) ya están
implementados y pasando con `go vet ./...`, `go build ./...` y `go test ./...`.
Las evidencias de pruebas Go están en:

- [`prueba_servicio.txt`](./prueba_servicio.txt)
- [`prueba_websocket.txt`](./prueba_websocket.txt)
- [`prueba_reto.txt`](./prueba_reto.txt)

## Evidencias en ejecución con PostgreSQL y navegador (completadas)

Se inició el PostgreSQL local `postgresql-x64-16` en `localhost:5433`, se
restauró la configuración de autenticación después de configurar el entorno y
se ejecutó la API con migraciones y siembra.

1. PostgreSQL y las cuentas `"Daniela Vera"` y `"Luis Mera"` quedaron activos.
2. Las cuatro columnas HTTP reales fueron completadas en
   [`peticiones.md`](./peticiones.md); las cuatro repeticiones dieron
   `401`, `200`, `200`, `403`, `200` y `422`.
3. Los tres `curl` reales contra `/ws` quedaron en `ws_sin_token.txt`,
   `ws_otro_origen.txt` y `ws_permitido.txt`, con `401`, `403` y `101`.
4. Se abrió la página local en dos ventanas con `"Daniela Vera"` y
   `"Luis Mera"`; el aviso llegó solo a Daniela al cambiar la rutina `1`.
   La captura está en [`aviso_navegador.png`](./aviso_navegador.png) y los
   logs sanitizados en [`logs_avisos.txt`](./logs_avisos.txt).

