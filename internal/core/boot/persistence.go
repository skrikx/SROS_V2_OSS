package boot

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"srosv2/internal/shared/config"
)

type Persistence struct {
	DB            *sql.DB
	Enabled       bool
	Connected     bool
	Driver        string
	URLPresent    bool
	MigrationsDir string
	Schema        string
}

func initPersistence(cfg config.Config) (*Persistence, error) {
	p := &Persistence{
		Enabled:       cfg.Database.Enabled,
		Driver:        cfg.Database.Driver,
		URLPresent:    strings.TrimSpace(cfg.Database.URL) != "",
		MigrationsDir: cfg.Database.MigrationsDir,
		Schema:        cfg.Database.Schema,
	}
	if err := os.MkdirAll(filepath.Join(cfg.ArtifactRoot, "bundles"), 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(cfg.ArtifactRoot, "receipts"), 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(cfg.ArtifactRoot, "releases"), 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(cfg.ArtifactRoot, "replays"), 0o755); err != nil {
		return nil, err
	}
	if !cfg.Database.Enabled || strings.TrimSpace(cfg.Database.URL) == "" {
		return p, nil
	}
	db, err := sql.Open(cfg.Database.Driver, cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return p, nil
	}
	p.DB = db
	p.Connected = true
	return p, nil
}

func (p *Persistence) Close() error {
	if p == nil || p.DB == nil {
		return nil
	}
	return p.DB.Close()
}

func (p *Persistence) Summary() map[string]any {
	return map[string]any{
		"enabled":        p.Enabled,
		"connected":      p.Connected,
		"driver":         p.Driver,
		"url_present":    p.URLPresent,
		"migrations_dir": p.MigrationsDir,
		"schema":         p.Schema,
	}
}
