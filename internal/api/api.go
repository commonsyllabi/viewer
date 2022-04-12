package api

import (
	"fmt"
	"net/http"

	zero "github.com/commonsyllabi/viewer/internal/logger"
)

const PORT string = "2046"

func StartServer() {
	zero.Log.Info().Msgf("Starting API on port %s", PORT)
	handler := http.HandlerFunc(handleRequest)
	http.Handle("/ping", handler)
	http.ListenAndServe(":"+PORT, nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	zero.Log.Debug().Msg("Received ping")
	fmt.Fprintf(w, "pong\n")
}
