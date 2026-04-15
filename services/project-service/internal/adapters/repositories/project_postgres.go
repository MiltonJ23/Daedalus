package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Daedalus/project-service/internal/core/domain"
)

// DBTX abstracts pgx connection vs transaction (Kliops pattern).
type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
}

// ProjectPostgres implements ports.ProjectRepository with pgx.
type ProjectPostgres struct {
	DB DBTX
}

func NewProjectPostgres(db DBTX) *ProjectPostgres {
	return &ProjectPostgres{DB: db}
}

func (r *ProjectPostgres) Create(ctx context.Context, p domain.Project) (domain.Project, error) {
	query := `
		INSERT INTO projects (id, name, industry_type, location, budget,
			floor_width, floor_depth, target_capacity, status, version,
			archived_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id, created_at, updated_at`

	err := r.DB.QueryRow(ctx, query,
		p.ID, p.Name, p.IndustryType, p.Location, p.Budget,
		p.FloorWidth, p.FloorDepth, p.TargetCapacity, p.Status, p.Version,
		p.ArchivedAt, p.CreatedAt, p.UpdatedAt,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return domain.Project{}, fmt.Errorf("failed to insert project: %w", err)
	}
	return p, nil
}

func (r *ProjectPostgres) GetByID(ctx context.Context, id string) (*domain.Project, error) {
	query := `
		SELECT id, name, industry_type, location, budget,
			floor_width, floor_depth, target_capacity, status, version,
			archived_at, created_at, updated_at
		FROM projects WHERE id = $1`

	p := &domain.Project{}
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.IndustryType, &p.Location, &p.Budget,
		&p.FloorWidth, &p.FloorDepth, &p.TargetCapacity, &p.Status, &p.Version,
		&p.ArchivedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project %s: %w", id, err)
	}
	return p, nil
}

func (r *ProjectPostgres) ListByStatus(ctx context.Context, status string) ([]domain.Project, error) {
	var query string
	switch status {
	case "archived":
		query = `SELECT id, name, industry_type, location, budget,
			floor_width, floor_depth, target_capacity, status, version,
			archived_at, created_at, updated_at
			FROM projects WHERE archived_at IS NOT NULL ORDER BY updated_at DESC`
	case "all":
		query = `SELECT id, name, industry_type, location, budget,
			floor_width, floor_depth, target_capacity, status, version,
			archived_at, created_at, updated_at
			FROM projects ORDER BY updated_at DESC`
	default:
		query = `SELECT id, name, industry_type, location, budget,
			floor_width, floor_depth, target_capacity, status, version,
			archived_at, created_at, updated_at
			FROM projects WHERE archived_at IS NULL ORDER BY updated_at DESC`
	}

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(
			&p.ID, &p.Name, &p.IndustryType, &p.Location, &p.Budget,
			&p.FloorWidth, &p.FloorDepth, &p.TargetCapacity, &p.Status, &p.Version,
			&p.ArchivedAt, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan project row: %w", err)
		}
		projects = append(projects, p)
	}

	if projects == nil {
		projects = []domain.Project{}
	}
	return projects, nil
}

func (r *ProjectPostgres) Update(ctx context.Context, p domain.Project) (domain.Project, error) {
	p.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE projects SET
			name=$2, industry_type=$3, location=$4, budget=$5,
			floor_width=$6, floor_depth=$7, target_capacity=$8,
			status=$9, version=$10, archived_at=$11, updated_at=$12
		WHERE id=$1`

	tag, err := r.DB.Exec(ctx, query,
		p.ID, p.Name, p.IndustryType, p.Location, p.Budget,
		p.FloorWidth, p.FloorDepth, p.TargetCapacity,
		p.Status, p.Version, p.ArchivedAt, p.UpdatedAt,
	)
	if err != nil {
		return domain.Project{}, fmt.Errorf("failed to update project %s: %w", p.ID, err)
	}
	if tag.RowsAffected() == 0 {
		return domain.Project{}, fmt.Errorf("project not found: %s", p.ID)
	}
	return p, nil
}

func (r *ProjectPostgres) HardDelete(ctx context.Context, id string) error {
	query := `DELETE FROM projects WHERE id = $1`
	tag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project %s: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("project not found: %s", id)
	}
	return nil
}
