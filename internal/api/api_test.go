package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/commonsyllabi/viewer/internal/api/models"
	"github.com/gin-gonic/gin"
)

const singleTestFile = "../../pkg/commoncartridge/test_files/test_01.imscc"

func TestLoadConfig(t *testing.T) {
	err := conf.load("../../internal/api/config.yml")

	if err != nil {
		t.Errorf("error loading conf file: %v", err)
	}
}

func TestHandlePing(t *testing.T) {
	router := mustSetupRouter(false)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("expected 200, got %v", res.Code)
	}
	if res.Body.String() != "pong" {
		t.Errorf("expected pong, got: %v", res.Body.String())
	}
}

func TestHandleUpload(t *testing.T) {
	router := mustSetupRouter(false)

	body, writer := createFormData("cartridge", singleTestFile, t)
	req, _ := http.NewRequest(http.MethodPost, "/api/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	result := res.Result()
	defer result.Body.Close()

	if res.Code != http.StatusOK {
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

func TestHandleUploadNoField(t *testing.T) {
	router := mustSetupRouter(false)

	body, writer := createFormData("bad_cartridge", singleTestFile, t)
	req, _ := http.NewRequest(http.MethodPost, "/api/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	result := res.Result()
	defer result.Body.Close()

	if res.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request response code, got %d", res.Code)
	}
}

func TestHandleUploadNoFile(t *testing.T) {
	router := mustSetupRouter(false)

	body, writer := createFormData("cartridge", "", t)
	req, _ := http.NewRequest(http.MethodPost, "/api/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	result := res.Result()
	defer result.Body.Close()

	if res.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request response code, got %d", res.Code)
	}
}

func TestHandleFile(t *testing.T) {
	TestHandleUpload(t)
	router := mustSetupRouter(false)

	req, _ := http.NewRequest(http.MethodGet, "/api/file/i3755487a331b36c76cec8bbbcdb7cc66?cartridge=test_01.imscc", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	result := res.Result()

	fmt.Println(res.Body)

	if res.Code != http.StatusOK {
		t.Errorf("expected 200 response code, got %d", res.Code)
	}

	defer result.Body.Close()
}

func TestHandleFileNoID(t *testing.T) {
	TestHandleUpload(t)
	router := mustSetupRouter(false)

	req, _ := http.NewRequest(http.MethodGet, "/api/file/WRONG-FILE?cartridge=test_01.imscc", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 response code, got %d", res.Code)
	}
}

func TestHandleFileNoCartridge(t *testing.T) {
	TestHandleUpload(t)
	router := mustSetupRouter(false)

	req, _ := http.NewRequest(http.MethodGet, "/api/file/i3755487a331b36c76cec8bbbcdb7cc66?cartridge=WRONG-CARTRIDGE.imscc", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 response code, got %d", res.Code)
	}
}

func TestHandleResource(t *testing.T) {
	TestHandleUpload(t)
	router := mustSetupRouter(false)

	req, _ := http.NewRequest(http.MethodGet, "/api/resource/i3755487a331b36c76cec8bbbcdb7cc66?cartridge=test_01.imscc", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	result := res.Result()

	if res.Code != http.StatusOK {
		t.Errorf("expected 200 response code, got %d", res.Code)
	}

	defer result.Body.Close()
}

func TestHandleResourceNoCartridge(t *testing.T) {
	TestHandleUpload(t)
	router := mustSetupRouter(false)

	req, _ := http.NewRequest(http.MethodGet, "/api/resource/i3755487a331b36c76cec8bbbcdb7cc66?cartridge=MISSING_CARTRIDGE", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	result := res.Result()

	if res.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 response code, got %d", res.Code)
	}

	defer result.Body.Close()
}

func TestHandleResourceNoID(t *testing.T) {
	TestHandleUpload(t)
	router := mustSetupRouter(false)

	req, _ := http.NewRequest(http.MethodGet, "/api/resource/MALFORMED_ID?cartridge=test_01.imscc", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	result := res.Result()

	if res.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 response code, got %d", res.Code)
	}

	defer result.Body.Close()
}

func createFormData(fieldName, fileName string, t *testing.T) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	if fileName != "" {
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
	}

	return b, w
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

func mustSetupRouter(debug bool) *gin.Engine {
	conf.defaults()
	conf.TemplatesDir = "../../internal/api/templates"

	//-- todo see db_test
	_, err := models.InitDB("postgres://test:test@localhost:5432/test")
	if err != nil {
		panic(err)
	}
	router, err := setupRouter(debug)
	if err != nil {
		panic(err)
	}
	return router
}
