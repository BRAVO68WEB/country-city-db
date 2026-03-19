package routes

import (
	"net/http"
	"strconv"

	apperrors "github.com/bravo68web/country-city-db/internal/errors"
	"github.com/bravo68web/country-city-db/internal/models"
	"github.com/bravo68web/country-city-db/internal/services"
	"github.com/gin-gonic/gin"
)

func registerRegionRoutes(rg *gin.RouterGroup, regionSvc *services.RegionService, subregionSvc *services.SubregionService, countrySvc *services.CountryService) {
	rg.GET("/regions", listRegions(regionSvc))
	rg.POST("/regions", listRegions(regionSvc))
	rg.GET("/regions/:id", getRegion(regionSvc))
	rg.GET("/regions/:id/subregions", listRegionSubregions(subregionSvc))
	rg.POST("/regions/:id/subregions", listRegionSubregions(subregionSvc))
	rg.GET("/regions/:id/countries", listRegionCountries(countrySvc))
	rg.POST("/regions/:id/countries", listRegionCountries(countrySvc))
}

func listRegions(svc *services.RegionService) gin.HandlerFunc {
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

func getRegion(svc *services.RegionService) gin.HandlerFunc {
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

func listRegionSubregions(svc *services.SubregionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		params := bindParams(c)
		result, err := svc.ListByRegionID(c.Request.Context(), id, params)
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func listRegionCountries(svc *services.CountryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		params := bindParams(c)
		result, err := svc.ListByRegionID(c.Request.Context(), id, params)
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func bindParams(c *gin.Context) models.QueryParams {
	var params models.QueryParams
	if c.Request.Method == http.MethodPost {
		c.ShouldBindJSON(&params)
	} else {
		c.ShouldBindQuery(&params)
	}
	return params
}

func handleError(c *gin.Context, err error) {
	status := apperrors.StatusCode(err)
	c.JSON(status, gin.H{"error": err.Error()})
}
