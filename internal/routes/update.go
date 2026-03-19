package routes

import (
	"net/http"

	"github.com/bravo68web/country-city-db/internal/migration"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func registerUpdateRoute(v1 *gin.RouterGroup, pool *pgxpool.Pool, redisClient *redis.Client, internalKey string) {
	v1.POST("/update", func(c *gin.Context) {
		if c.GetHeader("X-Internal-Key") != internalKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if err := migration.RunUpdate(c.Request.Context(), pool, redisClient); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "update completed"})
	})
}
