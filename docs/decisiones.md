# Registro de decisiones de arquitectura

## D1 Modelado de relación Cliente-Rutina
**Opciones:**
- Opción A: guardar toda la información de rutinas dentro de la entidad `Cliente` con anidación en un único objeto.
- Opción B: separar la información en dos entidades, `Cliente` (lado uno) y `Rutina` (lado muchos), con `ClienteID` como clave foránea.

**Qué elegimos:**
Elegimos la opción B: mantener `Cliente` y `Rutina` como entidades separadas con una relación 1 a N. Un cliente puede tener muchas rutinas, pero cada rutina pertenece a un único cliente.

**Por qué en nuestro negocio:**
En GymCoach, cada cliente tiene un historial de rutinas personalizadas, y cada rutina necesita su propio estado, título y trazabilidad. Separar las entidades facilita la consulta de clientes con sus rutinas, la actualización del estado de cada práctica y la gestión independiente de cada ejercicio o plan asignado.

**Qué pasaría con la otra opción:**
Si se guardaran todas las rutinas dentro de `Cliente`, el sistema se volvería más rígido y difícil de consultar por estado, además de complicar el CRUD de rutinas como recurso independiente. Un cliente con un historial grande terminaría con estructuras anidadas más complejas y más riesgo de inconsistencia al actualizar o eliminar una rutina en particular.

**Conclusión:**
La relación 1 a N es la opción más clara para un sistema de entrenamiento, porque preserva la integridad del dominio y facilita la evolución del servicio con nuevas reglas de negocio.