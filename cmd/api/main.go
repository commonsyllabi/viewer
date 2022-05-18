package main

import (
	"os"

	"github.com/commonsyllabi/viewer/internal/api"
	"github.com/commonsyllabi/viewer/internal/api/models"
	zero "github.com/commonsyllabi/viewer/internal/logger"
)

func main() {
	debug := false
	switch os.Getenv("DEBUG") {
	case "true":
		debug = true
		zero.InitLog(0)
	case "false":
		debug = false
		zero.InitLog(1)
	default:
		zero.Log.Warn().Msg("Missing env DEBUG, defaulting to false")
		zero.InitLog(1)
	}

	zero.Info("starting cosyl")

	var conf api.Config
	conf.DefaultConf()

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		zero.Log.Fatal().Msg("Missing env DATABASE_URL")
	}

	port := os.Getenv("PORT")
	if port == "" {
		zero.Log.Warn().Msg("Missing env PORT, defaulting to 8080")
		port = "8080"
	}

	_, err := models.InitDB(url, conf.FixturesDir)
	if err != nil {
		zero.Log.Fatal().Msgf("Error initializing D: %v", err)
	}

	err = api.StartServer(port, debug, conf)
	if err != nil {
		zero.Log.Fatal().Msgf("Error starting server: %v", err)
	}
}
