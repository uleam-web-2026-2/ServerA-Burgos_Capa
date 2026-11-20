#!/bin/sh
set -eu

psql -v ON_ERROR_STOP=1 --dbname "$POSTGRES_DB" \
    -c 'CREATE ROLE gymcoach_api;'

psql -v ON_ERROR_STOP=1 --dbname "$POSTGRES_DB" \
    -v api_db_password="$API_DB_PASSWORD" <<'SQL'
ALTER ROLE gymcoach_api LOGIN PASSWORD :'api_db_password';
SQL
