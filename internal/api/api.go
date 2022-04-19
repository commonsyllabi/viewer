package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/commonsyllabi/viewer/pkg/commoncartridge"
)

// todo, make this dependent on env (inside docker or not)
const uploadsDir = "uploads"
const tmpDir = "tmp"

func StartServer(port string) {
	zero.Log.Info().Msgf("Starting API on port %s", port)

	http.Handle("/ping", http.HandlerFunc(handlePing))
	http.Handle("/upload", http.HandlerFunc(handleUpload))
	http.Handle("/resource/", http.HandlerFunc(handleResource))
	http.Handle("/file/", http.HandlerFunc(handleFile))
	http.Handle("/tmp/", http.FileServer(http.Dir(tmpDir)))

	http.ListenAndServe(":"+port, nil)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	zero.Log.Debug().Msg("Received ping")
	fmt.Fprintf(w, "pong")
}

func handleFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		zero.Log.Warn().Msgf("Method not allowed: %s", r.Method)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/file/")
	cartridge := r.FormValue("cartridge")
	zero.Log.Info().Msgf("GET handleFile id: %v cartridge %v", id, cartridge)

	inputFile := filepath.Join(uploadsDir, cartridge)
	cc, err := commoncartridge.Load(inputFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zero.Log.Error().Msgf("error loading CC from disk: %v", err)
		return
	}

	file, err := cc.FindFile(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zero.Log.Error().Msgf("error finding finding file in CC: %v", err)
		return
	}

	info, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zero.Log.Error().Msgf("error getting file info: %v", err)
		return
	}

	path := filepath.Join(tmpDir, info.Name())
	dst, err := os.Create(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zero.Log.Error().Msgf("error creating dest tmp file: %v", err)
		return
	}

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zero.Log.Error().Msgf("error writing file to tmp: %v", err)
		return
	}

	//-- handle doc to pdf conversion
	ext := filepath.Ext(info.Name())
	match, err := regexp.Match(`(doc|docx|odt)`, []byte(ext))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zero.Log.Error().Msgf("error parsing file extension: %v", err)
		return
	}

	if match {
		libreoffice, err := exec.LookPath("libreoffice")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			zero.Log.Error().Msgf("error finding libreoffice: %v", err)
			return
		}

		cmd := exec.Command(libreoffice, "--headless", "--convert-to", "pdf", "--outdir", tmpDir, path)

		err = cmd.Run()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			zero.Log.Error().Msgf("error converting file to pdf: %v", err)
			return
		}

		path = strings.TrimSuffix(path, filepath.Ext(path)) + ".pdf"

	}

	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["path"] = path
	body, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zero.Log.Error().Msgf("error marshalling response to json: %v", err)
		return
	}
	fmt.Fprint(w, string(body))
}

func handleResource(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		zero.Log.Warn().Msgf("Method not allowed: %s", r.Method)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/resource/")
	cartridge := r.FormValue("cartridge")
	zero.Log.Info().Msgf("GET handleResource id: %v cartridge %v", id, cartridge)

	inputFile := filepath.Join(uploadsDir, cartridge)
	cc, err := commoncartridge.Load(inputFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zero.Log.Error().Msgf("error loading CC from disk: %v", err)
		return
	}

	file, err := cc.Find(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zero.Log.Error().Msgf("error finding resource in CC: %v", err)
		return
	}

	data, err := json.Marshal(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zero.Log.Error().Msg("error marshalling to json")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(data))
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

	err = os.MkdirAll(uploadsDir, os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// todo replace with Filepath.Join
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

	inputFile := filepath.Join(uploadsDir, fileHeader.Filename)
	cc, err := commoncartridge.Load(inputFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zero.Log.Error().Msg("error loading CC from filesystem")
		return
	}

	obj, err := cc.AsObject()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zero.Log.Error().Msg("error parsing manifest into JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["data"] = string(obj)

	fr, err := cc.Resources()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zero.Log.Error().Msgf("error getting resources: %v", err)
		return
	}
	sfr, err := json.Marshal(fr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		zero.Log.Error().Msgf("error getting resources: %v", err)
		return
	}
	resp["resources"] = string(sfr)

	body, _ := json.Marshal(resp)
	w.Write(body)
}
