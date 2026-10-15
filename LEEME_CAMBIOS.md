# S7 · Cambios para GymCoach (ServerA-Burgos_Capa)

Copiar estos archivos SOBRE el repositorio, respetando las rutas. Después, desde la raíz:

    go mod tidy      # descarga golang-jwt y bcrypt y los anota en go.mod / go.sum
    go build ./...
    go test ./...
    go run .

En .env: CLAVE_FIRMA con 32 letras y números o más (hoy dice replace-with-a-random-local-secret).

## Archivos
| Archivo | Qué es |
| --- | --- |
| internal/gymcoach/estado.go | NUEVO. Mapa `siguientes` y PATCH /rutinas/{id} (422 prohibida, 200 permitida) |
| internal/gymcoach/manejadores.go | `Rutas(r, tokens)` con tres niveles; el PUT ya no guarda el estado |
| internal/gymcoach/modelos.go | roles como alias de sesion; `puedeTransicionar` delega en `puedePasar` |
| internal/gymcoach/manejadores_test.go | `enrutadorDePrueba()` entra con sesión |
| internal/gymcoach/sesion_ayuda_test.go | NUEVO (anexo A.5) |
| internal/gymcoach/sin_token_test.go | NUEVO (anexo A.6), con las 7 rutas protegidas |
| internal/gymcoach/estado_test.go | NUEVO. Pruebas propias: PATCH, 403 del entrenador, 403 del cliente |
| internal/sesion/tokens.go | tal cual del taller |
| internal/sesion/guardia.go | del taller, con la envoltura de error de GymCoach |
| internal/sesion/usuarios.go | roles cliente y entrenador; ana = cliente, luis = entrenador |
| main.go | migra y siembra usuarios, POST /login, `Rutas(r, tokens)` |
| redocly.yaml | reglas del curso, sin `security-defined: off` |
| docs/openapi.yaml, docs/api.html | contrato con PATCH, /login, seguridad, 401 y 403 |
| docs/decisiones.md | el suyo + D1, D2, D3 de la semana 7 (BORRADOR: revisen las razones) |

## Las seis peticiones (crear antes una rutina nueva con luis)
1. GET /clientes, sin token → 401
2. POST /login con ana / ana-2026 → 200 y token
3. GET /rutinas/{id}, token de ana → 200
4. PATCH /rutinas/{id} {"estado":"en_progreso"}, token de luis (entrenador) → 403
5. PATCH /rutinas/{id} {"estado":"en_progreso"}, token de ana (cliente) → 200
6. PATCH /rutinas/{otra pendiente} {"estado":"completada"}, token de ana → 422
Para crear rutinas: POST /login con luis / luis-2026 y POST /rutinas con su token.

## Pendiente (no lo hace este paquete)
- pruebas.sh y la colección de Postman quedaron obsoletas: ahora piden token y el PUT ya no cambia el estado.
- docs/hito1_ficha_del_negocio.md (secciones 6 y 7) sigue diciendo que no hay autenticación y que el PUT cambia el estado.
- Quiten `.env` del commit (git status) y escriban la línea «IA: ...» en el segundo -m.
