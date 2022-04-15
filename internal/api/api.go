package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	zero "github.com/commonsyllabi/viewer/internal/logger"
)

func StartServer(port string) {
	zero.Log.Info().Msgf("Starting API on port %s", port)

	http.Handle("/ping", http.HandlerFunc(handlePing))
	http.Handle("/upload", http.HandlerFunc(handleUpload))

	http.ListenAndServe(":"+port, nil)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	zero.Log.Debug().Msg("Received ping")
	fmt.Fprintf(w, "pong")
}

func handleUpload(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		zero.Log.Warn().Msgf("Method not allowed: %s", r.Method)
		return
	}

	file, fileHeader, err := r.FormFile("cartridge")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		zero.Log.Warn().Msgf("Cannot find cartridge file: %s", r.FormValue("cartridge"))
		return
	}
	defer file.Close()

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(fmt.Sprintf("./uploads/%s", fileHeader.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["msg"] = "success"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		zero.Log.Error().Msg("cannot marshal response to json")
	}

	w.Write(jsonResp)
}
