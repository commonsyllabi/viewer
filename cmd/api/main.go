package main

import (
	"os"

	"github.com/commonsyllabi/viewer/internal/api"
	"github.com/commonsyllabi/viewer/internal/api/models"
	zero "github.com/commonsyllabi/viewer/internal/logger"
)

func main() {
	zero.InitLog(0)
	zero.Info("Starting CoSyl")

	// config should be loaded in main

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		zero.Log.Fatal().Msg("Missing env DATABASE_URL")
	}

	_, err := models.InitDB(url)
	if err != nil {
		zero.Log.Fatal().Msgf("Error initializing D: %v", err)
	}
	err = api.StartServer()
	if err != nil {
		zero.Log.Fatal().Msgf("Error starting server: %v", err)
	}
}
