package postgres

import (
	"context"
	"fmt"
	"go-template/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// New creates a pgxpool.Pool. Call Ping separately with a deadline-aware context.
func New(cfg config.Postgres, log *zap.Logger) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse pg config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pg pool: %w", err)
	}

	log.Info("postgres pool created", zap.String("host", cfg.Host), zap.Int("port", cfg.Port))
	return pool, nil
}
