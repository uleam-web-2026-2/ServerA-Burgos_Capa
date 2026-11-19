# REMID · Reparto Taller 12

## Responsabilidad inmediata: Persona 1

**Estado:** fase 0 completada y punto 1 completado. Persona 2 continúa con puntos 2 y 3.

### Hecho por Persona 1

- Revisada la estructura del repositorio y las rutas que reciben `estado` e `id`.
- Preparada `docs/taller12/patrones.txt` con los patrones oficiales.
- Preparada `docs/taller12/auditoria.md` con una fila por familia de linea.
- Comprobado que el filtro actual usa un placeholder GORM:
  `Where("estado = ?", estado)`.
- Comprobado que las consultas SQLC usan consultas fijas y argumentos separados.
- Ejecutado `go get` para incorporar chi 5.5.0, CORS y httprate.
- Resueltos los marcadores de merge preexistentes mediante la rama limpia basada en
  `origin/taller11`.
- `go build ./...` y `go test ./...` pasan.
- Generados `auditoria.txt` y `auditoria_despues.txt` mediante `git grep`.
- Levantada una sola API en el puerto 8090 y guardada la prueba real de inyección.
- Ejecutadas las seis peticiones de línea base: 401, 200, 200, 403, 200, 422.
- La rama `taller12` fue publicada en GitHub.

### Evidencia pendiente real de Persona 1

- Hecho: el commit del punto 1 es `7efbfdb` y fue publicado en `taller12`.

## Responsabilidad siguiente: Persona 2

### Punto 2 · CORS

- Obtener `git pull` despues de que Persona 1 publique el commit.
- Nombrar `http://localhost:5173` o el origen real del front en
  `r.Use(cors.Handler(...))`.
- Colocar CORS inmediatamente despues de `chi.NewRouter()` y antes de logger,
  guardia y rutas.
- Permitir los metodos HTTP reales y los headers `Authorization`,
  `Content-Type`.
- Crear `docs/taller12/cors_permitido.txt`,
  `cors_no_permitido.txt`, `cambio_estado.json` y
  `cors_curl_ejecuta.txt`.
- Crear el commit del punto 2 con el mismo formato de declaracion de IA.

### Punto 3 · Limite en login y D1-D4

- Aplicar el limite solo a `r.With(httprate.LimitBy(...)).Post("/login", ...)`.
- Mantener un limite distinto y explicado para el login.
- Usar `middleware.ClientIPFromRemoteAddr` y `porIP`.
- Ejecutar la rafaga y guardar `login_malo.json`, `rafaga.txt` y
  `logs_seis.txt`.
- Completar la columna «Con las defensas» de `peticiones.md`.
- Escribir D1-D4 en `docs/decisiones.md`.
- Reiniciar el servidor antes de probar las seis peticiones para limpiar el
  contador en memoria.
- Crear el commit del punto 3.

## Reto compartido: Persona 1 despues del punto 3

La guia propone que Persona 1 haga el reto al regresar el turno:

- Crear un usuario PostgreSQL propio sin superusuario.
- Otorgarle solo permisos `SELECT`, `INSERT`, `UPDATE` o `DELETE` segun cada
  tabla, incluyendo `USAGE` sobre secuencias.
- Separar las migraciones y la siembra del usuario usado por la API.
- Eliminar o retirar `AutoMigrate` del usuario operativo.
- Probar un arranque limpio desde Compose o PostgreSQL vacio.
- Guardar `reto_drop.txt` y `reto_usuarios.txt`.
- Completar la columna del reto en `peticiones.md` y escribir D5-D7.
- Crear el commit del reto.

## Reglas de entrega

- Trabajar en la rama `taller12`.
- Revisar cambios cruzados antes de cada commit.
- No incluir `.env`, tokens, contraseñas ni claves privadas.
- Declarar la herramienta de IA y su uso en el cuerpo de cada commit.
- Antes de publicar: `git status`, `git build`/`go build`, `go test` y
  `git log origin/taller12`.
