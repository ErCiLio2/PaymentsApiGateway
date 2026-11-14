package db

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS _schema_migrations (
			version    TEXT        PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("criar tabela de migrações: %w", err)
	}

	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return err
	}

	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".up.sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		version := strings.TrimSuffix(name, ".up.sql")

		var exists bool
		if err := pool.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM _schema_migrations WHERE version = $1)`, version,
		).Scan(&exists); err != nil {
			return err
		}
		if exists {
			continue
		}

		sql, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			return err
		}

		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("migração %s: %w", version, err)
		}
		if _, err := pool.Exec(ctx,
			`INSERT INTO _schema_migrations (version) VALUES ($1)`, version,
		); err != nil {
			return fmt.Errorf("registar migração %s: %w", version, err)
		}
		log.Printf("[kashe] migração aplicada: %s", version)
	}
	return nil
}
