package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func InitDB(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
    
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	fmt.Println("Connected to PostgreSQL database!")
	return pool, nil
}
