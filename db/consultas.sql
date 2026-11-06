-- name: ObtenerRutina :one
SELECT id, created_at, updated_at, titulo, descripcion, dias_semana, estado, cliente_id
FROM rutinas
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: ActualizarEstado :execrows
UPDATE rutinas
SET estado = sqlc.arg(nuevo_estado),
    updated_at = now()
WHERE id = sqlc.arg(id)
  AND estado = sqlc.arg(estado_actual)
  AND deleted_at IS NULL;
