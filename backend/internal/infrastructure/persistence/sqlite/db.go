package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaFS embed.FS

type DB struct {
	SQL *sql.DB
}

type Config struct {
	// Path to SQLite file, e.g. "data/app.db". If empty, defaults to "data/app.db".
	Path string
}

func Open(ctx context.Context, cfg Config) (*DB, error) {
	path := strings.TrimSpace(cfg.Path)
	if path == "" {
		path = filepath.Join("data", "app.db")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("mkdir db dir: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	// Keep a single connection so PRAGMA settings are consistent.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	// Basic ping with timeout.
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys = ON;"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	if err := ApplySchema(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &DB{SQL: db}, nil
}

func (d *DB) Close() error {
	if d == nil || d.SQL == nil {
		return nil
	}
	return d.SQL.Close()
}

func ApplySchema(ctx context.Context, db *sql.DB) error {
	schemaBytes, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("read embedded schema: %w", err)
	}

	schema := string(schemaBytes)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	for _, stmt := range strings.Split(schema, ";") {
		s := strings.TrimSpace(stmt)
		if s == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return fmt.Errorf("apply schema stmt %q: %w", s, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit schema tx: %w", err)
	}
	return nil
}

