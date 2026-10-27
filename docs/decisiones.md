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
