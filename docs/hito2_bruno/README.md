# Colección Bruno · Hito 2

1. Selecciona el entorno `Local` y configura `cliente_id` con un cliente que exista.
2. Ejecuta `preparacion/Login Entrenador.bru`; el token se guarda en `token_entrenador` y vence en una hora.
3. Ejecuta `0 - Preparar rutina.bru`; guarda el nuevo id automáticamente en `rutina_id`.
4. Ejecuta las peticiones 1 a 6 en orden. La petición 2 obtiene en vivo el token de Ana. La petición 4 verifica que Luis no puede cambiar el estado; la 5 usa Ana para iniciar la rutina y la 6 comprueba que no se permite volver a `pendiente`.

Para repetir el ensayo, inicia sesión como entrenador y ejecuta de nuevo la petición 0 para obtener un registro recién creado. La carpeta `preparacion` es una ayuda operativa; las seis peticiones calificadas siguen numeradas del 1 al 6, además de la petición 0.
