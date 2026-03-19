package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/bravo68web/country-city-db/internal/models"
)

func TestPing(t *testing.T) {
	w := makeRequest(http.MethodGet, "/ping", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]string
	parseJSON(w, &resp)
	if resp["message"] != "pong" {
		t.Fatalf("expected pong, got %s", resp["message"])
	}
}

func TestGetStats(t *testing.T) {
	w := makeRequest(http.MethodGet, "/api/v1/stats", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp models.StatsResponse
	parseJSON(w, &resp)
	if resp.Database.Regions == 0 {
		t.Fatal("expected non-zero region count")
	}
	if resp.Database.Countries == 0 {
		t.Fatal("expected non-zero country count")
	}
}

func TestListRegions(t *testing.T) {
	w := makeRequest(http.MethodGet, "/api/v1/regions", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp models.PaginatedResult[models.Region]
	parseJSON(w, &resp)
	if resp.Total == 0 {
		t.Fatal("expected non-zero total")
	}
	if len(resp.Data) == 0 {
		t.Fatal("expected data")
	}
}

func TestGetRegionByID(t *testing.T) {
	w := makeRequest(http.MethodGet, "/api/v1/regions/1", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp models.Region
	parseJSON(w, &resp)
	if resp.ID != 1 {
		t.Fatalf("expected id 1, got %d", resp.ID)
	}
}

func TestGetRegionNotFound(t *testing.T) {
	w := makeRequest(http.MethodGet, "/api/v1/regions/999999", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestListRegionSubregions(t *testing.T) {
	w := makeRequest(http.MethodGet, "/api/v1/regions/1/subregions", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp models.PaginatedResult[models.Subregion]
	parseJSON(w, &resp)
	// Should be a valid response shape
	if resp.Data == nil {
		t.Fatal("expected data array")
	}
}

func TestSearchCountriesGET(t *testing.T) {
	w := makeRequest(http.MethodGet, "/api/v1/countries?search=India", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp models.PaginatedResult[models.Country]
	parseJSON(w, &resp)
	if resp.Total == 0 {
		t.Fatal("expected results for India search")
	}
	found := false
	for _, c := range resp.Data {
		if c.Name == "India" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected India in results")
	}
}

func TestSearchCountriesPOST(t *testing.T) {
	body := map[string]string{"search": "India"}
	w := makeRequest(http.MethodPost, "/api/v1/countries", body)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp models.PaginatedResult[models.Country]
	parseJSON(w, &resp)
	if resp.Total == 0 {
		t.Fatal("expected results for India search via POST")
	}
}

func TestGetCountryByISO2(t *testing.T) {
	w := makeRequest(http.MethodGet, "/api/v1/countries/iso2/US", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp models.Country
	parseJSON(w, &resp)
	if resp.Name != "United States" {
		t.Fatalf("expected United States, got %s", resp.Name)
	}
}

func TestGetCountryByISO3(t *testing.T) {
	w := makeRequest(http.MethodGet, "/api/v1/countries/iso3/IND", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp models.Country
	parseJSON(w, &resp)
	if resp.Name != "India" {
		t.Fatalf("expected India, got %s", resp.Name)
	}
}

func TestGetCountryByName(t *testing.T) {
	w := makeRequest(http.MethodGet, "/api/v1/countries/name/India", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp models.Country
	parseJSON(w, &resp)
	if resp.Name != "India" {
		t.Fatalf("expected India, got %s", resp.Name)
	}
}

func TestPostCountriesWithISO2(t *testing.T) {
	body := map[string]string{"iso2": "US"}
	w := makeRequest(http.MethodPost, "/api/v1/countries", body)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp models.PaginatedResult[models.Country]
	parseJSON(w, &resp)
	if resp.Total == 0 {
		t.Fatal("expected results for US iso2 filter")
	}
	found := false
	for _, c := range resp.Data {
		if c.Name == "United States" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected United States in results")
	}
}

func TestListCountryStates(t *testing.T) {
	// Get India's ID first
	w := makeRequest(http.MethodGet, "/api/v1/countries/iso2/IN", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var country models.Country
	parseJSON(w, &country)

	w = makeRequest(http.MethodGet, "/api/v1/countries/"+itoa(country.ID)+"/states", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp models.PaginatedResult[models.State]
	parseJSON(w, &resp)
	if resp.Total == 0 {
		t.Fatal("expected states for India")
	}
}

func TestListStateCities(t *testing.T) {
	// Use state ID 1 (should exist)
	w := makeRequest(http.MethodGet, "/api/v1/states/1/cities", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp models.PaginatedResult[models.City]
	parseJSON(w, &resp)
	if resp.Data == nil {
		t.Fatal("expected data array")
	}
}

func TestPagination(t *testing.T) {
	w := makeRequest(http.MethodGet, "/api/v1/countries?limit=5&offset=2", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp models.PaginatedResult[models.Country]
	parseJSON(w, &resp)
	if len(resp.Data) > 5 {
		t.Fatalf("expected at most 5 items, got %d", len(resp.Data))
	}
	if resp.Limit != 5 {
		t.Fatalf("expected limit 5, got %d", resp.Limit)
	}
	if resp.Offset != 2 {
		t.Fatalf("expected offset 2, got %d", resp.Offset)
	}
}

func TestNoPageMode(t *testing.T) {
	w := makeRequest(http.MethodGet, "/api/v1/regions?no_page=true", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp models.PaginatedResult[models.Region]
	parseJSON(w, &resp)
	// In no_page mode, all results returned and limit/offset should be 0
	if resp.Limit != 0 {
		t.Fatalf("expected limit 0 in no_page mode, got %d", resp.Limit)
	}
	if resp.Offset != 0 {
		t.Fatalf("expected offset 0 in no_page mode, got %d", resp.Offset)
	}
	if int64(len(resp.Data)) != resp.Total {
		t.Fatalf("expected all %d results, got %d", resp.Total, len(resp.Data))
	}
}

func TestInvalidID(t *testing.T) {
	w := makeRequest(http.MethodGet, "/api/v1/regions/abc", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestNonExistentID(t *testing.T) {
	w := makeRequest(http.MethodGet, "/api/v1/countries/999999", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestSearchNoResults(t *testing.T) {
	w := makeRequest(http.MethodGet, "/api/v1/countries?search=ZZZZNOTACOUNTRY", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp models.PaginatedResult[models.Country]
	parseJSON(w, &resp)
	if len(resp.Data) != 0 {
		t.Fatalf("expected empty data array, got %d items", len(resp.Data))
	}
	if resp.Total != 0 {
		t.Fatalf("expected total 0, got %d", resp.Total)
	}
}

func TestUpdateEndpointDisabledWithoutKey(t *testing.T) {
	// The test router is built with empty internalKey, so the route should not be registered
	w := makeRequest(http.MethodPost, "/api/v1/update", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 when INTERNAL_KEY is empty, got %d", w.Code)
	}
}

func TestUpdateEndpointUnauthorized(t *testing.T) {
	// Build a router with an internal key to test auth
	router := buildRouterWithKey(testPool, testRedis, "test-secret")

	// No header
	req := makeRequestWithRouter(router, http.MethodPost, "/api/v1/update", nil, nil)
	if req.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without header, got %d", req.Code)
	}

	// Wrong header
	headers := map[string]string{"X-Internal-Key": "wrong-key"}
	req = makeRequestWithRouter(router, http.MethodPost, "/api/v1/update", nil, headers)
	if req.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 with wrong key, got %d", req.Code)
	}
}

func itoa(n int64) string {
	return fmt.Sprintf("%d", n)
}
