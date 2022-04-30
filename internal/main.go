package main

import (
	"os"

	"github.com/commonsyllabi/viewer/internal/api"
	"github.com/commonsyllabi/viewer/internal/db"
	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		zero.Log.Fatal().Msgf("Error loading .env file: %v", err)
	}

	zero.InitLog(0)
	zero.Log.Info().Msg("Starting CoSyl")

	db.Connect(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"))
	err = api.StartServer()
	if err != nil {
		zero.Log.Fatal().Msgf("Error starting server: %v", err)
	}
}
