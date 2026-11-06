CREATE TABLE entrenadores (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    nombre TEXT NOT NULL,
    correo TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_entrenadores_correo ON entrenadores (correo);
CREATE INDEX idx_entrenadores_deleted_at ON entrenadores (deleted_at);

CREATE TABLE clientes (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    nombre TEXT NOT NULL,
    correo TEXT NOT NULL,
    objetivo TEXT NOT NULL,
    entrenador_id BIGINT NOT NULL REFERENCES entrenadores(id)
);

CREATE UNIQUE INDEX idx_clientes_correo ON clientes (correo);
CREATE INDEX idx_clientes_deleted_at ON clientes (deleted_at);

CREATE TABLE rutinas (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    titulo TEXT NOT NULL,
    descripcion TEXT NOT NULL,
    dias_semana BIGINT NOT NULL,
    estado TEXT NOT NULL DEFAULT 'pendiente',
    cliente_id BIGINT NOT NULL REFERENCES clientes(id)
);

CREATE INDEX idx_rutinas_deleted_at ON rutinas (deleted_at);
