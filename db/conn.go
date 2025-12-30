package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool

func IntialiseDB() error {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		return fmt.Errorf("initialising DB Connection: %s", err.Error())
	}
	DBPool = pool
	return nil
}
