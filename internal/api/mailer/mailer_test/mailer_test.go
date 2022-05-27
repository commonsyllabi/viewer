package mailer_test

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/commonsyllabi/viewer/internal/api/mailer"
	"github.com/commonsyllabi/viewer/internal/api/models"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var syllabusID string

func setup(t *testing.T) func(t *testing.T) {
	mustSeedDB(t)
	return func(t *testing.T) {
		models.RemoveFixtures(t)
	}
}

func TestMagicLink(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	t.Run("Test get magic link", func(t *testing.T) {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		w.WriteField("id", syllabusID)
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
				Value: syllabusID,
			},
		}

		mailer.HandleMagicLink(c)
		fmt.Println(res.Body.String())
		assert.Equal(t, res.Code, http.StatusOK)
	})

}

func mustSeedDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	databaseTestURL := os.Getenv("DATABASE_TEST_URL")
	if databaseTestURL == "" {
		databaseTestURL = "postgres://cosyl:cosyl@localhost:5432/test"
	}
	_, err := models.InitDB(databaseTestURL)
	require.Nil(t, err)

	syll := models.Syllabus{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       "Test Title for Token",
		Description: "Test Description for Token",
	}
	s, err := models.AddNewSyllabus(&syll)
	require.Nil(t, err)

	syllabusID = strconv.Itoa(int(s.ID))

	hasher := sha256.New()
	hasher.Write([]byte(syll.Title))
	token := models.MagicToken{
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Token:           hasher.Sum(nil),
		SyllabusTokenID: s.ID,
	}
	token, err = models.AddNewToken(&token)

	assert.Nil(t, err)
}
