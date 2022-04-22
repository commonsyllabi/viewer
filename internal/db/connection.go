package db

import (
	"context"

	zero "github.com/commonsyllabi/viewer/internal/logger"

	"github.com/jackc/pgx/v4"
)

func Connect(user, password, name, host string) error {
	var db_url = "postgres://" + user + ":" + password + "@" + host + ":5432/" + name
	conn, err := pgx.Connect(context.Background(), db_url)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	zero.Log.Info().Msgf("Connected to database: %v", db_url)

	return nil
}
