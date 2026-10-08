# Hito 1 · GymCoach

Guion de contenido para las cinco diapositivas de la plantilla `S4_Hito1_Servidor_Plantilla_diapositivas.pptx`.

1. **El negocio:** GymCoach administra clientes y rutinas para entrenadores personales. Referencia: ALTR Project, SaaS para gimnasios y coaches, con ingresos mensuales reportados de USD 7.000–8.000 en Starter Story. Contraste: QuickCoach cerró tras sostener USD 11.000 MRR pagados frente a USD 31.000 de gastos, según su fundador.
2. **Adaptación al Ecuador:** cobro recurrente en USD con medios disponibles localmente; emisión de comprobantes electrónicos y validación del régimen tributario; uso móvil y tolerancia a conectividad variable; minimización de datos de bienestar.
3. **Modelo de datos:** Entrenador 1:N Cliente; Cliente 1:N Rutina. Atributos principales: nombre/correo, objetivo del cliente, descripción/días/estado de la rutina.
4. **Máquina de estados:** Rutina: `pendiente → en_progreso → completada`. El cliente inicia y termina el ciclo. No se permite completar una rutina que nunca empezó.
5. **API y pantallas:** `GET/POST /clientes`, `POST /rutinas`, `GET/PUT/DELETE /rutinas/{id}`; cada endpoint se relaciona con panel, alta de clientes o gestión de rutinas. Falta ejecutar pruebas reales y agregar autenticación/autorización.
