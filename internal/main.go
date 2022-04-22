package main

import (
	"github.com/commonsyllabi/viewer/internal/api"
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

	// db.Connect(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.GetEnv("DB_HOST"))
	api.StartServer()
}
