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

	_, err := models.InitDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		zero.Log.Fatal().Msgf("Error initializing D: %v", err)
	}
	err = api.StartServer()
	if err != nil {
		zero.Log.Fatal().Msgf("Error starting server: %v", err)
	}
}
