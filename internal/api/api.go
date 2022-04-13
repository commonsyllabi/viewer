package api

import (
	"fmt"
	"net/http"

	zero "github.com/commonsyllabi/viewer/internal/logger"
)

func StartServer(port string) {
	zero.Log.Info().Msgf("Starting API on port %s", port)
	handler := http.HandlerFunc(handlePing)
	http.Handle("/ping", handler)
	http.ListenAndServe(":"+port, nil)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	zero.Log.Debug().Msg("Received ping")
	fmt.Fprintf(w, "pong")
}
