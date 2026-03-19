package repositories

import (
	"context"
	"fmt"

	"github.com/bravo68web/country-city-db/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubregionRepository struct {
	pool *pgxpool.Pool
}

func NewSubregionRepository(pool *pgxpool.Pool) *SubregionRepository {
	return &SubregionRepository{pool: pool}
}

func (r *SubregionRepository) GetByID(ctx context.Context, id int64) (*models.Subregion, error) {
	var s models.Subregion
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, translations, region_id, created_at, updated_at, flag, "wikiDataId" FROM subregions WHERE id = $1`, id).
		Scan(&s.ID, &s.Name, &s.Translations, &s.RegionID, &s.CreatedAt, &s.UpdatedAt, &s.Flag, &s.WikiDataID)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SubregionRepository) List(ctx context.Context, params models.QueryParams) ([]models.Subregion, int64, error) {
	return r.listWithFilter(ctx, params, "", nil)
}

func (r *SubregionRepository) ListByRegionID(ctx context.Context, regionID int64, params models.QueryParams) ([]models.Subregion, int64, error) {
	return r.listWithFilter(ctx, params, "region_id", regionID)
}

func (r *SubregionRepository) listWithFilter(ctx context.Context, params models.QueryParams, filterCol string, filterVal any) ([]models.Subregion, int64, error) {
	where := "WHERE 1=1"
	args := []any{}
	argIdx := 1

	if filterCol != "" {
		where += fmt.Sprintf(" AND %s = $%d", filterCol, argIdx)
		args = append(args, filterVal)
		argIdx++
	}
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
	err := r.pool.QueryRow(ctx, "SELECT count(*) FROM subregions "+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, name, translations, region_id, created_at, updated_at, flag, "wikiDataId" FROM subregions ` + where + " ORDER BY id"
	if !params.NoPage {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
		args = append(args, params.Limit, params.Offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []models.Subregion
	for rows.Next() {
		var s models.Subregion
		if err := rows.Scan(&s.ID, &s.Name, &s.Translations, &s.RegionID, &s.CreatedAt, &s.UpdatedAt, &s.Flag, &s.WikiDataID); err != nil {
			return nil, 0, err
		}
		items = append(items, s)
	}
	return items, total, rows.Err()
}
