package routes

import (
	"net/http"
	"strconv"

	"github.com/bravo68web/country-city-db/internal/services"
	"github.com/gin-gonic/gin"
)

func registerStateRoutes(rg *gin.RouterGroup, stateSvc *services.StateService, citySvc *services.CityService) {
	rg.GET("/states", listStates(stateSvc))
	rg.POST("/states", listStates(stateSvc))
	rg.GET("/states/:id", getState(stateSvc))
	rg.GET("/states/:id/cities", listStateCities(citySvc))
	rg.POST("/states/:id/cities", listStateCities(citySvc))
}

func listStates(svc *services.StateService) gin.HandlerFunc {
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

func getState(svc *services.StateService) gin.HandlerFunc {
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

func listStateCities(svc *services.CityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		params := bindParams(c)
		result, err := svc.ListByStateID(c.Request.Context(), id, params)
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
