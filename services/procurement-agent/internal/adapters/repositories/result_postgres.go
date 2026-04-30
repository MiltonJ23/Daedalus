package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Daedalus/procurement-agent/internal/core/domain"
)

// ResultPostgres implements ports.ResultRepository with pgx.
type ResultPostgres struct {
	DB DBTX
}

func NewResultPostgres(db DBTX) *ResultPostgres {
	return &ResultPostgres{DB: db}
}

func (r *ResultPostgres) BulkCreate(ctx context.Context, results []domain.EquipmentResult) error {
	if len(results) == 0 {
		return nil
	}

	const query = `
		INSERT INTO equipment_results
			(id, search_id, name, model, supplier, supplier_rating,
			 price_usd, price_xaf, lead_time_days, spec_match, score,
			 specifications, width_m, depth_m, height_m, power_kw,
			 country, decision, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)`

	for _, item := range results {
		specs, err := json.Marshal(item.Specifications)
		if err != nil {
			return fmt.Errorf("failed to marshal specifications: %w", err)
		}
		if _, err := r.DB.Exec(ctx, query,
			item.ID, item.SearchID, item.Name, item.Model, item.Supplier, item.SupplierRating,
			item.PriceUSD, item.PriceXAF, item.LeadTimeDays, item.SpecMatch, item.Score,
			specs, item.Dimensions.WidthM, item.Dimensions.DepthM, item.Dimensions.HeightM,
			item.PowerKW, item.Country, item.Decision, item.CreatedAt, item.UpdatedAt,
		); err != nil {
			return fmt.Errorf("failed to insert result %s: %w", item.ID, err)
		}
	}
	return nil
}

func (r *ResultPostgres) GetByID(ctx context.Context, id string) (*domain.EquipmentResult, error) {
	const query = `
		SELECT id, search_id, name, model, supplier, supplier_rating,
			price_usd, price_xaf, lead_time_days, spec_match, score,
			specifications, width_m, depth_m, height_m, power_kw,
			country, decision, created_at, updated_at
		FROM equipment_results WHERE id = $1`

	row := r.DB.QueryRow(ctx, query, id)
	res, err := scanResult(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get result %s: %w", id, err)
	}
	return &res, nil
}

func (r *ResultPostgres) ListBySearch(ctx context.Context, searchID string) ([]domain.EquipmentResult, error) {
	const query = `
		SELECT id, search_id, name, model, supplier, supplier_rating,
			price_usd, price_xaf, lead_time_days, spec_match, score,
			specifications, width_m, depth_m, height_m, power_kw,
			country, decision, created_at, updated_at
		FROM equipment_results WHERE search_id = $1 ORDER BY score DESC`

	rows, err := r.DB.Query(ctx, query, searchID)
	if err != nil {
		return nil, fmt.Errorf("failed to list results: %w", err)
	}
	defer rows.Close()

	var out []domain.EquipmentResult
	for rows.Next() {
		res, err := scanResult(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, res)
	}
	if out == nil {
		out = []domain.EquipmentResult{}
	}
	return out, nil
}

func (r *ResultPostgres) Update(ctx context.Context, item domain.EquipmentResult) (domain.EquipmentResult, error) {
	item.UpdatedAt = time.Now().UTC()

	specs, err := json.Marshal(item.Specifications)
	if err != nil {
		return domain.EquipmentResult{}, fmt.Errorf("failed to marshal specifications: %w", err)
	}

	const query = `
		UPDATE equipment_results SET
			name=$2, model=$3, supplier=$4, supplier_rating=$5,
			price_usd=$6, price_xaf=$7, lead_time_days=$8, spec_match=$9, score=$10,
			specifications=$11, width_m=$12, depth_m=$13, height_m=$14, power_kw=$15,
			country=$16, decision=$17, updated_at=$18
		WHERE id=$1`

	tag, err := r.DB.Exec(ctx, query,
		item.ID, item.Name, item.Model, item.Supplier, item.SupplierRating,
		item.PriceUSD, item.PriceXAF, item.LeadTimeDays, item.SpecMatch, item.Score,
		specs, item.Dimensions.WidthM, item.Dimensions.DepthM, item.Dimensions.HeightM,
		item.PowerKW, item.Country, item.Decision, item.UpdatedAt,
	)
	if err != nil {
		return domain.EquipmentResult{}, fmt.Errorf("failed to update result %s: %w", item.ID, err)
	}
	if tag.RowsAffected() == 0 {
		return domain.EquipmentResult{}, fmt.Errorf("result not found: %s", item.ID)
	}
	return item, nil
}

func scanResult(row scannable) (domain.EquipmentResult, error) {
	var (
		res   domain.EquipmentResult
		specs []byte
	)
	err := row.Scan(
		&res.ID, &res.SearchID, &res.Name, &res.Model, &res.Supplier, &res.SupplierRating,
		&res.PriceUSD, &res.PriceXAF, &res.LeadTimeDays, &res.SpecMatch, &res.Score,
		&specs, &res.Dimensions.WidthM, &res.Dimensions.DepthM, &res.Dimensions.HeightM,
		&res.PowerKW, &res.Country, &res.Decision, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		return res, err
	}
	if len(specs) > 0 {
		if err := json.Unmarshal(specs, &res.Specifications); err != nil {
			return res, fmt.Errorf("failed to unmarshal specifications: %w", err)
		}
	} else {
		res.Specifications = map[string]interface{}{}
	}
	return res, nil
}
