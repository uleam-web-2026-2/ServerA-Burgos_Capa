# Registro de decisiones de arquitectura

## D1 Modelo Entrenador-Cliente-Rutina

Separamos el entrenador, el cliente y la rutina en entidades relacionadas. Un entrenador administra varios clientes y cada cliente puede tener distintas rutinas a lo largo del tiempo. Esto evita guardar programas completos dentro del perfil del cliente y permite filtrar y cambiar el estado de cada plan de forma independiente.

## D2 Objetivo y datos de salud

El MVP guarda un objetivo general como texto, pero no guarda diagnósticos, lesiones, peso, medidas ni historial médico. Estos datos pueden ser sensibles; se debe validar consentimiento, controles de acceso y política de retención antes de ampliar el modelo.

## D3 Máquina de estados

La rutina empieza `pendiente`, pasa a `en_progreso` cuando el cliente inicia el plan y termina `completada` cuando finaliza el ciclo. Una rutina completada no se reinicia; el entrenador crea una nueva para conservar trazabilidad.

## D4 Permisos frente a implementación

Entrenador y cliente tienen permisos distintos definidos en la ficha. Las rutas actuales no autentican la identidad ni aplican autorización; son una API de desarrollo, no un control de acceso.
