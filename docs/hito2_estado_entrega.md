# Hito 2 · Estado del trabajo y guía para completar la entrega

## Resumen

El proyecto contiene el trabajo del Hito 2 descrito en la hoja de demostración. En esta revisión se añadieron espacios documentales para las evidencias; no se modificó el código de la API ni se generaron capturas.

## Trabajo realizado en esta conversación

- Se añadió la sección **5 · Evidencias** a `docs/hito2_hoja.md`, con espacios para seis respuestas de la API y una ejecución de pruebas.
- Se actualizaron las instrucciones de `docs/hito2_capturas/README.md` con los nombres de archivo correspondientes.
- Se ejecutó `go test ./...`; todas las pruebas terminaron correctamente.
- Se verificó que `.env` no está versionado en Git.

## Estado de las evidencias

Las imágenes todavía no están en `docs/hito2_capturas/`. Guarda las capturas reales con estos nombres para que aparezcan en la hoja:

| Archivo | Evidencia que debe mostrar |
| --- | --- |
| `01-sin-token.png` | Petición GET sin token y respuesta HTTP 401. |
| `02-login-cliente.png` | Inicio de sesión de Ana y respuesta HTTP 200 con token. Evita mostrar credenciales o el token completo. |
| `03-rutina-con-token.png` | Consulta de rutina autenticada y respuesta HTTP 200. |
| `04-entrenador-prohibido.png` | Intento del entrenador de cambiar el estado y respuesta HTTP 403. |
| `05-transicion-valida.png` | Cambio válido del estado por Ana y respuesta HTTP 200. |
| `06-transicion-prohibida.png` | Intento de volver a `pendiente` y respuesta HTTP 422 con `transicion_prohibida`. |
| `07-go-test.png` | Terminal mostrando `go test ./...` completado sin errores. |

Usa datos de demostración locales. Comprueba que cada captura sea legible y que no muestre secretos ni información personal real.

## Pasos pendientes

1. Inicia la API localmente y selecciona el entorno `Local` en Bruno.
2. Configura `cliente_id` con un cliente existente.
3. Ejecuta `preparacion/Login Entrenador.bru` para obtener el token del entrenador.
4. Ejecuta `0 - Preparar rutina.bru` para crear una rutina nueva.
5. Ejecuta las peticiones `1` a `6` en orden y toma las seis capturas indicadas arriba.
6. Ejecuta `go test ./...` y guarda la captura de la ejecución exitosa.
7. Ejecuta `npx @redocly/cli lint docs/openapi.yaml --format=stylish` y confirma que termina sin errores ni avisos.
8. Revisa cada elemento del recuadro de cumplimiento en `docs/hito2_hoja.md` y marca solo los puntos que hayas comprobado.
9. Antes de crear un commit de entrega, revisa `LEEME_CAMBIOS.md` y añade la línea `IA: ...` indicada allí, si corresponde a las instrucciones del curso.

## Verificaciones registradas

| Verificación | Resultado |
| --- | --- |
| `go test ./...` | Correcto en esta sesión. |
| `.env` versionado en Git | No; no aparece como archivo rastreado. |
| Lint de `docs/openapi.yaml` | Pendiente de ejecutar. |
| Capturas de Bruno y pruebas | Pendientes de añadir. |

## Archivos relacionados

- `docs/hito2_hoja.md`: hoja de demostración y espacios para las evidencias.
- `docs/hito2_capturas/README.md`: instrucciones y nombres de las capturas.
- `docs/hito2_bruno/README.md`: orden para preparar y ejecutar la colección Bruno.
- `docs/decisiones.md`: decisiones de arquitectura y estructura por capas.
- `LEEME_CAMBIOS.md`: notas operativas del taller y pasos de entrega.
