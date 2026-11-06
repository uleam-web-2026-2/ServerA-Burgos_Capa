CREATE TABLE rutinas (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    titulo TEXT NOT NULL,
    descripcion TEXT NOT NULL,
    dias_semana BIGINT NOT NULL,
    estado TEXT NOT NULL DEFAULT 'pendiente',
    cliente_id BIGINT NOT NULL
);

CREATE INDEX idx_rutinas_deleted_at ON rutinas (deleted_at);
