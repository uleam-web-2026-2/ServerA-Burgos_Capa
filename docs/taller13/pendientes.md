# Taller 13 · pendientes y bloqueos

## Línea de base HTTP

**Responsable: Persona 2, por acuerdo de la pareja.**

Se intentó iniciar GymCoach localmente con `EJECUTAR_MIGRACIONES=0` y
`MIGRACION_SOLO=0`, para no modificar el esquema ni sembrar datos. La API no
respondió en `http://localhost:8080`; las pruebas con `curl` terminaron en
timeout o reinicio de conexión, no en una respuesta HTTP de la aplicación.

Por eso `peticiones.md` conserva sus seis respuestas como pendientes: esos
intentos no cuentan como evidencia. Antes de repetirlos, confirmar que el
servicio PostgreSQL local esté disponible y que la configuración de `.env` sea
válida, sin copiar valores secretos a los archivos del taller.

## Prueba de navegador

**Responsable: Persona 2, por acuerdo de la pareja.**

La página de dos ventanas está preparada, pero la ruta `/ws` todavía no existe
en esta rama. Persona 2 debe integrar y probar el punto 2 antes de hacer la
prueba visual, tomar la captura o registrar `logs_avisos.txt`.
