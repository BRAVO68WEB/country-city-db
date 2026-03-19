package repositories

import (
	"context"
	"fmt"

	"github.com/bravo68web/country-city-db/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CountryRepository struct {
	pool *pgxpool.Pool
}

func NewCountryRepository(pool *pgxpool.Pool) *CountryRepository {
	return &CountryRepository{pool: pool}
}

const countryCols = `id, name, iso3, numeric_code, iso2, phonecode, capital, currency, currency_name, currency_symbol, tld, native, population, gdp, region, region_id, subregion, subregion_id, nationality, area_sq_km, postal_code_format, postal_code_regex, timezones, translations, latitude, longitude, emoji, "emojiU", created_at, updated_at, flag, "wikiDataId"`

func scanCountry(scanner interface{ Scan(dest ...any) error }) (*models.Country, error) {
	var c models.Country
	err := scanner.Scan(
		&c.ID, &c.Name, &c.ISO3, &c.NumericCode, &c.ISO2, &c.PhoneCode,
		&c.Capital, &c.Currency, &c.CurrencyName, &c.CurrencySymbol, &c.TLD,
		&c.Native, &c.Population, &c.GDP, &c.Region, &c.RegionID, &c.Subregion,
		&c.SubregionID, &c.Nationality, &c.AreaSqKm, &c.PostalCodeFormat,
		&c.PostalCodeRegex, &c.Timezones, &c.Translations, &c.Latitude,
		&c.Longitude, &c.Emoji, &c.EmojiU, &c.CreatedAt, &c.UpdatedAt,
		&c.Flag, &c.WikiDataID,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CountryRepository) GetByID(ctx context.Context, id int64) (*models.Country, error) {
	row := r.pool.QueryRow(ctx, "SELECT "+countryCols+" FROM countries WHERE id = $1", id)
	return scanCountry(row)
}

func (r *CountryRepository) GetByISO2(ctx context.Context, code string) (*models.Country, error) {
	row := r.pool.QueryRow(ctx, "SELECT "+countryCols+" FROM countries WHERE UPPER(iso2) = UPPER($1)", code)
	return scanCountry(row)
}

func (r *CountryRepository) GetByISO3(ctx context.Context, code string) (*models.Country, error) {
	row := r.pool.QueryRow(ctx, "SELECT "+countryCols+" FROM countries WHERE UPPER(iso3) = UPPER($1)", code)
	return scanCountry(row)
}

func (r *CountryRepository) GetByName(ctx context.Context, name string) (*models.Country, error) {
	row := r.pool.QueryRow(ctx, "SELECT "+countryCols+" FROM countries WHERE LOWER(name) = LOWER($1)", name)
	return scanCountry(row)
}

func (r *CountryRepository) List(ctx context.Context, params models.QueryParams) ([]models.Country, int64, error) {
	return r.listWithFilter(ctx, params, "", nil)
}

func (r *CountryRepository) ListByRegionID(ctx context.Context, regionID int64, params models.QueryParams) ([]models.Country, int64, error) {
	return r.listWithFilter(ctx, params, "region_id", regionID)
}

func (r *CountryRepository) ListBySubregionID(ctx context.Context, subregionID int64, params models.QueryParams) ([]models.Country, int64, error) {
	return r.listWithFilter(ctx, params, "subregion_id", subregionID)
}

func (r *CountryRepository) listWithFilter(ctx context.Context, params models.QueryParams, filterCol string, filterVal any) ([]models.Country, int64, error) {
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
	if params.ISO2 != "" {
		where += fmt.Sprintf(" AND UPPER(iso2) = UPPER($%d)", argIdx)
		args = append(args, params.ISO2)
		argIdx++
	}
	if params.ISO3 != "" {
		where += fmt.Sprintf(" AND UPPER(iso3) = UPPER($%d)", argIdx)
		args = append(args, params.ISO3)
		argIdx++
	}

	var total int64
	err := r.pool.QueryRow(ctx, "SELECT count(*) FROM countries "+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := "SELECT " + countryCols + " FROM countries " + where + " ORDER BY id"
	if !params.NoPage {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
		args = append(args, params.Limit, params.Offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []models.Country
	for rows.Next() {
		c, err := scanCountry(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *c)
	}
	return items, total, rows.Err()
}
