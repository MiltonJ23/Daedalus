-- +migrate Up
CREATE TABLE projects (
    id              VARCHAR(36) PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    industry_type   VARCHAR(120) NOT NULL,
    location        VARCHAR(255) NOT NULL,
    budget          DOUBLE PRECISION NOT NULL DEFAULT 0,
    floor_width     DOUBLE PRECISION NOT NULL,
    floor_depth     DOUBLE PRECISION NOT NULL,
    target_capacity VARCHAR(255),
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    version         INTEGER NOT NULL DEFAULT 1,
    archived_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_archived_at ON projects(archived_at);
CREATE INDEX idx_projects_updated_at ON projects(updated_at DESC);

-- Auto-update trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_projects_updated_at
BEFORE UPDATE ON projects
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_projects_updated_at ON projects;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS projects;
