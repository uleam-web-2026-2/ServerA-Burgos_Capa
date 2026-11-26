# Regresión HTTP · Taller 13

Persona 2 debe completar esta tabla con respuestas reales ejecutadas en la rama
`taller13`. La columna
«Código esperado» es una referencia del Hito 2, no evidencia de una ejecución.
No pegar tokens ni contraseñas en este archivo.

| # | Petición Bruno | Código esperado | Antes del Taller 13 | Con el punto 1 | Con WebSocket | Con el reto |
| --- | --- | ---: | --- | --- | --- | --- |
| 1 | GET Rutina Sin Token | 401 | Pendiente | Pendiente | Pendiente | Pendiente |
| 2 | Login Cliente | 200 | Pendiente | Pendiente | Pendiente | Pendiente |
| 3 | GET Rutina Con Token | 200 | Pendiente | Pendiente | Pendiente | Pendiente |
| 4 | PATCH Entrenador Prohibido | 403 | Pendiente | Pendiente | Pendiente | Pendiente |
| 5 | PATCH Transición Válida | 200 | Pendiente | Pendiente | Pendiente | Pendiente |
| 6 | PATCH Transición Prohibida | 422 | Pendiente | Pendiente | Pendiente | Pendiente |

La petición 5 consume la rutina al cambiar su estado; crear otra antes de
repetir la secuencia. El límite de inicios de sesión del Taller 12 puede
responder 429: detener y volver a iniciar el servidor antes de repetir las seis
peticiones.

Al completar cada columna, registrar el código observado y guardar la evidencia
correspondiente. La respuesta completa de la petición 3 se necesita para la
comparación con el Hito 2.
