package migration

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const sourceURL = "https://github.com/dr5hn/countries-states-cities-database/raw/refs/heads/master/psql/world.sql.gz"

var mu sync.Mutex

// IsMigrated checks if the database has been populated with data.
func IsMigrated(ctx context.Context, pool *pgxpool.Pool) (bool, error) {
	var exists bool
	err := pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = 'regions'
		)
	`).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking migration status: %w", err)
	}
	if !exists {
		return false, nil
	}

	var count int64
	err = pool.QueryRow(ctx, `SELECT count(*) FROM regions`).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("counting regions: %w", err)
	}
	return count > 0, nil
}

func downloadAndDecompress(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("downloading world.sql.gz: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d downloading world.sql.gz", resp.StatusCode)
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("creating gzip reader: %w", err)
	}
	defer gz.Close()

	data, err := io.ReadAll(gz)
	if err != nil {
		return nil, fmt.Errorf("decompressing world.sql.gz: %w", err)
	}

	return data, nil
}

func sanitizeSQL(raw []byte) []string {
	scanner := bufio.NewScanner(bytes.NewReader(raw))
	// Increase scanner buffer for long lines in the SQL dump
	scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024)

	var stmts []string
	var buf strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, `\`) {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}
		buf.WriteString(line)
		buf.WriteByte('\n')
		if strings.HasSuffix(trimmed, ";") {
			stmts = append(stmts, buf.String())
			buf.Reset()
		}
	}
	// Flush any remaining content without a trailing semicolon
	if buf.Len() > 0 {
		stmts = append(stmts, buf.String())
	}
	return stmts
}

func executeSQL(ctx context.Context, pool *pgxpool.Pool, stmts []string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	total := len(stmts)
	for i, stmt := range stmts {
		if (i+1)%100 == 0 || i+1 == total {
			log.Printf("Executing statement %d/%d...", i+1, total)
		}
		if _, err := tx.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("executing statement %d/%d: %w", i+1, total, err)
		}
	}

	// Reset search_path — the SQL dump sets it to '' at the session level,
	// which would poison this connection when returned to the pool.
	if _, err := tx.Exec(ctx, "SET search_path TO public;"); err != nil {
		return fmt.Errorf("resetting search_path: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

// RunMigration downloads and imports the world database.
func RunMigration(ctx context.Context, pool *pgxpool.Pool) error {
	mu.Lock()
	defer mu.Unlock()

	log.Println("Downloading world.sql.gz...")
	raw, err := downloadAndDecompress(ctx)
	if err != nil {
		return err
	}
	log.Printf("Downloaded and decompressed (%d bytes). Sanitizing...", len(raw))

	stmts := sanitizeSQL(raw)
	log.Printf("Sanitized into %d statements. Executing SQL import...", len(stmts))

	if err := executeSQL(ctx, pool, stmts); err != nil {
		return err
	}

	log.Println("SQL import complete.")
	return nil
}

// RunUpdate re-imports the world database and flushes the Redis cache.
func RunUpdate(ctx context.Context, pool *pgxpool.Pool, redisClient *redis.Client) error {
	if err := RunMigration(ctx, pool); err != nil {
		return err
	}

	log.Println("Flushing Redis cache...")
	if err := redisClient.FlushDB(ctx).Err(); err != nil {
		return fmt.Errorf("flushing redis: %w", err)
	}
	log.Println("Redis cache flushed.")
	return nil
}
