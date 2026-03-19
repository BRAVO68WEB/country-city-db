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

type CityService struct {
	repo  *repositories.CityRepository
	redis *redis.Client
}

func NewCityService(repo *repositories.CityRepository, redis *redis.Client) *CityService {
	return &CityService{repo: repo, redis: redis}
}

func (s *CityService) GetByID(ctx context.Context, id int64) (*models.City, error) {
	key := fmt.Sprintf("city:%d", id)
	cached, _ := cache.Get[models.City](ctx, s.redis, key)
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

func (s *CityService) List(ctx context.Context, params models.QueryParams) (*models.PaginatedResult[models.City], error) {
	clampParams(&params)
	key := fmt.Sprintf("cities:list:%s:%s:%d:%d:%v", params.Search, params.Name, params.Limit, params.Offset, params.NoPage)
	cached, _ := cache.Get[models.PaginatedResult[models.City]](ctx, s.redis, key)
	if cached != nil {
		return cached, nil
	}

	data, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []models.City{}
	}

	result := &models.PaginatedResult[models.City]{Data: data, Total: total}
	if !params.NoPage {
		result.Limit = params.Limit
		result.Offset = params.Offset
	}

	cache.Set(ctx, s.redis, key, result, cache.ListTTL)
	return result, nil
}

func (s *CityService) ListByCountryID(ctx context.Context, countryID int64, params models.QueryParams) (*models.PaginatedResult[models.City], error) {
	clampParams(&params)
	key := fmt.Sprintf("cities:country:%d:%s:%s:%d:%d:%v", countryID, params.Search, params.Name, params.Limit, params.Offset, params.NoPage)
	cached, _ := cache.Get[models.PaginatedResult[models.City]](ctx, s.redis, key)
	if cached != nil {
		return cached, nil
	}

	data, total, err := s.repo.ListByCountryID(ctx, countryID, params)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []models.City{}
	}

	result := &models.PaginatedResult[models.City]{Data: data, Total: total}
	if !params.NoPage {
		result.Limit = params.Limit
		result.Offset = params.Offset
	}

	cache.Set(ctx, s.redis, key, result, cache.ListTTL)
	return result, nil
}

func (s *CityService) ListByStateID(ctx context.Context, stateID int64, params models.QueryParams) (*models.PaginatedResult[models.City], error) {
	clampParams(&params)
	key := fmt.Sprintf("cities:state:%d:%s:%s:%d:%d:%v", stateID, params.Search, params.Name, params.Limit, params.Offset, params.NoPage)
	cached, _ := cache.Get[models.PaginatedResult[models.City]](ctx, s.redis, key)
	if cached != nil {
		return cached, nil
	}

	data, total, err := s.repo.ListByStateID(ctx, stateID, params)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []models.City{}
	}

	result := &models.PaginatedResult[models.City]{Data: data, Total: total}
	if !params.NoPage {
		result.Limit = params.Limit
		result.Offset = params.Offset
	}

	cache.Set(ctx, s.redis, key, result, cache.ListTTL)
	return result, nil
}
