# Country City DB

A fast, self-migrating REST API serving worldwide geographic data — countries, states, cities, regions, and subregions — backed by PostgreSQL and Redis.

## Features

- **250+ countries**, 5,000+ states, 150,000+ cities with metadata (ISO codes, coordinates, population, timezones, translations)
- **Auto-migration** — on first boot, the database is populated automatically
- **Redis caching** for sub-millisecond repeated lookups
- **Pagination** with configurable `limit`, `offset`, and `no_page` mode
- **Full-text search** across all entity types
- **Lookup by ISO2, ISO3, or name** for countries
- **Interactive API docs** via Scalar at `/docs`
- **CORS enabled** — open to all origins

## Credits

Data Collected from [https://github.com/dr5hn/countries-states-cities-database](https://github.com/dr5hn/countries-states-cities-database/)

## Quick Start

### Docker Compose (recommended)

```bash
docker compose up
```

The API will be available at `http://localhost:8080`. On first startup, it downloads and imports the world database automatically.

### Local Development

```bash
# Start dependencies
docker compose -f docker-compose.dev.yaml up -d

# Copy and edit env
cp .env.local .env

# Run the server
go run ./cmd/main.go
```

## API Reference

Interactive docs: **[http://localhost:8080/docs](http://localhost:8080/docs)**

Raw OpenAPI spec: **[http://localhost:8080/openapi](http://localhost:8080/openapi)**

### Health Check

```
GET /ping
→ { "message": "pong" }
```

### Regions

```
GET    /api/v1/regions                  # List all regions
POST   /api/v1/regions                  # List with JSON body filters
GET    /api/v1/regions/:id              # Get region by ID
GET    /api/v1/regions/:id/subregions   # Subregions in a region
GET    /api/v1/regions/:id/countries    # Countries in a region
```

### Subregions

```
GET    /api/v1/subregions               # List all subregions
POST   /api/v1/subregions               # List with JSON body filters
GET    /api/v1/subregions/:id           # Get subregion by ID
```

### Countries

```
GET    /api/v1/countries                # List all countries
POST   /api/v1/countries                # List with JSON body filters
GET    /api/v1/countries/:id            # Get country by ID
GET    /api/v1/countries/iso2/:code     # Lookup by ISO2 (e.g. US, IN)
GET    /api/v1/countries/iso3/:code     # Lookup by ISO3 (e.g. USA, IND)
GET    /api/v1/countries/name/:name     # Lookup by name
GET    /api/v1/countries/:id/states     # States in a country
GET    /api/v1/countries/:id/cities     # Cities in a country
```

### States

```
GET    /api/v1/states                   # List all states
POST   /api/v1/states                   # List with JSON body filters
GET    /api/v1/states/:id               # Get state by ID
GET    /api/v1/states/:id/cities        # Cities in a state
```

### Cities

```
GET    /api/v1/cities                   # List all cities
POST   /api/v1/cities                   # List with JSON body filters
GET    /api/v1/cities/:id               # Get city by ID
```

### Stats

```
GET    /api/v1/stats                    # Database & cache statistics
```

## Query Parameters

All list endpoints accept these parameters (via query string on GET, JSON body on POST):

| Parameter | Type    | Default | Description                         |
|-----------|---------|---------|-------------------------------------|
| `search`  | string  |         | Full-text search                    |
| `name`    | string  |         | Filter by exact name                |
| `iso2`    | string  |         | Filter by ISO2 code (countries)     |
| `iso3`    | string  |         | Filter by ISO3 code (countries)     |
| `limit`   | integer | 20      | Max results per page (max 100)      |
| `offset`  | integer | 0       | Number of results to skip           |
| `no_page` | boolean | false   | Return all results without paging   |

## Examples

```bash
# Search for countries matching "united"
curl "http://localhost:8080/api/v1/countries?search=united"

# Get India by ISO2
curl "http://localhost:8080/api/v1/countries/iso2/IN"

# List states in the US with pagination
curl "http://localhost:8080/api/v1/countries/233/states?limit=10&offset=0"

# Search cities via POST
curl -X POST http://localhost:8080/api/v1/cities \
  -H "Content-Type: application/json" \
  -d '{"search": "New York", "limit": 5}'

# Get all regions (no pagination)
curl "http://localhost:8080/api/v1/regions?no_page=true"
```

## Response Format

All list endpoints return a paginated wrapper:

```json
{
  "data": [...],
  "total": 250,
  "limit": 20,
  "offset": 0
}
```

Single-resource endpoints return the object directly.

## Docker Image

Pre-built images are published to GHCR on every tagged release:

```bash
docker pull ghcr.io/bravo68web/country-city-db:latest
```

## Running Tests

```bash
# Start test dependencies
docker compose -f docker-compose.dev.yaml up -d

# Run all tests
DATABASE_URL=postgresql://postgres:postgres@localhost:5456/postgres \
REDIS_URL=localhost:6374 \
go test -v ./tests/...
```

## Environment Variables

| Variable       | Default                                            | Description              |
|----------------|----------------------------------------------------|--------------------------|
| `DATABASE_URL` | `postgresql://postgres:postgres@localhost:5432/postgres` | PostgreSQL connection string |
| `REDIS_URL`    | `localhost:6379`                                   | Redis address            |
| `PORT`         | `8080`                                             | Server port              |
| `INTERNAL_KEY` |                                                    | Key for `/api/v1/update` |

## License

[./LICENSE](LICENSE)
