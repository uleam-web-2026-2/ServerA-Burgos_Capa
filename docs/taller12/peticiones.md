# Peticiones de linea base

Se ejecutaron el 2026-11-19 contra `http://localhost:8090`.
La evidencia completa esta en [fase0_seis.txt](./fase0_seis.txt).

| Peticion | Antes (taller 11) | Con las defensas | Con el usuario de la API (reto) |
| --- | --- | --- | --- |
| 1 · Sin token | `401 Unauthorized` | Pendiente Persona 2 | Pendiente reto |
| 2 · Login correcto | `200 OK` | Pendiente Persona 2 | Pendiente reto |
| 3 · Con token | `200 OK` | Pendiente Persona 2 | Pendiente reto |
| 4 · Rol prohibido | `403 Forbidden` | Pendiente Persona 2 | Pendiente reto |
| 5 · Transicion valida | `200 OK` | Pendiente Persona 2 | Pendiente reto |
| 6 · Transicion prohibida | `422 Unprocessable Entity` | Pendiente Persona 2 | Pendiente reto |
