package repositories

import (
	"context"
	"fmt"

	"github.com/bravo68web/country-city-db/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StateRepository struct {
	pool *pgxpool.Pool
}

func NewStateRepository(pool *pgxpool.Pool) *StateRepository {
	return &StateRepository{pool: pool}
}

const stateCols = `id, name, country_id, country_code, fips_code, iso2, iso3166_2, type, level, parent_id, native, latitude, longitude, timezone, translations, created_at, updated_at, flag, "wikiDataId", population`

func scanState(scanner interface{ Scan(dest ...any) error }) (*models.State, error) {
	var s models.State
	err := scanner.Scan(
		&s.ID, &s.Name, &s.CountryID, &s.CountryCode, &s.FipsCode, &s.ISO2,
		&s.ISO3166_2, &s.Type, &s.Level, &s.ParentID, &s.Native, &s.Latitude,
		&s.Longitude, &s.Timezone, &s.Translations, &s.CreatedAt, &s.UpdatedAt,
		&s.Flag, &s.WikiDataID, &s.Population,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StateRepository) GetByID(ctx context.Context, id int64) (*models.State, error) {
	row := r.pool.QueryRow(ctx, "SELECT "+stateCols+" FROM states WHERE id = $1", id)
	return scanState(row)
}

func (r *StateRepository) List(ctx context.Context, params models.QueryParams) ([]models.State, int64, error) {
	return r.listWithFilter(ctx, params, "", nil)
}

func (r *StateRepository) ListByCountryID(ctx context.Context, countryID int64, params models.QueryParams) ([]models.State, int64, error) {
	return r.listWithFilter(ctx, params, "country_id", countryID)
}

func (r *StateRepository) listWithFilter(ctx context.Context, params models.QueryParams, filterCol string, filterVal any) ([]models.State, int64, error) {
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
	err := r.pool.QueryRow(ctx, "SELECT count(*) FROM states "+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := "SELECT " + stateCols + " FROM states " + where + " ORDER BY id"
	if !params.NoPage {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
		args = append(args, params.Limit, params.Offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []models.State
	for rows.Next() {
		s, err := scanState(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *s)
	}
	return items, total, rows.Err()
}
