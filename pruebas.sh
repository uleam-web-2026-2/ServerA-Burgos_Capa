#!/usr/bin/env bash
# Uso: ./pruebas.sh (contra localhost:8080)
# Uso alternativo: ./pruebas.sh http://192.168.1.20:8080

B=${1:-http://localhost:8080}

c() {
  echo
  echo "### $1 [esperado: $2]"
  shift 2
  curl -s -i "$@" | sed -n '1p; /^{/p'
}

c "GET lista de clientes" "200" "$B/clientes"
c "GET clientes con rutinas pendientes" "200" "$B/clientes?estado=pendiente"
c "GET filtro de estado inválido" "422" "$B/clientes?estado=terminada"
c "POST cliente válido" "201" -X POST "$B/clientes" -H "Content-Type: application/json" \
  -d '{"nombre":"María Solórzano","correo":"maria@example.test","objetivo":"Mejorar fuerza general","entrenador_id":1}'
c "POST rutina válida" "201" -X POST "$B/rutinas" -H "Content-Type: application/json" \
  -d '{"titulo":"Fuerza inicial","descripcion":"Cuerpo completo con cargas progresivas","dias_semana":3,"cliente_id":1}'
c "GET rutina" "200" "$B/rutinas/1"
c "GET rutina inexistente" "404" "$B/rutinas/99999"
c "GET id no numérico" "400" "$B/rutinas/abc"
c "POST JSON roto" "400" -X POST "$B/rutinas" -H "Content-Type: application/json" -d '{"titulo":"x'
c "POST días fuera de rango" "422" -X POST "$B/rutinas" -H "Content-Type: application/json" \
  -d '{"titulo":"Fuerza","descripcion":"Rutina","dias_semana":8,"cliente_id":1}'
c "PUT transición prohibida" "409" -X PUT "$B/rutinas/1" -H "Content-Type: application/json" \
  -d '{"titulo":"Fuerza inicial","descripcion":"Cuerpo completo","dias_semana":3,"estado":"completada","cliente_id":1}'
c "DELETE rutina" "200" -X DELETE "$B/rutinas/1"
c "ruta inexistente" "404" "$B/entrenadores"
