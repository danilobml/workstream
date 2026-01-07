package db

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type DBInterface interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}
