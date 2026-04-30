package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Daedalus/procurement-agent/internal/core/domain"
)

// DBTX abstracts pgx connection vs transaction (Kliops pattern).
type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
}

// SearchPostgres implements ports.SearchRepository with pgx.
type SearchPostgres struct {
	DB DBTX
}

func NewSearchPostgres(db DBTX) *SearchPostgres {
	return &SearchPostgres{DB: db}
}

func (r *SearchPostgres) Create(ctx context.Context, s domain.ProcurementSearch) (domain.ProcurementSearch, error) {
	specs, err := json.Marshal(s.ExtractedSpec)
	if err != nil {
		return domain.ProcurementSearch{}, fmt.Errorf("failed to marshal extracted_spec: %w", err)
	}
	query := `
		INSERT INTO procurement_searches
			(id, project_id, query, category, max_budget_usd, status, cache_key, extracted_spec, expires_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id, created_at, updated_at`

	err = r.DB.QueryRow(ctx, query,
		s.ID, s.ProjectID, s.Query, s.Category, s.MaxBudgetUSD,
		s.Status, s.CacheKey, specs, s.ExpiresAt, s.CreatedAt, s.UpdatedAt,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return domain.ProcurementSearch{}, fmt.Errorf("failed to insert search: %w", err)
	}
	return s, nil
}

func (r *SearchPostgres) GetByID(ctx context.Context, id string) (*domain.ProcurementSearch, error) {
	return r.fetchOne(ctx, `WHERE id = $1`, id)
}

func (r *SearchPostgres) GetByCacheKey(ctx context.Context, cacheKey string) (*domain.ProcurementSearch, error) {
	return r.fetchOne(ctx, `WHERE cache_key = $1 AND expires_at > NOW() ORDER BY created_at DESC LIMIT 1`, cacheKey)
}

func (r *SearchPostgres) ListByProject(ctx context.Context, projectID string) ([]domain.ProcurementSearch, error) {
	const base = `
		SELECT id, project_id, query, category, max_budget_usd, status,
			cache_key, extracted_spec, expires_at, created_at, updated_at
		FROM procurement_searches`

	var (
		rows pgx.Rows
		err  error
	)
	if projectID == "" {
		rows, err = r.DB.Query(ctx, base+` ORDER BY created_at DESC`)
	} else {
		rows, err = r.DB.Query(ctx, base+` WHERE project_id = $1 ORDER BY created_at DESC`, projectID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list searches: %w", err)
	}
	defer rows.Close()

	var out []domain.ProcurementSearch
	for rows.Next() {
		s, err := scanSearch(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	if out == nil {
		out = []domain.ProcurementSearch{}
	}
	return out, nil
}

func (r *SearchPostgres) Update(ctx context.Context, s domain.ProcurementSearch) (domain.ProcurementSearch, error) {
	s.UpdatedAt = time.Now().UTC()

	specs, err := json.Marshal(s.ExtractedSpec)
	if err != nil {
		return domain.ProcurementSearch{}, fmt.Errorf("failed to marshal extracted_spec: %w", err)
	}

	query := `
		UPDATE procurement_searches SET
			project_id=$2, query=$3, category=$4, max_budget_usd=$5,
			status=$6, cache_key=$7, extracted_spec=$8, expires_at=$9, updated_at=$10
		WHERE id=$1`

	tag, err := r.DB.Exec(ctx, query,
		s.ID, s.ProjectID, s.Query, s.Category, s.MaxBudgetUSD,
		s.Status, s.CacheKey, specs, s.ExpiresAt, s.UpdatedAt,
	)
	if err != nil {
		return domain.ProcurementSearch{}, fmt.Errorf("failed to update search %s: %w", s.ID, err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ProcurementSearch{}, fmt.Errorf("search not found: %s", s.ID)
	}
	return s, nil
}

func (r *SearchPostgres) fetchOne(ctx context.Context, where string, args ...any) (*domain.ProcurementSearch, error) {
	query := `
		SELECT id, project_id, query, category, max_budget_usd, status,
			cache_key, extracted_spec, expires_at, created_at, updated_at
		FROM procurement_searches ` + where

	row := r.DB.QueryRow(ctx, query, args...)
	s, err := scanSearchRow(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch search: %w", err)
	}
	return &s, nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scanSearch(row scannable) (domain.ProcurementSearch, error) {
	return scanSearchRow(row)
}

func scanSearchRow(row scannable) (domain.ProcurementSearch, error) {
	var (
		s     domain.ProcurementSearch
		specs []byte
	)
	err := row.Scan(
		&s.ID, &s.ProjectID, &s.Query, &s.Category, &s.MaxBudgetUSD,
		&s.Status, &s.CacheKey, &specs, &s.ExpiresAt, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return s, err
	}
	if len(specs) > 0 {
		if err := json.Unmarshal(specs, &s.ExtractedSpec); err != nil {
			return s, fmt.Errorf("failed to unmarshal extracted_spec: %w", err)
		}
	} else {
		s.ExtractedSpec = map[string]interface{}{}
	}
	return s, nil
}
