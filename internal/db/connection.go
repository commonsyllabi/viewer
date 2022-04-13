package db

import (
	"context"
	"os"

	zero "github.com/commonsyllabi/viewer/internal/logger"

	"github.com/jackc/pgx/v4"
)

func Connect(user, password, name string) {
	var db_url = "postgres://" + user + ":" + password + "@db:5432/" + name
	conn, err := pgx.Connect(context.Background(), db_url)
	if err != nil {
		zero.Log.Error().Msgf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	zero.Log.Info().Msgf("Connected to database: %v", db_url)
}
