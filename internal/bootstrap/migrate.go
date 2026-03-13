package bootstrap

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gym-pro-2026-ptit/internal/infrastructure/database"
	"gym-pro-2026-ptit/internal/infrastructure/logger"

	"github.com/jackc/pgx/v5"
)

const createSchemaMigrationsTable = `CREATE TABLE IF NOT EXISTS schema_migrations (
	version BIGINT PRIMARY KEY,
	name TEXT NOT NULL,
	applied_at TIMESTAMPTZ DEFAULT NOW()
);`

type migrationFile struct {
	Version uint
	Name    string
}

func listMigrationFiles(dir string) ([]migrationFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	upRe := regexp.MustCompile(`^(\d+)_(.+)\.up\.sql$`)
	var list []migrationFile
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		sub := upRe.FindStringSubmatch(name)
		if len(sub) < 3 {
			continue
		}
		ver, err := strconv.ParseUint(sub[1], 10, 64)
		if err != nil {
			continue
		}
		list = append(list, migrationFile{Version: uint(ver), Name: name})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Version < list[j].Version })
	return list, nil
}

func runSQLScript(ctx context.Context, db *database.DB, content string) error {
	parts := regexp.MustCompile(`;\s*\n`).Split(content, -1)
	for _, p := range parts {
		s := strings.TrimSpace(p)
		if s == "" {
			continue
		}
		if _, err := db.Exec(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func RunAutoMigrate(ctx context.Context, db *database.DB, migrationsPath string) error {
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("migrations path: %w", err)
	}
	if _, err := os.Stat(absPath); err != nil {
		if os.IsNotExist(err) {
			logger.Info("auto migrate skipped", "reason", "migrations directory not found", "path", absPath)
			return nil
		}
		return fmt.Errorf("migrations path: %w", err)
	}

	if _, err := db.Exec(ctx, createSchemaMigrationsTable); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	files, err := listMigrationFiles(absPath)
	if err != nil {
		return fmt.Errorf("list migrations: %w", err)
	}
	if len(files) == 0 {
		logger.Info("auto migrate", "migrations_dir", absPath, "status", "no migration files")
		return nil
	}

	logger.Info("auto migrate", "migrations_dir", absPath, "migration_files", len(files))

	for _, f := range files {
		var exists int
		err := db.QueryRow(ctx, "SELECT 1 FROM schema_migrations WHERE version = $1", f.Version).Scan(&exists)
		if err == nil {
			continue
		}
		if err != pgx.ErrNoRows {
			return fmt.Errorf("check migration version %d: %w", f.Version, err)
		}

		path := filepath.Join(absPath, f.Name)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", f.Name, err)
		}

		if err := runSQLScript(ctx, db, string(content)); err != nil {
			return fmt.Errorf("run migration %s: %w", f.Name, err)
		}

		if _, err := db.Exec(ctx, "INSERT INTO schema_migrations (version, name) VALUES ($1, $2)", f.Version, f.Name); err != nil {
			return fmt.Errorf("record migration %s: %w", f.Name, err)
		}

		logger.Info("auto migrate applied", "file", f.Name)
	}

	logger.Info("auto migrate", "status", "done")
	return nil
}
