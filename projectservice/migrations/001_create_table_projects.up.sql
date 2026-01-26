CREATE TABLE IF NOT EXISTS projects(
    id SERIAL PRIMARY KEY,
    owner_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT unique_owner_id UNIQUE (owner_id, name)
);

CREATE INDEX idx_projects_owner_id ON projects(owner_id);