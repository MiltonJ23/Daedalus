-- +migrate Up
CREATE TABLE IF NOT EXISTS goals (
    id UUID PRIMARY KEY,
    user_id VARCHAR(120) NOT NULL,
    project_id VARCHAR(120) NOT NULL DEFAULT '',
    description TEXT NOT NULL,
    status VARCHAR(40) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_goals_user ON goals(user_id);

CREATE TABLE IF NOT EXISTS sub_tasks (
    id UUID PRIMARY KEY,
    goal_id UUID NOT NULL REFERENCES goals(id) ON DELETE CASCADE,
    type VARCHAR(80) NOT NULL,
    status VARCHAR(40) NOT NULL DEFAULT 'pending',
    payload JSONB NOT NULL DEFAULT '{}',
    depends_on JSONB NOT NULL DEFAULT '[]',
    stream_id VARCHAR(120) NOT NULL DEFAULT '',
    error TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_sub_tasks_goal ON sub_tasks(goal_id);
CREATE INDEX IF NOT EXISTS idx_sub_tasks_status ON sub_tasks(status);

-- +migrate Down
DROP TABLE IF EXISTS sub_tasks;
DROP TABLE IF EXISTS goals;
