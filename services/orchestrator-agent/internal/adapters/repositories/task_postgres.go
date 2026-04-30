package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Daedalus/orchestrator-agent/internal/core/domain"
)

type TaskPostgres struct{ DB *pgxpool.Pool }

func NewTaskPostgres(db *pgxpool.Pool) *TaskPostgres { return &TaskPostgres{DB: db} }

func (r *TaskPostgres) BulkCreate(ctx context.Context, tasks []domain.SubTask) error {
	const q = `INSERT INTO sub_tasks (id, goal_id, type, status, payload, depends_on, stream_id, error, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	for _, t := range tasks {
		payload, err := json.Marshal(t.Payload)
		if err != nil {
			return err
		}
		deps, err := json.Marshal(t.DependsOn)
		if err != nil {
			return err
		}
		if _, err := r.DB.Exec(ctx, q, t.ID, t.GoalID, t.Type, t.Status, payload, deps, t.StreamID, t.Error, t.CreatedAt, t.UpdatedAt); err != nil {
			return fmt.Errorf("insert task %s: %w", t.ID, err)
		}
	}
	return nil
}

func (r *TaskPostgres) GetByID(ctx context.Context, id string) (*domain.SubTask, error) {
	const q = `SELECT id, goal_id, type, status, payload, depends_on, stream_id, error, created_at, updated_at FROM sub_tasks WHERE id=$1`
	t, err := scanTask(r.DB.QueryRow(ctx, q, id))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskPostgres) ListByGoal(ctx context.Context, goalID string) ([]domain.SubTask, error) {
	const q = `SELECT id, goal_id, type, status, payload, depends_on, stream_id, error, created_at, updated_at
		FROM sub_tasks WHERE goal_id=$1 ORDER BY created_at`
	rows, err := r.DB.Query(ctx, q, goalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.SubTask{}
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func (r *TaskPostgres) Update(ctx context.Context, t domain.SubTask) (domain.SubTask, error) {
	t.UpdatedAt = time.Now().UTC()
	payload, err := json.Marshal(t.Payload)
	if err != nil {
		return domain.SubTask{}, err
	}
	const q = `UPDATE sub_tasks SET status=$2, payload=$3, stream_id=$4, error=$5, updated_at=$6 WHERE id=$1`
	tag, err := r.DB.Exec(ctx, q, t.ID, t.Status, payload, t.StreamID, t.Error, t.UpdatedAt)
	if err != nil {
		return domain.SubTask{}, err
	}
	if tag.RowsAffected() == 0 {
		return domain.SubTask{}, fmt.Errorf("task not found: %s", t.ID)
	}
	return t, nil
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func scanTask(row scannable) (domain.SubTask, error) {
	var (
		t       domain.SubTask
		payload []byte
		deps    []byte
	)
	if err := row.Scan(&t.ID, &t.GoalID, &t.Type, &t.Status, &payload, &deps, &t.StreamID, &t.Error, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return t, err
	}
	if len(payload) > 0 {
		_ = json.Unmarshal(payload, &t.Payload)
	}
	if t.Payload == nil {
		t.Payload = map[string]interface{}{}
	}
	if len(deps) > 0 {
		_ = json.Unmarshal(deps, &t.DependsOn)
	}
	if t.DependsOn == nil {
		t.DependsOn = []string{}
	}
	return t, nil
}
