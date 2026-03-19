package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/bravo68web/country-city-db/internal/cache"
	"github.com/bravo68web/country-city-db/internal/db"
	"github.com/bravo68web/country-city-db/internal/migration"
	"github.com/bravo68web/country-city-db/internal/models"
	"github.com/bravo68web/country-city-db/internal/repositories"
	"github.com/bravo68web/country-city-db/internal/routes"
	"github.com/bravo68web/country-city-db/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var (
	testRouter *gin.Engine
	testPool   *pgxpool.Pool
	testRedis  *redis.Client
)

func init() {
	gin.SetMode(gin.TestMode)

	mode := os.Getenv("TEST_MODE")
	if mode == "real" {
		setupRealServer()
	} else {
		setupMockServer()
	}
}

func ensureMigrated(ctx context.Context, pool *pgxpool.Pool) {
	migrated, err := migration.IsMigrated(ctx, pool)
	if err != nil {
		panic("failed to check migration status: " + err.Error())
	}
	if !migrated {
		if err := migration.RunMigration(ctx, pool); err != nil {
			panic("migration failed: " + err.Error())
		}
	}
}

func setupRealServer() {
	ctx := context.Background()

	if url := os.Getenv("TEST_DATABASE_URL"); url != "" {
		os.Setenv("DATABASE_URL", url)
	}
	if url := os.Getenv("TEST_REDIS_URL"); url != "" {
		os.Setenv("REDIS_URL", url)
	}

	pool, err := db.NewPool(ctx)
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}
	testPool = pool

	ensureMigrated(ctx, pool)

	redisClient, err := cache.NewRedisClient(ctx)
	if err != nil {
		panic("failed to connect to test redis: " + err.Error())
	}
	testRedis = redisClient

	testRouter = buildRouter(pool, redisClient)
}

func setupMockServer() {
	ctx := context.Background()

	// Use a minimal Redis mock via miniredis or just use a real Redis if available
	// For mock mode, we try to connect to Redis; if unavailable, tests that need it will skip
	addr := os.Getenv("REDIS_URL")
	if addr == "" {
		addr = "localhost:6374"
	}

	redisClient := redis.NewClient(&redis.Options{Addr: addr})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		// If Redis not available, create a client that will just error on cache ops
		// The service layer handles cache errors gracefully
		redisClient = redis.NewClient(&redis.Options{Addr: "localhost:6374"})
	}
	testRedis = redisClient

	// For mock mode, try to connect to real DB too — the tests are E2E
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://postgres:postgres@localhost:5456/postgres"
		os.Setenv("DATABASE_URL", dbURL)
	}

	pool, err := db.NewPool(ctx)
	if err != nil {
		panic("failed to connect to database for mock mode (ensure docker-compose is up): " + err.Error())
	}
	testPool = pool

	ensureMigrated(ctx, pool)

	testRouter = buildRouter(pool, redisClient)
}

func buildRouter(pool *pgxpool.Pool, redisClient *redis.Client) *gin.Engine {
	regionRepo := repositories.NewRegionRepository(pool)
	subregionRepo := repositories.NewSubregionRepository(pool)
	countryRepo := repositories.NewCountryRepository(pool)
	stateRepo := repositories.NewStateRepository(pool)
	cityRepo := repositories.NewCityRepository(pool)

	regionSvc := services.NewRegionService(regionRepo, redisClient)
	subregionSvc := services.NewSubregionService(subregionRepo, redisClient)
	countrySvc := services.NewCountryService(countryRepo, redisClient)
	stateSvc := services.NewStateService(stateRepo, redisClient)
	citySvc := services.NewCityService(cityRepo, redisClient)
	statsSvc := services.NewStatsService(pool, redisClient)

	r := gin.New()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	routes.Register(r, regionSvc, subregionSvc, countrySvc, stateSvc, citySvc, statsSvc, pool, redisClient, "")
	return r
}

func makeRequest(method, path string, body any) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func parseJSON(w *httptest.ResponseRecorder, v any) error {
	return json.Unmarshal(w.Body.Bytes(), v)
}

func buildRouterWithKey(pool *pgxpool.Pool, redisClient *redis.Client, internalKey string) *gin.Engine {
	regionRepo := repositories.NewRegionRepository(pool)
	subregionRepo := repositories.NewSubregionRepository(pool)
	countryRepo := repositories.NewCountryRepository(pool)
	stateRepo := repositories.NewStateRepository(pool)
	cityRepo := repositories.NewCityRepository(pool)

	regionSvc := services.NewRegionService(regionRepo, redisClient)
	subregionSvc := services.NewSubregionService(subregionRepo, redisClient)
	countrySvc := services.NewCountryService(countryRepo, redisClient)
	stateSvc := services.NewStateService(stateRepo, redisClient)
	citySvc := services.NewCityService(cityRepo, redisClient)
	statsSvc := services.NewStatsService(pool, redisClient)

	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	routes.Register(r, regionSvc, subregionSvc, countrySvc, stateSvc, citySvc, statsSvc, pool, redisClient, internalKey)
	return r
}

func makeRequestWithRouter(router *gin.Engine, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// Suppress unused import warnings
var (
	_ = time.Now
	_ models.Region
)
