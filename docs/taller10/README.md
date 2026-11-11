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

### Verificación previa a los últimos ajustes

La Persona 1 registró seis peticiones contra una versión local anterior con
sqlc y los primeros cambios DTO; la salida de `go test ./...` quedó en
[go_test.txt](./go_test.txt). Esos resultados son un antecedente de verificación,
no sustituyen la repetición de Fase 4 contra los últimos cambios descritos en
este documento.

## Persona 2 · DTO y capa HTTP/aplicación

Implementación completada localmente y validada con pruebas automatizadas.

1. Se definieron los DTO de entrada y salida de clientes y rutinas en
   `internal/gymcoach/aplicacion/dtos.go`. Los DTO de salida exponen solo los
   campos públicos de la API y no serializan el modelo completo de GORM.
   La lista de clientes tiene un DTO específico con sus relaciones de
   entrenador y rutinas.
2. Se adaptó `internal/gymcoach/aplicacion/rutinas.go` para que los casos de uso
   reciban DTO de entrada y devuelvan DTO de salida, centralizando la
   conversión desde y hacia el dominio.
3. Se adaptó `internal/gymcoach/manejadores.go` para decodificar los DTO de
   entrada y responder con DTOs de salida, sin cambiar las rutas, códigos HTTP,
   validaciones ni envoltura JSON ya acordadas. El DTO del `PUT` no incluye
   `estado`; ese campo se cambia exclusivamente mediante `PATCH`.
4. Se evitó que los modelos del dominio fueran el contrato JSON de la API,
   manteniendo en dominio solo los campos y reglas de negocio necesarias.
5. Se añadieron pruebas para la conversión y para verificar que las respuestas
   JSON mantienen la forma esperada, evitando fugas de campos internos.
6. Se ejecutó `gofmt` y `go test ./...` con resultado correcto. También se
   actualizó `docs/openapi.yaml`, se validó con Redocly y se regeneró
   `docs/api.html`.
7. Se documentó la sección de transacciones: la migración de datos antiguos
   usa una transacción GORM; el cambio de estado usa una actualización atómica
   condicional con sqlc para manejar concurrencia.

Estado: la implementación de Persona 2 quedó completada en el árbol de trabajo.
La verificación manual de las seis peticiones de Fase 4 debe repetirse contra
esta revisión junto con PostgreSQL antes de declarar cerrada la integración
publicada. El repositorio conserva GORM para las operaciones no cubiertas por
sqlc y usa el adaptador `RutinasSqlc` donde corresponde.

## Integración entre ambas personas

1. La Persona 2 terminó la refactorización DTO, las pruebas automatizadas y la
   actualización de la documentación en el árbol de trabajo.
2. Falta integrar/publicar esta revisión en `taller10` y repetir la verificación
   de la [Fase 4](./peticiones.md) contra la versión exacta que se entregará.

La implementación de la Persona 2 está completa; el cierre de integración queda
sujeto a la verificación manual final en el entorno con PostgreSQL.
