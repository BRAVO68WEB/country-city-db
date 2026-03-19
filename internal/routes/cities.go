package routes

import (
	"net/http"
	"strconv"

	"github.com/bravo68web/country-city-db/internal/services"
	"github.com/gin-gonic/gin"
)

func registerCityRoutes(rg *gin.RouterGroup, svc *services.CityService) {
	rg.GET("/cities", listCities(svc))
	rg.POST("/cities", listCities(svc))
	rg.GET("/cities/:id", getCity(svc))
}

func listCities(svc *services.CityService) gin.HandlerFunc {
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

func getCity(svc *services.CityService) gin.HandlerFunc {
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
