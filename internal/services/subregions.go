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

type SubregionService struct {
	repo  *repositories.SubregionRepository
	redis *redis.Client
}

func NewSubregionService(repo *repositories.SubregionRepository, redis *redis.Client) *SubregionService {
	return &SubregionService{repo: repo, redis: redis}
}

func (s *SubregionService) GetByID(ctx context.Context, id int64) (*models.Subregion, error) {
	key := fmt.Sprintf("subregion:%d", id)
	cached, _ := cache.Get[models.Subregion](ctx, s.redis, key)
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

func (s *SubregionService) List(ctx context.Context, params models.QueryParams) (*models.PaginatedResult[models.Subregion], error) {
	clampParams(&params)
	key := fmt.Sprintf("subregions:list:%s:%s:%d:%d:%v", params.Search, params.Name, params.Limit, params.Offset, params.NoPage)
	cached, _ := cache.Get[models.PaginatedResult[models.Subregion]](ctx, s.redis, key)
	if cached != nil {
		return cached, nil
	}

	data, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []models.Subregion{}
	}

	result := &models.PaginatedResult[models.Subregion]{Data: data, Total: total}
	if !params.NoPage {
		result.Limit = params.Limit
		result.Offset = params.Offset
	}

	cache.Set(ctx, s.redis, key, result, cache.ListTTL)
	return result, nil
}

func (s *SubregionService) ListByRegionID(ctx context.Context, regionID int64, params models.QueryParams) (*models.PaginatedResult[models.Subregion], error) {
	clampParams(&params)
	key := fmt.Sprintf("subregions:region:%d:%s:%s:%d:%d:%v", regionID, params.Search, params.Name, params.Limit, params.Offset, params.NoPage)
	cached, _ := cache.Get[models.PaginatedResult[models.Subregion]](ctx, s.redis, key)
	if cached != nil {
		return cached, nil
	}

	data, total, err := s.repo.ListByRegionID(ctx, regionID, params)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []models.Subregion{}
	}

	result := &models.PaginatedResult[models.Subregion]{Data: data, Total: total}
	if !params.NoPage {
		result.Limit = params.Limit
		result.Offset = params.Offset
	}

	cache.Set(ctx, s.redis, key, result, cache.ListTTL)
	return result, nil
}
