package mailer

import (
	"bytes"
	"crypto/sha256"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/commonsyllabi/viewer/internal/api/models"
	"github.com/gin-gonic/gin"
)

const singleTestFile = "../../../pkg/commoncartridge/test_files/test_01.imscc"

func TestGetMagicLink(t *testing.T) {
	mustSeedDB(t)
	gin.SetMode(gin.TestMode)

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	w.WriteField("id", "1")
	w.WriteField("email", "pierre.depaz@gmail.com")
	w.Close()

	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", w.FormDataContentType())
	c.Request.Body = io.NopCloser(&body)
	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: "1",
		},
	}

	HandleMagicLink(c)
	if res.Code != 200 {
		t.Errorf("Expected 200, got %v - %v", res.Code, res)
	}
}

func mustSeedDB(t *testing.T) {
	databaseTestURL := os.Getenv("DATABASE_TEST_URL")
	if databaseTestURL == "" {
		databaseTestURL = "postgres://cosyl:cosyl@localhost:5432/test"
	}
	_, err := models.InitDB(databaseTestURL, "../models/fixtures")
	if err != nil {
		panic(err)
	}

	// err = models.SetupTables(true)
	// if err != nil {
	// 	panic(err)
	// }

	syll := models.Syllabus{Title: "Test Title 1", Description: "Test Description 1"}
	_, err = models.AddNewSyllabus(&syll)
	if err != nil {
		panic(err)
	}

	hasher := sha256.New()
	hasher.Write([]byte(syll.Title))
	token := models.MagicToken{Token: hasher.Sum(nil), SyllabusTokenID: syll.ID}
	token, err = models.AddNewToken(&token)
	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadFile(singleTestFile)
	if err != nil {
		t.Error(err)
	}
	att := models.Attachment{
		Name: "test_01.imscc",
		File: bytes,
		Type: "zip",
	}

	_, err = models.AddNewAttachment(&att)
	if err != nil {
		panic(err)
	}
}
