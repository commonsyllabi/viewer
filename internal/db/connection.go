package db

import (
	"context"
	"os"

	zero "github.com/commonsyllabi/viewer/internal/logger"

	"github.com/jackc/pgx/v4"
)

const DATABASE_URL = "postgres://postgres:postgres@db:5432/postgres"

func Connect() {
	conn, err := pgx.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		zero.Log.Error().Msgf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	zero.Log.Info().Msgf("Connected to database: %v", DATABASE_URL)
}
