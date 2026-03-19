package models

import "time"

type Region struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Translations *string    `json:"translations,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Flag         int16      `json:"flag"`
	WikiDataID   *string    `json:"wiki_data_id,omitempty"`
}

type Subregion struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Translations *string    `json:"translations,omitempty"`
	RegionID     int64      `json:"region_id"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Flag         int16      `json:"flag"`
	WikiDataID   *string    `json:"wiki_data_id,omitempty"`
}

type Country struct {
	ID               int64      `json:"id"`
	Name             string     `json:"name"`
	ISO3             *string    `json:"iso3,omitempty"`
	NumericCode      *string    `json:"numeric_code,omitempty"`
	ISO2             *string    `json:"iso2,omitempty"`
	PhoneCode        *string    `json:"phonecode,omitempty"`
	Capital          *string    `json:"capital,omitempty"`
	Currency         *string    `json:"currency,omitempty"`
	CurrencyName     *string    `json:"currency_name,omitempty"`
	CurrencySymbol   *string    `json:"currency_symbol,omitempty"`
	TLD              *string    `json:"tld,omitempty"`
	Native           *string    `json:"native,omitempty"`
	Population       *int64     `json:"population,omitempty"`
	GDP              *int64     `json:"gdp,omitempty"`
	Region           *string    `json:"region,omitempty"`
	RegionID         *int64     `json:"region_id,omitempty"`
	Subregion        *string    `json:"subregion,omitempty"`
	SubregionID      *int64     `json:"subregion_id,omitempty"`
	Nationality      *string    `json:"nationality,omitempty"`
	AreaSqKm         *float64   `json:"area_sq_km,omitempty"`
	PostalCodeFormat *string    `json:"postal_code_format,omitempty"`
	PostalCodeRegex  *string    `json:"postal_code_regex,omitempty"`
	Timezones        *string    `json:"timezones,omitempty"`
	Translations     *string    `json:"translations,omitempty"`
	Latitude         *float64   `json:"latitude,omitempty"`
	Longitude        *float64   `json:"longitude,omitempty"`
	Emoji            *string    `json:"emoji,omitempty"`
	EmojiU           *string    `json:"emoji_u,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time  `json:"updated_at"`
	Flag             int16      `json:"flag"`
	WikiDataID       *string    `json:"wiki_data_id,omitempty"`
}

type State struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	CountryID   int64      `json:"country_id"`
	CountryCode string     `json:"country_code"`
	FipsCode    *string    `json:"fips_code,omitempty"`
	ISO2        *string    `json:"iso2,omitempty"`
	ISO3166_2   *string    `json:"iso3166_2,omitempty"`
	Type        *string    `json:"type,omitempty"`
	Level       *int       `json:"level,omitempty"`
	ParentID    *int64     `json:"parent_id,omitempty"`
	Native      *string    `json:"native,omitempty"`
	Latitude    *float64   `json:"latitude,omitempty"`
	Longitude   *float64   `json:"longitude,omitempty"`
	Timezone    *string    `json:"timezone,omitempty"`
	Translations *string   `json:"translations,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Flag        int16      `json:"flag"`
	WikiDataID  *string    `json:"wiki_data_id,omitempty"`
	Population  *string    `json:"population,omitempty"`
}

type City struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	StateID     int64      `json:"state_id"`
	StateCode   string     `json:"state_code"`
	CountryID   int64      `json:"country_id"`
	CountryCode string     `json:"country_code"`
	Type        *string    `json:"type,omitempty"`
	Level       *int       `json:"level,omitempty"`
	ParentID    *int64     `json:"parent_id,omitempty"`
	Latitude    float64    `json:"latitude"`
	Longitude   float64    `json:"longitude"`
	Native      *string    `json:"native,omitempty"`
	Population  *int64     `json:"population,omitempty"`
	Timezone    *string    `json:"timezone,omitempty"`
	Translations *string   `json:"translations,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Flag        int16      `json:"flag"`
	WikiDataID  *string    `json:"wiki_data_id,omitempty"`
}

type QueryParams struct {
	Search string `form:"search" json:"search"`
	Name   string `form:"name" json:"name"`
	ISO2   string `form:"iso2" json:"iso2"`
	ISO3   string `form:"iso3" json:"iso3"`
	Limit  int    `form:"limit,default=20" json:"limit"`
	Offset int    `form:"offset,default=0" json:"offset"`
	NoPage bool   `form:"no_page" json:"no_page"`
}

type PaginatedResult[T any] struct {
	Data   []T   `json:"data"`
	Total  int64 `json:"total"`
	Limit  int   `json:"limit,omitempty"`
	Offset int   `json:"offset,omitempty"`
}

type StatsResponse struct {
	Database DatabaseStats `json:"database"`
	Cache    CacheStats    `json:"cache"`
}

type DatabaseStats struct {
	Regions    int64 `json:"regions"`
	Subregions int64 `json:"subregions"`
	Countries  int64 `json:"countries"`
	States     int64 `json:"states"`
	Cities     int64 `json:"cities"`
}

type CacheStats struct {
	Keys      int64  `json:"keys"`
	Memory    string `json:"memory_used"`
	Connected bool   `json:"connected"`
}
