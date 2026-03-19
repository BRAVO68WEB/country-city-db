package services

import "github.com/bravo68web/country-city-db/internal/models"

func clampParams(p *models.QueryParams) {
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if !p.NoPage && p.Limit > 100 {
		p.Limit = 100
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}
