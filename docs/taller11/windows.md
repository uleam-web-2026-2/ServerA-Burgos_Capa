# Notas de validación en Windows

| Qué falló | Mensaje observado | Causa comprobada | Solución aplicada | Evidencia |
|---|---|---|---|---|
| Conexión a la API publicada en el puerto `8080` | `curl` recibió `Recv failure: Connection was reset` | `Get-NetTCPConnection` mostró que el puerto de host `8080` ya tenía un listener cuyo proceso era `postgres`. | Se cambió solo el puerto local de prueba a `18080` para Compose y `18081` para el contenedor independiente; ambos respondieron correctamente. | `docker compose ps`, `peticiones.md` y `compose.txt` |

No se detuvo ni modificó el proceso que ocupaba el puerto `8080`.
