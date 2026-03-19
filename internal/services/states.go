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

type StateService struct {
	repo  *repositories.StateRepository
	redis *redis.Client
}

func NewStateService(repo *repositories.StateRepository, redis *redis.Client) *StateService {
	return &StateService{repo: repo, redis: redis}
}

func (s *StateService) GetByID(ctx context.Context, id int64) (*models.State, error) {
	key := fmt.Sprintf("state:%d", id)
	cached, _ := cache.Get[models.State](ctx, s.redis, key)
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

func (s *StateService) List(ctx context.Context, params models.QueryParams) (*models.PaginatedResult[models.State], error) {
	clampParams(&params)
	key := fmt.Sprintf("states:list:%s:%s:%d:%d:%v", params.Search, params.Name, params.Limit, params.Offset, params.NoPage)
	cached, _ := cache.Get[models.PaginatedResult[models.State]](ctx, s.redis, key)
	if cached != nil {
		return cached, nil
	}

	data, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []models.State{}
	}

	result := &models.PaginatedResult[models.State]{Data: data, Total: total}
	if !params.NoPage {
		result.Limit = params.Limit
		result.Offset = params.Offset
	}

	cache.Set(ctx, s.redis, key, result, cache.ListTTL)
	return result, nil
}

func (s *StateService) ListByCountryID(ctx context.Context, countryID int64, params models.QueryParams) (*models.PaginatedResult[models.State], error) {
	clampParams(&params)
	key := fmt.Sprintf("states:country:%d:%s:%s:%d:%d:%v", countryID, params.Search, params.Name, params.Limit, params.Offset, params.NoPage)
	cached, _ := cache.Get[models.PaginatedResult[models.State]](ctx, s.redis, key)
	if cached != nil {
		return cached, nil
	}

	data, total, err := s.repo.ListByCountryID(ctx, countryID, params)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []models.State{}
	}

	result := &models.PaginatedResult[models.State]{Data: data, Total: total}
	if !params.NoPage {
		result.Limit = params.Limit
		result.Offset = params.Offset
	}

	cache.Set(ctx, s.redis, key, result, cache.ListTTL)
	return result, nil
}
