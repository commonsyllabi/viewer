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
	"path/filepath"
	"testing"

	"github.com/commonsyllabi/viewer/internal/api/models"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

const singleTestFile = "../../tests/samples/test_01.imscc"

var router *gin.Engine

func setup(t *testing.T) func(t *testing.T) {
	err := os.MkdirAll(filepath.Join(conf.TmpDir, conf.FilesDir), os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	gin.SetMode(gin.TestMode)
	router = mustSetupRouter()

	return func(t *testing.T) {
		t.Log("tearing down api")
	}
}

func TestApi(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	t.Run("Testing ping", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, res.Code, http.StatusOK, "expected 200, got %v", res.Code)
		assert.Equal(t, res.Body.String(), "pong", "expected pong, got: %v", res.Body.String())
	})

	t.Run("Testing upload", func(t *testing.T) {
		body, writer := createFormData("cartridge", singleTestFile, t)
		req, _ := http.NewRequest(http.MethodPost, "/api/upload", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		result := res.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusOK, res.Code, "expected 200, got %v", res.Code)

		var response map[string]string
		json.Unmarshal(res.Body.Bytes(), &response)

		assert.NotEmpty(t, response["data"])
		assert.NotEmpty(t, response["items"])
		assert.NotEmpty(t, response["resources"])
	})

	t.Run("Test upload no field", func(t *testing.T) {
		body, writer := createFormData("bad_cartridge", singleTestFile, t)
		req, _ := http.NewRequest(http.MethodPost, "/api/upload", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		result := res.Result()
		defer result.Body.Close()
		assert.Equal(t, res.Code, http.StatusBadRequest, "expected 400, got %v", res.Code)
	})

	t.Run("Test upload no file", func(t *testing.T) {
		body, writer := createFormData("cartridge", "", t)
		req, _ := http.NewRequest(http.MethodPost, "/api/upload", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		result := res.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	//-- todo check for content-type in header
	t.Run("Test handle file", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/file/i3755487a331b36c76cec8bbbcdb7cc66?cartridge=test_01.imscc", nil)
		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		result := res.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Test handle file without ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/file/WRONG-FILE?cartridge=test_01.imscc", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusInternalServerError, "expected 500, got %v", res.Code)
	})

	t.Run("Test handle file without cartridge", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/file/i3755487a331b36c76cec8bbbcdb7cc66?cartridge=WRONG-CARTRIDGE.imscc", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusInternalServerError, "expected 500, got %v", res.Code)
	})

	//-- todo check headers for content-type
	t.Run("Test handle resource", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/resource/i3755487a331b36c76cec8bbbcdb7cc66?cartridge=test_01.imscc", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusOK, "expected 200, got %v", res.Code)
	})

	t.Run("Test resource no cartridge", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/resource/i3755487a331b36c76cec8bbbcdb7cc66?cartridge=MISSING_CARTRIDGE", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusInternalServerError, "expected 500, got %v", res.Code)
	})

	t.Run("Test resource no ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/resource/MALFORMED_ID?cartridge=test_01.imscc", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusInternalServerError, "expected 500, got %v", res.Code)
	})
}

func TestLoadConfig(t *testing.T) {
	err := conf.LoadConf("../../internal/api/config.yml")
	assert.NotNil(t, err)
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

		if _, err = io.Copy(fw, file); err != nil {
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

func mustUploadFile(t *testing.T, router *gin.Engine) {
	body, writer := createFormData("cartridge", singleTestFile, t)
	req, _ := http.NewRequest(http.MethodPost, "/api/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
}

func mustSetupRouter() *gin.Engine {
	conf.DefaultConf()
	conf.TemplatesDir = "../../internal/api/templates"
	conf.FixturesDir = "../../internal/api/models/fixtures"

	databaseTestURL := os.Getenv("DATABASE_TEST_URL")
	if databaseTestURL == "" {
		databaseTestURL = "postgres://cosyl:cosyl@localhost:5432/viewer-test"
		fmt.Printf("didn't get db test url from env, defaulting t %v\n", databaseTestURL)

	}

	db, err := models.InitDB(databaseTestURL)
	models.RunFixtures(db, conf.FixturesDir)
	if err != nil {
		panic(err)
	}
	router, err := setupRouter()
	if err != nil {
		panic(err)
	}
	return router
}
