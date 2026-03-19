package repositories

import (
	"context"
	"fmt"

	"github.com/bravo68web/country-city-db/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RegionRepository struct {
	pool *pgxpool.Pool
}

func NewRegionRepository(pool *pgxpool.Pool) *RegionRepository {
	return &RegionRepository{pool: pool}
}

func (r *RegionRepository) GetByID(ctx context.Context, id int64) (*models.Region, error) {
	var region models.Region
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, translations, created_at, updated_at, flag, "wikiDataId" FROM regions WHERE id = $1`, id).
		Scan(&region.ID, &region.Name, &region.Translations, &region.CreatedAt, &region.UpdatedAt, &region.Flag, &region.WikiDataID)
	if err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *RegionRepository) List(ctx context.Context, params models.QueryParams) ([]models.Region, int64, error) {
	where := "WHERE 1=1"
	args := []any{}
	argIdx := 1

	if params.Search != "" {
		where += fmt.Sprintf(" AND name ILIKE $%d", argIdx)
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}
	if params.Name != "" {
		where += fmt.Sprintf(" AND name ILIKE $%d", argIdx)
		args = append(args, "%"+params.Name+"%")
		argIdx++
	}

	var total int64
	err := r.pool.QueryRow(ctx, "SELECT count(*) FROM regions "+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, name, translations, created_at, updated_at, flag, "wikiDataId" FROM regions ` + where + " ORDER BY id"
	if !params.NoPage {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
		args = append(args, params.Limit, params.Offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var regions []models.Region
	for rows.Next() {
		var region models.Region
		if err := rows.Scan(&region.ID, &region.Name, &region.Translations, &region.CreatedAt, &region.UpdatedAt, &region.Flag, &region.WikiDataID); err != nil {
			return nil, 0, err
		}
		regions = append(regions, region)
	}
	return regions, total, rows.Err()
}
