package routes

import (
	"net/http"

	"github.com/bravo68web/country-city-db/internal/services"
	"github.com/gin-gonic/gin"
)

func registerStatsRoutes(rg *gin.RouterGroup, svc *services.StatsService) {
	rg.GET("/stats", getStats(svc))
}

func getStats(svc *services.StatsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := svc.GetStats(c.Request.Context())
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
