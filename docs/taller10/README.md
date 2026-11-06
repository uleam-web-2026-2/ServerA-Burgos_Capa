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

### Pendiente de cierre

La verificación de las seis peticiones sigue bloqueada. El PostgreSQL 16 local
está escuchando en `localhost:8080`, por lo que las peticiones HTTP a ese puerto
se reinician; la conexión al `localhost:5432` indicado en `.env` no está
disponible. La autenticación con las credenciales configuradas contra el
PostgreSQL en `8080` también fue rechazada. Por ello, `db/schema.sql` se armó a
partir de los modelos del proyecto, pero todavía no se contrastó con `\d
rutinas`.

La Fase 4 tampoco está completada: cuando la Persona 2 publique sus cambios, la
Persona 1 debe actualizar `taller10`, disponer de la API y la base de datos,
repetir las seis peticiones de la [hoja de verificación](./peticiones.md),
llenar la columna «Después» y guardar la salida real de `go test ./...` en
`docs/taller10/go_test.txt`.

La columna «Antes» indica el bloqueo observado y no presenta respuestas no
verificadas como resultados.

## Persona 2 · DTO y capa HTTP/aplicación

Completar y publicar la refactorización de DTO que está en curso:

1. Definir los DTO de entrada y salida de rutinas y clientes en
   `internal/gymcoach/aplicacion/dtos.go`. Los DTO de salida deben exponer solo
   los campos de la API y representar el cliente sin serializar el modelo
   completo de GORM.
2. Adaptar `internal/gymcoach/aplicacion/rutinas.go` para que los casos de uso
   reciban DTO de entrada y devuelvan DTO de salida, convirtiendo los modelos
   del dominio en un único punto.
3. Adaptar `internal/gymcoach/manejadores.go` para decodificar los DTO de
   entrada y responder con los DTO de salida, conservando las rutas, códigos
   HTTP, validaciones y formato JSON actuales.
4. Evitar que los modelos del dominio sean el contrato JSON de la API. Mantener
   en dominio los campos y reglas que requiere la lógica de negocio.
5. Agregar o actualizar pruebas para la conversión y para comprobar que las
   respuestas JSON mantienen la forma esperada.
6. Ejecutar `gofmt` y `go test ./...`, revisar que el cambio funciona junto al
   repositorio sqlc ya integrado y subir los commits de la Persona 2 a
   `taller10`.

Los cambios locales de DTO ya presentes en `aplicacion/rutinas.go`,
`aplicacion/dtos.go`, `dominio/modelos.go` y `manejadores.go` son el punto de
partida; deben revisarse y completarse antes de darlos por terminados. La
refactorización del repositorio GORM en `repositorio/rutinas.go` también está
pendiente en el árbol de trabajo: conservarla solo si sigue siendo necesaria
para las operaciones que sqlc no cubre, y verificarla junto con el adaptador
`RutinasSqlc`.

## Integración entre ambas personas

1. La Persona 2 termina, prueba y publica sus cambios en `taller10`.
2. La Persona 1 integra esos cambios y comprueba que no se hayan perdido
   archivos o modificaciones de ninguna de las dos partes.
3. La Persona 1 ejecuta la verificación final de la [Fase 4](./peticiones.md)
   con el servidor y PostgreSQL disponibles, registra los resultados reales y
   actualiza `go_test.txt`.
