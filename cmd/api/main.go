package main

import (
	"os"

	"github.com/commonsyllabi/viewer/internal/api"
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

	zero.Info("starting IMSCC viewer")

	var conf api.Config
	conf.DefaultConf()

	port := os.Getenv("PORT")
	if port == "" {
		zero.Log.Warn().Msg("Missing env PORT, defaulting to 8080")
		port = "8080"
	}

	err := api.StartServer(port, debug, conf)
	if err != nil {
		zero.Log.Fatal().Msgf("Error starting server: %v", err)
	}
}
