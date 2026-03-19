package services

import (
	"context"
	"fmt"

	"github.com/bravo68web/country-city-db/internal/cache"
	apperrors "github.com/bravo68web/country-city-db/internal/errors"
	"github.com/bravo68web/country-city-db/internal/models"
	"github.com/bravo68web/country-city-db/internal/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type CountryService struct {
	repo  *repositories.CountryRepository
	redis *redis.Client
}

func NewCountryService(repo *repositories.CountryRepository, redis *redis.Client) *CountryService {
	return &CountryService{repo: repo, redis: redis}
}

func (s *CountryService) GetByID(ctx context.Context, id int64) (*models.Country, error) {
	key := fmt.Sprintf("country:%d", id)
	cached, _ := cache.Get[models.Country](ctx, s.redis, key)
	if cached != nil {
		return cached, nil
	}

	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	cache.Set(ctx, s.redis, key, item, cache.EntityTTL)
	return item, nil
}

func (s *CountryService) GetByISO2(ctx context.Context, code string) (*models.Country, error) {
	key := fmt.Sprintf("country:iso2:%s", code)
	cached, _ := cache.Get[models.Country](ctx, s.redis, key)
	if cached != nil {
		return cached, nil
	}

	item, err := s.repo.GetByISO2(ctx, code)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	cache.Set(ctx, s.redis, key, item, cache.EntityTTL)
	return item, nil
}

func (s *CountryService) GetByISO3(ctx context.Context, code string) (*models.Country, error) {
	key := fmt.Sprintf("country:iso3:%s", code)
	cached, _ := cache.Get[models.Country](ctx, s.redis, key)
	if cached != nil {
		return cached, nil
	}

	item, err := s.repo.GetByISO3(ctx, code)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	cache.Set(ctx, s.redis, key, item, cache.EntityTTL)
	return item, nil
}

func (s *CountryService) GetByName(ctx context.Context, name string) (*models.Country, error) {
	key := fmt.Sprintf("country:name:%s", name)
	cached, _ := cache.Get[models.Country](ctx, s.redis, key)
	if cached != nil {
		return cached, nil
	}

	item, err := s.repo.GetByName(ctx, name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	cache.Set(ctx, s.redis, key, item, cache.EntityTTL)
	return item, nil
}

func (s *CountryService) List(ctx context.Context, params models.QueryParams) (*models.PaginatedResult[models.Country], error) {
	clampParams(&params)
	key := fmt.Sprintf("countries:list:%s:%s:%s:%s:%d:%d:%v", params.Search, params.Name, params.ISO2, params.ISO3, params.Limit, params.Offset, params.NoPage)
	cached, _ := cache.Get[models.PaginatedResult[models.Country]](ctx, s.redis, key)
	if cached != nil {
		return cached, nil
	}

	data, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []models.Country{}
	}

	result := &models.PaginatedResult[models.Country]{Data: data, Total: total}
	if !params.NoPage {
		result.Limit = params.Limit
		result.Offset = params.Offset
	}

	cache.Set(ctx, s.redis, key, result, cache.ListTTL)
	return result, nil
}

func (s *CountryService) ListByRegionID(ctx context.Context, regionID int64, params models.QueryParams) (*models.PaginatedResult[models.Country], error) {
	clampParams(&params)
	key := fmt.Sprintf("countries:region:%d:%s:%s:%d:%d:%v", regionID, params.Search, params.Name, params.Limit, params.Offset, params.NoPage)
	cached, _ := cache.Get[models.PaginatedResult[models.Country]](ctx, s.redis, key)
	if cached != nil {
		return cached, nil
	}

	data, total, err := s.repo.ListByRegionID(ctx, regionID, params)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []models.Country{}
	}

	result := &models.PaginatedResult[models.Country]{Data: data, Total: total}
	if !params.NoPage {
		result.Limit = params.Limit
		result.Offset = params.Offset
	}

	cache.Set(ctx, s.redis, key, result, cache.ListTTL)
	return result, nil
}

func (s *CountryService) ListBySubregionID(ctx context.Context, subregionID int64, params models.QueryParams) (*models.PaginatedResult[models.Country], error) {
	clampParams(&params)
	key := fmt.Sprintf("countries:subregion:%d:%s:%s:%d:%d:%v", subregionID, params.Search, params.Name, params.Limit, params.Offset, params.NoPage)
	cached, _ := cache.Get[models.PaginatedResult[models.Country]](ctx, s.redis, key)
	if cached != nil {
		return cached, nil
	}

	data, total, err := s.repo.ListBySubregionID(ctx, subregionID, params)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []models.Country{}
	}

	result := &models.PaginatedResult[models.Country]{Data: data, Total: total}
	if !params.NoPage {
		result.Limit = params.Limit
		result.Offset = params.Offset
	}

	cache.Set(ctx, s.redis, key, result, cache.ListTTL)
	return result, nil
}
