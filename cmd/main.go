package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	spec "github.com/bravo68web/country-city-db"
	"github.com/bravo68web/country-city-db/internal/cache"
	"github.com/bravo68web/country-city-db/internal/db"
	"github.com/bravo68web/country-city-db/internal/migration"
	"github.com/bravo68web/country-city-db/internal/repositories"
	"github.com/bravo68web/country-city-db/internal/routes"
	"github.com/bravo68web/country-city-db/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	runtime.GOMAXPROCS(4)

	ctx := context.Background()

	pool, err := db.NewPool(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Auto-migration: detect empty DB and import data
	migrated, err := migration.IsMigrated(ctx, pool)
	if err != nil {
		log.Fatalf("Failed to check migration status: %v", err)
	}
	if !migrated {
		log.Println("Database not migrated. Downloading and importing...")
		if err := migration.RunMigration(ctx, pool); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migration completed.")
	}

	redisClient, err := cache.NewRedisClient(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Repositories
	regionRepo := repositories.NewRegionRepository(pool)
	subregionRepo := repositories.NewSubregionRepository(pool)
	countryRepo := repositories.NewCountryRepository(pool)
	stateRepo := repositories.NewStateRepository(pool)
	cityRepo := repositories.NewCityRepository(pool)

	// Services
	regionSvc := services.NewRegionService(regionRepo, redisClient)
	subregionSvc := services.NewSubregionService(subregionRepo, redisClient)
	countrySvc := services.NewCountryService(countryRepo, redisClient)
	stateSvc := services.NewStateService(stateRepo, redisClient)
	citySvc := services.NewCityService(cityRepo, redisClient)
	statsSvc := services.NewStatsService(pool, redisClient)

	internalKey := os.Getenv("INTERNAL_KEY")

	// Router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Internal-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/openapi", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/yaml", spec.OpenAPISpec)
	})

	r.GET("/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`<!doctype html>
<html>
<head><title>API Docs</title><meta charset="utf-8" /></head>
<body>
  <script id="api-reference" data-url="/openapi"></script>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`))
	})

	routes.Register(r, regionSvc, subregionSvc, countrySvc, stateSvc, citySvc, statsSvc, pool, redisClient, internalKey)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exited")
}
