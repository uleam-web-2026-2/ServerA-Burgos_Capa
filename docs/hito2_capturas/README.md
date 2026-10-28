# Capturas de respaldo del Hito 2

Guarda aquí capturas tomadas al ejecutar la API y las pruebas localmente. No se incluyen imágenes simuladas porque deben mostrar respuestas reales de este repositorio.

Incluye una captura por cada petición de demostración y una de `go test ./...` en verde. Usa estos nombres, que corresponden a los espacios de evidencia de `docs/hito2_hoja.md`:

1. `01-sin-token.png` — GET sin token (401).
2. `02-login-cliente.png` — inicio de sesión de Ana (200 y token).
3. `03-rutina-con-token.png` — GET con token de Ana (200).
4. `04-entrenador-prohibido.png` — entrenador intenta cambiar el estado (403).
5. `05-transicion-valida.png` — Ana cambia el estado válidamente (200).
6. `06-transicion-prohibida.png` — transición prohibida (422).
7. `07-go-test.png` — `go test ./...` en verde.
