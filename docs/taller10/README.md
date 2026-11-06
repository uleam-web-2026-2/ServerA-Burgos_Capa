# Taller 10 · División del trabajo

El taller se divide entre dos personas para trabajar en paralelo sobre la API
GymCoach. La Persona 1 prepara la verificación y la persistencia con sqlc; la
Persona 2 define los DTO de la capa de aplicación y adapta el manejo HTTP para
usarlos. Después se integran ambos trabajos y se repiten las pruebas.

## Persona 1 · Base de datos y verificación

### Trabajo realizado

- Creó y publicó la rama `taller10` antes de integrar los cambios de DTO.
- Tomó como línea base el commit `a756a7c`, punto de partida de `taller10`, y
  confirmó que `internal/gymcoach/aplicacion/rutinas.go` no importa GORM.
- Instaló sqlc 1.31.1 y creó `sqlc.yaml`, `db/schema.sql` y `db/consultas.sql`.
- Definió `ObtenerRutina` (`:one`) y `ActualizarEstado` (`:execrows`). Ambas
  consultas excluyen filas con `deleted_at`; la actualización también fija
  `updated_at = now()`.
- Generó el código sqlc y creó el adaptador `RutinasSqlc`, que traduce
  `sql.ErrNoRows` a `dominio.ErrRutinaNoEncontrada`.
- Conectó sqlc en `main.go` usando `db.DB()` y añadió una prueba para la
  traducción de error.
- Subió estos cambios a `origin/taller10`. `go test ./...`, `go vet ./...` y
  `sqlc generate` terminaron correctamente.
- Contrastó `db/schema.sql` con `\d rutinas`, `\d clientes` y
  `\d entrenadores` en PostgreSQL 16. Ejecutó las seis peticiones «Antes» con
  la API base en una base aislada; todos los códigos y transiciones coincidieron
  con lo esperado.

### Cierre de la fase

La refactorización DTO de la Persona 2 quedó integrada y validada en el árbol de
trabajo actual. Como comprobación adicional, la Persona 1 ejecutó las seis
peticiones contra la versión local con los cambios DTO y guardó la salida de
`go test ./...` en [go_test.txt](./go_test.txt). En la versión final del taller,
la parte de DTO quedó cerrada y la verificación de la API con los cambios
recientes se considera completada.

## Persona 2 · DTO y capa HTTP/aplicación

Trabajo completado y validado.

1. Se definieron los DTO de entrada y salida de clientes y rutinas en
   `internal/gymcoach/aplicacion/dtos.go`. Los DTO de salida exponen solo los
   campos públicos de la API y no serializan el modelo completo de GORM.
2. Se adaptó `internal/gymcoach/aplicacion/rutinas.go` para que los casos de uso
   reciban DTO de entrada y devuelvan DTO de salida, centralizando la
   conversión desde y hacia el dominio.
3. Se adaptó `internal/gymcoach/manejadores.go` para decodificar los DTO de
   entrada y responder con DTOs de salida, sin cambiar las rutas, códigos HTTP,
   validaciones ni formato JSON ya acordados.
4. Se evitó que los modelos del dominio fueran el contrato JSON de la API,
   manteniendo en dominio solo los campos y reglas de negocio necesarias.
5. Se añadieron pruebas para la conversión y para verificar que las respuestas
   JSON mantienen la forma esperada, evitando fugas de campos internos.
6. Se ejecutó `gofmt` y `go test ./...` con resultado correcto, confirmando
   que el cambio funciona junto con el repositorio sqlc ya integrado.

Estado final: la refactorización de DTO de la Persona 2 quedó completada,
validada y consistente con el repositorio actual. La revisión del repositorio
GORM en `repositorio/rutinas.go` quedó subordinada a la necesidad real del
proyecto: si permanece, debe seguir siendo compatible con el adaptador
`RutinasSqlc`, pero la tarea principal de DTO y capa HTTP quedó cerrada con la
integración del servicio y los handlers actualizados.

## Integración entre ambas personas

1. La Persona 2 terminó, probó y dejó cerrada la refactorización de DTO en el
   árbol de trabajo del taller.
2. La Persona 1 integró y revisó que no se hubieran perdido archivos ni cambios
   de ninguna de las dos partes.
3. La Persona 1 ejecutó la verificación final de la [Fase 4](./peticiones.md)
   con el servidor y PostgreSQL disponibles y registró los resultados relevantes
   en `go_test.txt`.

En este punto, la división del trabajo queda cerrada y la parte correspondiente
al trabajo de la Persona 2 queda marcada como completada.
