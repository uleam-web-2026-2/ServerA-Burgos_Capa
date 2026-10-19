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

**Transición:** de `pendiente` a `completada`.

**Por qué el negocio no la permite:** el cliente no puede marcar como terminado un plan que todavía no ha empezado. Primero debe pasar a `en_progreso`.

**Qué responde y dónde se valida:** HTTP `422`, código `transicion_prohibida` y mensaje `de pendiente no se puede pasar a completada`. La transición se comprueba en `internal/gymcoach/estado.go`, función `cambiarEstado`, usando el mapa `siguientes`.
