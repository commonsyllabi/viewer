package main

import (
	"github.com/commonsyllabi/viewer/internal/api"
	zero "github.com/commonsyllabi/viewer/internal/logger"
)

func main() {
	zero.InitLog(true)
	zero.Log.Info().Msg("Starting Cosyl")
	// db.Connect()
	api.StartServer()
}
