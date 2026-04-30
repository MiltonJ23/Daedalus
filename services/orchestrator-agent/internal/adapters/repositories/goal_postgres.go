package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Daedalus/orchestrator-agent/internal/core/domain"
)

type GoalPostgres struct{ DB *pgxpool.Pool }

func NewGoalPostgres(db *pgxpool.Pool) *GoalPostgres { return &GoalPostgres{DB: db} }

func (r *GoalPostgres) Create(ctx context.Context, g domain.Goal) (domain.Goal, error) {
	const q = `INSERT INTO goals (id, user_id, project_id, description, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, created_at, updated_at`
	if err := r.DB.QueryRow(ctx, q, g.ID, g.UserID, g.ProjectID, g.Description, g.Status, g.CreatedAt, g.UpdatedAt).
		Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt); err != nil {
		return domain.Goal{}, fmt.Errorf("insert goal: %w", err)
	}
	return g, nil
}

func (r *GoalPostgres) GetByID(ctx context.Context, id string) (*domain.Goal, error) {
	const q = `SELECT id, user_id, project_id, description, status, created_at, updated_at FROM goals WHERE id=$1`
	var g domain.Goal
	err := r.DB.QueryRow(ctx, q, id).Scan(&g.ID, &g.UserID, &g.ProjectID, &g.Description, &g.Status, &g.CreatedAt, &g.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GoalPostgres) List(ctx context.Context, userID string) ([]domain.Goal, error) {
	q := `SELECT id, user_id, project_id, description, status, created_at, updated_at FROM goals`
	args := []interface{}{}
	if userID != "" {
		q += ` WHERE user_id=$1`
		args = append(args, userID)
	}
	q += ` ORDER BY created_at DESC`
	rows, err := r.DB.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Goal{}
	for rows.Next() {
		var g domain.Goal
		if err := rows.Scan(&g.ID, &g.UserID, &g.ProjectID, &g.Description, &g.Status, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, nil
}

func (r *GoalPostgres) Update(ctx context.Context, g domain.Goal) (domain.Goal, error) {
	g.UpdatedAt = time.Now().UTC()
	const q = `UPDATE goals SET status=$2, updated_at=$3 WHERE id=$1`
	tag, err := r.DB.Exec(ctx, q, g.ID, g.Status, g.UpdatedAt)
	if err != nil {
		return domain.Goal{}, err
	}
	if tag.RowsAffected() == 0 {
		return domain.Goal{}, fmt.Errorf("goal not found: %s", g.ID)
	}
	return g, nil
}
