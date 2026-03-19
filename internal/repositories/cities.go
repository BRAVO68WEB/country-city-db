package repositories

import (
	"context"
	"fmt"

	"github.com/bravo68web/country-city-db/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CityRepository struct {
	pool *pgxpool.Pool
}

func NewCityRepository(pool *pgxpool.Pool) *CityRepository {
	return &CityRepository{pool: pool}
}

const cityCols = `id, name, state_id, state_code, country_id, country_code, type, level, parent_id, latitude, longitude, native, population, timezone, translations, created_at, updated_at, flag, "wikiDataId"`

func scanCity(scanner interface{ Scan(dest ...any) error }) (*models.City, error) {
	var c models.City
	err := scanner.Scan(
		&c.ID, &c.Name, &c.StateID, &c.StateCode, &c.CountryID, &c.CountryCode,
		&c.Type, &c.Level, &c.ParentID, &c.Latitude, &c.Longitude, &c.Native,
		&c.Population, &c.Timezone, &c.Translations, &c.CreatedAt, &c.UpdatedAt,
		&c.Flag, &c.WikiDataID,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CityRepository) GetByID(ctx context.Context, id int64) (*models.City, error) {
	row := r.pool.QueryRow(ctx, "SELECT "+cityCols+" FROM cities WHERE id = $1", id)
	return scanCity(row)
}

func (r *CityRepository) List(ctx context.Context, params models.QueryParams) ([]models.City, int64, error) {
	return r.listWithFilter(ctx, params, "", nil)
}

func (r *CityRepository) ListByCountryID(ctx context.Context, countryID int64, params models.QueryParams) ([]models.City, int64, error) {
	return r.listWithFilter(ctx, params, "country_id", countryID)
}

func (r *CityRepository) ListByStateID(ctx context.Context, stateID int64, params models.QueryParams) ([]models.City, int64, error) {
	return r.listWithFilter(ctx, params, "state_id", stateID)
}

func (r *CityRepository) listWithFilter(ctx context.Context, params models.QueryParams, filterCol string, filterVal any) ([]models.City, int64, error) {
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
	err := r.pool.QueryRow(ctx, "SELECT count(*) FROM cities "+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := "SELECT " + cityCols + " FROM cities " + where + " ORDER BY id"
	if !params.NoPage {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
		args = append(args, params.Limit, params.Offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []models.City
	for rows.Next() {
		c, err := scanCity(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *c)
	}
	return items, total, rows.Err()
}
