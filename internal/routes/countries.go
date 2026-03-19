package routes

import (
	"net/http"
	"strconv"

	"github.com/bravo68web/country-city-db/internal/services"
	"github.com/gin-gonic/gin"
)

func registerCountryRoutes(rg *gin.RouterGroup, countrySvc *services.CountryService, stateSvc *services.StateService, citySvc *services.CityService) {
	rg.GET("/countries", listCountries(countrySvc))
	rg.POST("/countries", listCountries(countrySvc))
	rg.GET("/countries/:id", getCountry(countrySvc))
	rg.GET("/countries/iso2/:code", getCountryByISO2(countrySvc))
	rg.GET("/countries/iso3/:code", getCountryByISO3(countrySvc))
	rg.GET("/countries/name/:name", getCountryByName(countrySvc))
	rg.GET("/countries/:id/states", listCountryStates(stateSvc))
	rg.POST("/countries/:id/states", listCountryStates(stateSvc))
	rg.GET("/countries/:id/cities", listCountryCities(citySvc))
	rg.POST("/countries/:id/cities", listCountryCities(citySvc))
}

func listCountries(svc *services.CountryService) gin.HandlerFunc {
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

func getCountry(svc *services.CountryService) gin.HandlerFunc {
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

func getCountryByISO2(svc *services.CountryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := svc.GetByISO2(c.Request.Context(), c.Param("code"))
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func getCountryByISO3(svc *services.CountryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := svc.GetByISO3(c.Request.Context(), c.Param("code"))
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func getCountryByName(svc *services.CountryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := svc.GetByName(c.Request.Context(), c.Param("name"))
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func listCountryStates(svc *services.StateService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		params := bindParams(c)
		result, err := svc.ListByCountryID(c.Request.Context(), id, params)
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func listCountryCities(svc *services.CityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		params := bindParams(c)
		result, err := svc.ListByCountryID(c.Request.Context(), id, params)
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
