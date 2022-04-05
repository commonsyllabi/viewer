package db

import (
	"context"
	"os"

	zero "commonsyllabi/internal/logger"

	"github.com/jackc/pgx/v4"
)

const DATABASE_URL = "postgres://postgres@localhost:5432/testdatabase"

func Connect() {
	conn, err := pgx.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		zero.Log.Error().Msgf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	zero.Log.Info().Msgf("Connected to database: %v", DATABASE_URL)

	rows, err := conn.Query(context.Background(), "select * from testtable")
	if err != nil {
		zero.Log.Error().Msgf("Error getting row: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var num int
		var name string

		if err := rows.Scan(&num, &name); err != nil {
			zero.Log.Error().Msgf("Error looping over row: %v", err)
		}

		zero.Log.Info().Msgf("reading name %s and num %d", name, num)
	}

	if err := rows.Err(); err != nil {
		zero.Log.Error().Msgf("Error final over row: %v", err)
	}
}
