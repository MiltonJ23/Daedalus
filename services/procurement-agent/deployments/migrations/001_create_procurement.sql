-- +migrate Up
CREATE TABLE procurement_searches (
    id              VARCHAR(36) PRIMARY KEY,
    project_id      VARCHAR(36) NOT NULL,
    query           TEXT NOT NULL,
    category        VARCHAR(120),
    max_budget_usd  DOUBLE PRECISION,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    cache_key       VARCHAR(64) NOT NULL,
    extracted_spec  JSONB NOT NULL DEFAULT '{}'::jsonb,
    expires_at      TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_procurement_searches_project ON procurement_searches(project_id);
CREATE INDEX idx_procurement_searches_cache_key ON procurement_searches(cache_key);
CREATE INDEX idx_procurement_searches_expires_at ON procurement_searches(expires_at);

CREATE TABLE equipment_results (
    id                  VARCHAR(36) PRIMARY KEY,
    search_id           VARCHAR(36) NOT NULL REFERENCES procurement_searches(id) ON DELETE CASCADE,
    name                VARCHAR(255) NOT NULL,
    model               VARCHAR(255) NOT NULL,
    supplier            VARCHAR(255) NOT NULL,
    supplier_rating     DOUBLE PRECISION NOT NULL DEFAULT 0,
    price_usd           DOUBLE PRECISION NOT NULL,
    price_xaf           DOUBLE PRECISION NOT NULL,
    lead_time_days      INTEGER NOT NULL DEFAULT 0,
    spec_match          DOUBLE PRECISION NOT NULL DEFAULT 0,
    score               DOUBLE PRECISION NOT NULL DEFAULT 0,
    specifications      JSONB NOT NULL DEFAULT '{}'::jsonb,
    width_m             DOUBLE PRECISION NOT NULL DEFAULT 0,
    depth_m             DOUBLE PRECISION NOT NULL DEFAULT 0,
    height_m            DOUBLE PRECISION NOT NULL DEFAULT 0,
    power_kw            DOUBLE PRECISION NOT NULL DEFAULT 0,
    country             VARCHAR(80) NOT NULL DEFAULT '',
    decision            VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_equipment_results_search ON equipment_results(search_id);
CREATE INDEX idx_equipment_results_score ON equipment_results(search_id, score DESC);
CREATE INDEX idx_equipment_results_country ON equipment_results(search_id, country);
CREATE INDEX idx_equipment_results_price ON equipment_results(search_id, price_usd);

CREATE OR REPLACE FUNCTION update_proc_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_procurement_searches_updated_at
BEFORE UPDATE ON procurement_searches
FOR EACH ROW EXECUTE FUNCTION update_proc_updated_at_column();

CREATE TRIGGER update_equipment_results_updated_at
BEFORE UPDATE ON equipment_results
FOR EACH ROW EXECUTE FUNCTION update_proc_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_equipment_results_updated_at ON equipment_results;
DROP TRIGGER IF EXISTS update_procurement_searches_updated_at ON procurement_searches;
DROP FUNCTION IF EXISTS update_proc_updated_at_column();
DROP TABLE IF EXISTS equipment_results;
DROP TABLE IF EXISTS procurement_searches;
