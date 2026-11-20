# Peticiones de linea base

Se ejecutaron el 2026-11-19 contra `http://localhost:8090`.
La evidencia completa esta en [fase0_seis.txt](./fase0_seis.txt).

| Peticion | Antes (taller 11) | Con las defensas | Con el usuario de la API (reto) |
| --- | --- | --- | --- |
| 1 · Sin token | `401 Unauthorized` | Pendiente: requiere API activa | Pendiente: requiere usuario operativo |
| 2 · Login correcto | `200 OK` | Pendiente: requiere API activa | Pendiente: requiere usuario operativo |
| 3 · Con token | `200 OK` | Pendiente: requiere API activa | Pendiente: requiere usuario operativo |
| 4 · Rol prohibido | `403 Forbidden` | Pendiente: requiere API activa | Pendiente: requiere usuario operativo |
| 5 · Transicion valida | `200 OK` | Pendiente: requiere API activa | Pendiente: requiere usuario operativo |
| 6 · Transicion prohibida | `422 Unprocessable Entity` | Pendiente: requiere API activa | Pendiente: requiere usuario operativo |
