package db

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v4"
)

func InitDB(dsn string) (*pgx.Conn, error) {

    db, err := pgx.Connect(context.Background(), dsn)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(context.Background()); err != nil {
        db.Close(context.Background())
        return nil, err
    }
	
    fmt.Println("Connected to PostgreSQL database!")
    return db, nil
}