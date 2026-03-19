package services

import (
	"context"
	"strings"

	"github.com/bravo68web/country-city-db/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type StatsService struct {
	pool  *pgxpool.Pool
	redis *redis.Client
}

func NewStatsService(pool *pgxpool.Pool, redis *redis.Client) *StatsService {
	return &StatsService{pool: pool, redis: redis}
}

func (s *StatsService) GetStats(ctx context.Context) (*models.StatsResponse, error) {
	var dbStats models.DatabaseStats

	tables := []struct {
		name  string
		dest  *int64
	}{
		{"regions", &dbStats.Regions},
		{"subregions", &dbStats.Subregions},
		{"countries", &dbStats.Countries},
		{"states", &dbStats.States},
		{"cities", &dbStats.Cities},
	}

	for _, t := range tables {
		if err := s.pool.QueryRow(ctx, "SELECT count(*) FROM "+t.name).Scan(t.dest); err != nil {
			return nil, err
		}
	}

	cacheStats := models.CacheStats{}
	info, err := s.redis.Info(ctx, "memory").Result()
	if err == nil {
		cacheStats.Connected = true
		for _, line := range strings.Split(info, "\r\n") {
			if strings.HasPrefix(line, "used_memory_human:") {
				cacheStats.Memory = strings.TrimPrefix(line, "used_memory_human:")
			}
		}
	}

	dbSize, err := s.redis.DBSize(ctx).Result()
	if err == nil {
		cacheStats.Keys = dbSize
	}

	return &models.StatsResponse{
		Database: dbStats,
		Cache:    cacheStats,
	}, nil
}
