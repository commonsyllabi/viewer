package handlers_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/commonsyllabi/viewer/internal/api/handlers"
	"github.com/commonsyllabi/viewer/internal/api/models"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const singleTestFile = "../../../../pkg/commoncartridge/test_files/test_01.imscc"

var syllabusID string

func setup(t *testing.T) func(t *testing.T) {
	mustSeedDB(t)
	gin.SetMode(gin.TestMode)
	return func(t *testing.T) {
	}
}

func TestSyllabusHandler(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	t.Run("Test get all syllabi", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		handlers.AllSyllabi(c)
		assert.Equal(t, res.Code, http.StatusOK)
	})

	t.Run("Test add new syllabus", func(t *testing.T) {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		w.WriteField("title", "Test Syllabus Handling")
		w.WriteField("description", "This is a test for the syllabus handling")
		w.WriteField("email", "correct@host.com")
		w.Close()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Request.Method = "POST"
		c.Request.Header.Set("Content-Type", w.FormDataContentType())
		c.Request.Body = io.NopCloser(&body)

		handlers.NewSyllabus(c)

		assert.Equal(t, res.Code, http.StatusOK)
	})

	t.Run("Test new syllabus with wrong values", func(t *testing.T) {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		w.WriteField("title", "")
		w.WriteField("description", "This is a test for the syllabus handling")
		w.WriteField("email", "name@host.com")
		w.Close()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Request.Method = "POST"
		c.Request.Header.Set("Content-Type", w.FormDataContentType())
		c.Request.Body = io.NopCloser(&body)

		handlers.NewSyllabus(c)

		assert.Equal(t, res.Code, http.StatusBadRequest)
	})

	t.Run("Test new syllabus with email too short", func(t *testing.T) {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		w.WriteField("title", "")
		w.WriteField("description", "This is a test for the syllabus handling")
		w.WriteField("email", "a@b.fr")
		w.Close()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Request.Method = "POST"
		c.Request.Header.Set("Content-Type", w.FormDataContentType())
		c.Request.Body = io.NopCloser(&body)

		handlers.NewSyllabus(c)

		assert.Equal(t, res.Code, http.StatusBadRequest)
	})

	t.Run("Test new syllabus invalid email", func(t *testing.T) {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		w.WriteField("title", "")
		w.WriteField("description", "This is a test for the syllabus handling")
		w.WriteField("email", "obviously not an email")
		w.Close()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Request.Method = "POST"
		c.Request.Header.Set("Content-Type", w.FormDataContentType())
		c.Request.Body = io.NopCloser(&body)

		handlers.NewSyllabus(c)

		assert.Equal(t, res.Code, http.StatusBadRequest)
	})

	t.Run("Test new syllabus single attachment", func(t *testing.T) {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		w.WriteField("title", "Test Syllabus Handling")
		w.WriteField("description", "This is a test for the syllabus handling")
		w.WriteField("email", "name@host.com")

		var fw io.Writer
		file := mustOpen(singleTestFile)
		fw, err := w.CreateFormFile("attachments[]", file.Name())
		if err != nil {
			t.Errorf("Cannot create form file %s", file.Name())
		}

		if _, err = io.Copy(fw, file); err != nil {
			t.Error("error copying file")
		}
		w.Close()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Request.Method = "POST"
		c.Request.Header.Set("Content-Type", w.FormDataContentType())
		c.Request.Body = io.NopCloser(&body)

		handlers.NewSyllabus(c)

		assert.Equal(t, res.Code, http.StatusOK)
	})

	t.Run("Test update syllabus", func(t *testing.T) {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		w.WriteField("title", "Updated")
		w.WriteField("description", "This is a test for the syllabus handling")
		w.Close()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Request.Method = "PATCH"
		c.Request.Header.Set("Content-Type", w.FormDataContentType())
		c.Request.Body = io.NopCloser(&body)
		c.Params = []gin.Param{
			{
				Key:   "id",
				Value: "1",
			},
		}

		handlers.UpdateSyllabus(c)
		assert.Equal(t, res.Code, http.StatusOK)

		var syll models.Syllabus
		err := c.Bind(&syll)
		require.Nil(t, err)
		assert.Equal(t, syll.Title, "Updated")
	})

	t.Run("Test update syllabus partial", func(t *testing.T) {
		var body bytes.Buffer
		w := multipart.NewWriter(&body)
		w.WriteField("description", "Updated")
		w.Close()

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Request.Method = "PATCH"
		c.Request.Header.Set("Content-Type", w.FormDataContentType())
		c.Request.Body = io.NopCloser(&body)
		c.Params = []gin.Param{
			{
				Key:   "id",
				Value: "1",
			},
		}

		handlers.UpdateSyllabus(c)
		assert.Equal(t, res.Code, http.StatusOK)

		var syll models.Syllabus
		err := c.Bind(&syll)
		require.Nil(t, err)
		assert.Equal(t, syll.Description, "Updated")
	})

	t.Run("Test get syllabus", func(t *testing.T) {

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Request.Method = "GET"
		c.Params = []gin.Param{
			{
				Key:   "id",
				Value: syllabusID,
			},
		}

		handlers.GetSyllabus(c)
		assert.Equal(t, res.Code, http.StatusOK)
	})

	t.Run("Test display magic link", func(t *testing.T) {
		id, _ := strconv.Atoi(syllabusID)
		token, err := models.GetTokenSyllabus(id)
		if err != nil {
			t.Error(err)
		}

		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Request.Method = "GET"
		c.Params = []gin.Param{
			{
				Key:   "id",
				Value: syllabusID,
			},
		}
		c.Request.URL, _ = url.Parse("?token=" + base64.URLEncoding.EncodeToString(token.Token))

		handlers.DisplayMagicLink(c)
		assert.Equal(t, res.Code, http.StatusOK)
	})

	t.Run("Test delete syllabus", func(t *testing.T) {
		res := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(res)
		c.Request = &http.Request{
			Header: make(http.Header),
		}

		c.Request.Method = "DELETE"
		c.Params = []gin.Param{
			{
				Key:   "id",
				Value: "1",
			},
		}

		handlers.DeleteSyllabus(c)
		assert.Equal(t, res.Code, http.StatusOK)
		//-- todo check the actual id returned
	})

	// todo test syllabus with already existing email
}

func mustSeedDB(t *testing.T) {
	databaseTestURL := os.Getenv("DATABASE_TEST_URL")
	if databaseTestURL == "" {
		databaseTestURL = "postgres://cosyl:cosyl@localhost:5432/test"
	}
	db, err := models.InitDB(databaseTestURL)
	models.RunFixtures(db, models.Basepath+"/../fixtures")
	require.Nil(t, err)

	syll := models.Syllabus{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       "Test Title 1",
		Description: "Test Description 1",
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
		SyllabusTokenID: syll.ID,
	}
	token, err = models.AddNewToken(&token)
	require.Nil(t, err)

	bytes, err := ioutil.ReadFile(singleTestFile)
	if err != nil {
		t.Error(err)
	}
	att := models.Attachment{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      "test_01.imscc",
		File:      bytes,
		Type:      "zip",
	}

	_, err = models.AddNewAttachment(&att)
	require.Nil(t, err)
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
