package database

import (
	"context"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
)

// NewDB establishes a database connection using the DATABASE_URL environment variable and
// returns it as *sqlx.DB. If the DATABASE_URL environment variable is not set, it
// logs a fatal error and exits the program. It also registers a hook with the given
// fx.Lifecycle to close the database connection when the lifecycle is stopped.
func NewDB(lc fx.Lifecycle) (*sqlx.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("FATAL: DATABASE_URL environment variable is not set")
	}

	db, err := sqlx.Connect("pgx", dbURL)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Closing database connection.")
			return db.Close()
		},
	})

	log.Println("Database connection established successfully.")
	return db, nil
}
