package routes

import (
	"github.com/bravo68web/country-city-db/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Register(r *gin.Engine,
	regionSvc *services.RegionService,
	subregionSvc *services.SubregionService,
	countrySvc *services.CountryService,
	stateSvc *services.StateService,
	citySvc *services.CityService,
	statsSvc *services.StatsService,
	pool *pgxpool.Pool,
	redisClient *redis.Client,
	internalKey string,
) {
	v1 := r.Group("/api/v1")

	registerRegionRoutes(v1, regionSvc, subregionSvc, countrySvc)
	registerSubregionRoutes(v1, subregionSvc)
	registerCountryRoutes(v1, countrySvc, stateSvc, citySvc)
	registerStateRoutes(v1, stateSvc, citySvc)
	registerCityRoutes(v1, citySvc)
	registerStatsRoutes(v1, statsSvc)

	if internalKey != "" {
		registerUpdateRoute(v1, pool, redisClient, internalKey)
	}
}
