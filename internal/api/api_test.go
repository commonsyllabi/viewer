package api

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const singleTestFile = "../../pkg/commoncartridge/test_files/test_01.imscc"

func TestLoadConfig(t *testing.T) {
	err := conf.loadConfig("../../internal/api/config.yml")

	if err != nil {
		t.Errorf("error loading conf file: %v", err)
	}
}

func TestHandlePing(t *testing.T) {
	conf.defaults()
	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != 200 {
		t.Errorf("expected 200, got %v", res.Code)
	}
	if res.Body.String() != "pong" {
		t.Errorf("expected pong, got: %v", res.Body.String())
	}
}

func TestHandleUpload(t *testing.T) {
	conf.defaults()
	router := setupRouter()

	body, writer := createFormData("cartridge", singleTestFile, t)
	req, _ := http.NewRequest(http.MethodPost, "/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	result := res.Result()
	defer result.Body.Close()

	if res.Code != 200 {
		t.Errorf("expected 200 response code, got %d", res.Code)
	}

	var response map[string]string
	json.Unmarshal(res.Body.Bytes(), &response)
	_, exists := response["data"]
	if !exists {
		t.Errorf("Expected to have a JSON object with a \"data\" field")
	}

	_, exists = response["items"]
	if !exists {
		t.Errorf("Expected to have a JSON object with a \"items\" field")
	}

	_, exists = response["resources"]
	if !exists {
		t.Errorf("Expected to have a JSON object with a \"resources\" field")
	}
}

func TestHandleFile(t *testing.T) {
	TestHandleUpload(t)
	router := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/file/i3755487a331b36c76cec8bbbcdb7cc66?cartridge=test_01.imscc", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	result := res.Result()

	if res.Code != 200 {
		t.Errorf("expected 200 response code, got %d", res.Code)
	}

	var response map[string]string
	json.Unmarshal(res.Body.Bytes(), &response)
	_, exists := response["path"]
	if !exists {
		t.Errorf("Expected to have a JSON object with a \"path\" field")
	}

	defer result.Body.Close()
}

func TestHandleResource(t *testing.T) {
	TestHandleUpload(t)
	router := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/resource/i3755487a331b36c76cec8bbbcdb7cc66?cartridge=test_01.imscc", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	result := res.Result()

	if res.Code != 200 {
		t.Errorf("expected 200 response code, got %d", res.Code)
	}

	defer result.Body.Close()
}

func createFormData(fieldName, fileName string, t *testing.T) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	var fw io.Writer
	file := mustOpen(fileName)

	fw, err := w.CreateFormFile(fieldName, file.Name())
	if err != nil {
		t.Errorf("Cannot create form file %s", file.Name())
	}

	_, err = io.Copy(fw, file)
	if err != nil {
		t.Error("error copying file")
	}
	w.Close()

	return b, w
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
