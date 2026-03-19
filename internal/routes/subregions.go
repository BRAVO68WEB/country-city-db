package routes

import (
	"net/http"
	"strconv"

	"github.com/bravo68web/country-city-db/internal/services"
	"github.com/gin-gonic/gin"
)

func registerSubregionRoutes(rg *gin.RouterGroup, svc *services.SubregionService) {
	rg.GET("/subregions", listSubregions(svc))
	rg.POST("/subregions", listSubregions(svc))
	rg.GET("/subregions/:id", getSubregion(svc))
}

func listSubregions(svc *services.SubregionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := bindParams(c)
		result, err := svc.List(c.Request.Context(), params)
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func getSubregion(svc *services.SubregionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		result, err := svc.GetByID(c.Request.Context(), id)
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
