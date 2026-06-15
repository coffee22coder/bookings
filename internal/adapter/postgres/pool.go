package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/coffee22coder/bookings/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	pool *pgxpool.Pool
}

func NewPool(ctx context.Context, cfg config.Config) (*Pool, error) {
	pool, err := pgxpool.New(ctx, config.DSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &Pool{pool: pool}, nil
}

func (p *Pool) Close() {
	p.pool.Close()
}

func (p *Pool) Ping(ctx context.Context) error {
	if err := p.pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}
	return nil
}

// func (p *Pool) FligthsCount(ctx context.Context) (int, error) {
// 	var count int
// 	err := p.pool.QueryRow(ctx, "SELECT COUNT(*) FROM bookings.flights").Scan(&count)
// 	return count, err
// }
