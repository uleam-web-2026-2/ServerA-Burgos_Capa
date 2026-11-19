# Auditoria del punto 1 - Inyeccion

| Linea | ¿Arma SQL? | ¿Llega texto del cliente? | ¿Va aparte? | Que hicimos |
| --- | --- | --- | --- | --- |
| `internal/gymcoach/repositorio/rutinas.go:19` | Si | Si, `estado` llega desde `r.URL.Query().Get("estado")` | Si, GORM usa `estado = ?` y pasa `estado` como argumento | Nada: segura; conserva el texto como dato. |
| `internal/gymcoach/repositorio/rutinas.go:27,39,48,57,66,84` | Si | Los metodos reciben `uint`, IDs validados antes de llegar al repositorio, o estados internos | Si, GORM usa parametros `?` | Nada: las rutas parsean el ID con `ParseUint` y GORM separa los valores. |
| `internal/gymcoach/manejadores.go:110,115,133` | No arma SQL directamente | Si, lee el filtro `estado` de la URL | La validacion ocurre antes de enviar el filtro al repositorio | No es vulnerable: el filtro se valida contra la lista `estadosValidos` y el repositorio usa `?`. |
| `internal/gymcoach/repositorio/sqlc/consultas.sql.go` | Si | Si, recibe IDs y estados como argumentos tipados | Si, SQLC genera consultas fijas y usa parametros | Nada: el archivo generado no interpola texto del cliente. |
| `internal/sesion/usuarios.go:46,80` | Si | Si, el nombre o usuario viene de entrada HTTP | Si, GORM usa `nombre = ?` y pasa el valor como argumento | Nada: segura. |
| `internal/gymcoach/migracion.go` | Si | Las consultas usan datos semilla o columnas fijas, no texto libre recibido por HTTP | Algunas son consultas fijas y GORM usa placeholders donde hay datos externos | Revisar como parte de mantenimiento; no es una entrada de una ruta. |
| `main.go:82` y `internal/config/config.go:35,48` | No | No reciben texto de consulta SQL | No aplica | Mensajes de configuracion, puerto o errores, no SQL. |

## Resultado

Se revisaron 40 líneas coincidentes. No se encontró una consulta donde texto libre
de una ruta se concatene dentro del SQL. El filtro principal usa GORM con
`Where("estado = ?", estado)`, y las rutas consultadas mediante SQLC usan parámetros
generados por SQLC.

## Prueba real

En `inyeccion.txt`:

- El filtro `estado=pendiente` respondió `200 OK` y devolvió una rutina pendiente.
- El filtro `"estado=' OR 1=1 --"` respondió `422` y `estado_invalido`, sin resultados extra.
- El ID con texto respondió `400` y `id_invalido`, sin llegar a la consulta.

No fue necesario cambiar código para cerrar una inyección: la validación y los
placeholders ya estaban presentes antes del Taller 12.
