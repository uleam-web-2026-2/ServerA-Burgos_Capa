#!/usr/bin/env bash
# Requiere Bash, curl y jq. Uso: ./pruebas.sh (contra localhost:8080)
# Uso alternativo: ./pruebas.sh http://192.168.1.20:8080
set -euo pipefail

B=${1:-http://localhost:8080}

command -v jq >/dev/null || { echo "Se requiere jq para leer el token del login." >&2; exit 1; }

login() {
  curl -fsS -X POST "$B/login" -H "Content-Type: application/json" -d "$1" |
    jq -er '.datos.token'
}

TOKEN_ENTRENADOR=$(login '{"usuario":"luis","clave":"luis-2026"}')
TOKEN_CLIENTE=$(login '{"usuario":"ana","clave":"ana-2026"}')
AUTH_ENTRENADOR=(-H "Authorization: Bearer $TOKEN_ENTRENADOR")
AUTH_CLIENTE=(-H "Authorization: Bearer $TOKEN_CLIENTE")

c() {
  local nombre esperado respuesta codigo cuerpo
  nombre=$1
  esperado=$2
  shift 2
  echo
  echo "### $nombre [esperado: $esperado]"
  respuesta=$(curl -sS -w $'\n%{http_code}' "$@")
  codigo=${respuesta##*$'\n'}
  cuerpo=${respuesta%$'\n'*}
  RESPUESTA_BODY=$cuerpo
  printf 'HTTP %s\n' "$codigo"
  if [[ -n $cuerpo ]]; then printf '%s\n' "$cuerpo"; fi
  if [[ $codigo != "$esperado" ]]; then
    echo "Se esperaba HTTP $esperado y llegó $codigo." >&2
    return 1
  fi
}

c "GET lista de clientes sin token" "401" "$B/clientes"
c "GET clientes como cliente" "403" "${AUTH_CLIENTE[@]}" "$B/clientes"
c "GET lista de clientes como entrenador" "200" "${AUTH_ENTRENADOR[@]}" "$B/clientes"
CLIENTE_ID=$(printf '%s' "$RESPUESTA_BODY" | jq -er '.datos[0].id // empty') || {
  echo "No hay clientes para ejecutar las pruebas; cree uno primero." >&2
  exit 1
}
c "GET lista de rutinas" "200" "${AUTH_ENTRENADOR[@]}" "$B/rutinas"
c "GET filtro de estado inválido" "422" "${AUTH_ENTRENADOR[@]}" "$B/rutinas?estado=terminada"
c "POST rutina como cliente" "403" -X POST "${AUTH_CLIENTE[@]}" "$B/rutinas" \
  -H "Content-Type: application/json" \
  -d '{"titulo":"Fuerza","descripcion":"Rutina","dias_semana":3,"cliente_id":1}'
c "POST rutina válida" "201" -X POST "${AUTH_ENTRENADOR[@]}" "$B/rutinas" \
  -H "Content-Type: application/json" \
  -d "{\"titulo\":\"Fuerza integradora\",\"descripcion\":\"Cuerpo completo\",\"dias_semana\":3,\"cliente_id\":$CLIENTE_ID}"
RUTINA_ID=$(printf '%s' "$RESPUESTA_BODY" | jq -er '.datos.id')
c "GET rutina" "200" "${AUTH_ENTRENADOR[@]}" "$B/rutinas/$RUTINA_ID"
c "GET rutina inexistente" "404" "${AUTH_ENTRENADOR[@]}" "$B/rutinas/99999"
c "GET id no numérico" "400" "${AUTH_ENTRENADOR[@]}" "$B/rutinas/abc"
c "PATCH transición prohibida" "422" -X PATCH "${AUTH_CLIENTE[@]}" "$B/rutinas/$RUTINA_ID" \
  -H "Content-Type: application/json" -d '{"estado":"completada"}'
c "PATCH transición permitida" "200" -X PATCH "${AUTH_CLIENTE[@]}" "$B/rutinas/$RUTINA_ID" \
  -H "Content-Type: application/json" -d '{"estado":"en_progreso"}'
c "PUT no cambia el estado" "200" -X PUT "${AUTH_ENTRENADOR[@]}" "$B/rutinas/$RUTINA_ID" \
  -H "Content-Type: application/json" \
  -d "{\"titulo\":\"Fuerza integradora\",\"descripcion\":\"Cuerpo completo\",\"dias_semana\":3,\"estado\":\"completada\",\"cliente_id\":$CLIENTE_ID}"
ESTADO=$(printf '%s' "$RESPUESTA_BODY" | jq -er '.datos.estado')
if [[ $ESTADO != "en_progreso" ]]; then
  echo "PUT cambió el estado a $ESTADO; se esperaba en_progreso." >&2
  exit 1
fi
c "PATCH completar rutina" "200" -X PATCH "${AUTH_CLIENTE[@]}" "$B/rutinas/$RUTINA_ID" \
  -H "Content-Type: application/json" -d '{"estado":"completada"}'
c "DELETE rutina" "200" -X DELETE "${AUTH_ENTRENADOR[@]}" "$B/rutinas/$RUTINA_ID"
c "GET ruta inexistente" "404" "${AUTH_ENTRENADOR[@]}" "$B/entrenadores"
