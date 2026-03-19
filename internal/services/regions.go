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

type RegionService struct {
	repo  *repositories.RegionRepository
	redis *redis.Client
}

func NewRegionService(repo *repositories.RegionRepository, redis *redis.Client) *RegionService {
	return &RegionService{repo: repo, redis: redis}
}

func (s *RegionService) GetByID(ctx context.Context, id int64) (*models.Region, error) {
	key := fmt.Sprintf("region:%d", id)
	cached, _ := cache.Get[models.Region](ctx, s.redis, key)
	if cached != nil {
		return cached, nil
	}

	region, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	cache.Set(ctx, s.redis, key, region, cache.EntityTTL)
	return region, nil
}

func (s *RegionService) List(ctx context.Context, params models.QueryParams) (*models.PaginatedResult[models.Region], error) {
	clampParams(&params)
	key := fmt.Sprintf("regions:list:%s:%s:%d:%d:%v", params.Search, params.Name, params.Limit, params.Offset, params.NoPage)
	cached, _ := cache.Get[models.PaginatedResult[models.Region]](ctx, s.redis, key)
	if cached != nil {
		return cached, nil
	}

	data, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []models.Region{}
	}

	result := &models.PaginatedResult[models.Region]{Data: data, Total: total}
	if !params.NoPage {
		result.Limit = params.Limit
		result.Offset = params.Offset
	}

	cache.Set(ctx, s.redis, key, result, cache.ListTTL)
	return result, nil
}
