package handlers

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
	"testing"

	"github.com/commonsyllabi/viewer/internal/api/models"
	"github.com/gin-gonic/gin"
)

const singleTestFile = "../../../pkg/commoncartridge/test_files/test_01.imscc"

func TestAllSyllabi(t *testing.T) {
	mustSeedDB(t)
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	AllSyllabi(c)
	if res.Code != 200 {
		t.Errorf("Expected 200, got %v", res.Code)
	}
}

func TestNewSyllabus(t *testing.T) {
	mustSeedDB(t)

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	w.WriteField("title", "Test Syllabus Handling")
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

	NewSyllabus(c)

	if res.Code != 200 {
		t.Errorf("Expected 200, got %v: %v", res.Code, res.Body)
	}
}

func TestNewSyllabusWrongValue(t *testing.T) {
	mustSeedDB(t)

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

	NewSyllabus(c)

	if res.Code != 400 {
		t.Errorf("Expected 200, got %v: %v", res.Code, res.Body)
	}
}

func TestNewSyllabusSingleAttachment(t *testing.T) {
	mustSeedDB(t)

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

	NewSyllabus(c)

	if res.Code != 200 {
		t.Errorf("Expected 200, got %v : %s", res.Code, res.Body.String())
	}
}

func TestUpdateSyllabus(t *testing.T) {
	mustSeedDB(t)

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

	UpdateSyllabus(c)

	if res.Code != 200 {
		t.Errorf("Expected 200, got %v: %v", res.Code, res.Body)
	}

	var syll models.Syllabus
	err := c.Bind(&syll)
	if err != nil {
		t.Error(err)
	}

	if syll.Title != "Updated" {
		t.Errorf("Expected to have updated title, got %v", syll.Title)
	}
}

func TestUpdateSyllabusPartial(t *testing.T) {
	mustSeedDB(t)

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

	UpdateSyllabus(c)

	if res.Code != 200 {
		t.Errorf("Expected 200, got %v: %v", res.Code, res.Body)
	}

	var syll models.Syllabus
	err := c.Bind(&syll)
	if err != nil {
		t.Error(err)
	}

	if syll.Description != "Updated" {
		t.Errorf("Expected to have updated description, got %v", syll.Description)
	}
}

func TestGetSyllabus(t *testing.T) {
	mustSeedDB(t)
	gin.SetMode(gin.TestMode)

	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Request.Method = "GET"
	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: "1",
		},
	}
	GetSyllabus(c)
	if res.Code != 200 {
		t.Errorf("Expected 200, got %v", res.Code)
	}
}

func TestDisplayMagicLink(t *testing.T) {
	mustSeedDB(t)
	gin.SetMode(gin.TestMode)
	token, err := models.GetTokenSyllabus(1)
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
			Value: "1",
		},
	}
	c.Request.URL, _ = url.Parse("?token=" + base64.URLEncoding.EncodeToString(token.Token))

	DisplayMagicLink(c)
	if res.Code != 200 {
		t.Errorf("Expected 200, got %v: %v", res.Code, res)
	}
}

func TestDeleteSyllabus(t *testing.T) {
	mustSeedDB(t)

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

	DeleteSyllabus(c)

	if res.Code != 200 {
		t.Errorf("Expected 200, got %v: %v", res.Code, res.Body)
	}

	//-- todo check the actual id returned
	// var body struct{ id string }
	// err := c.Bind(&body)
	// if err != nil {
	// 	t.Error(err)
	// }

	// fmt.Println(body)

	// if body.id != "1" {
	// 	t.Errorf("Expected to have deleted id 1, got %v", body.id)
	// }
}

func mustSeedDB(t *testing.T) {
	_, err := models.InitDB("postgres://test:test@localhost:5432/test")
	if err != nil {
		panic(err)
	}

	err = models.SetupTables(true)
	if err != nil {
		panic(err)
	}

	syll := models.Syllabus{Title: "Test Title 1", Description: "Test Description 1"}
	_, err = models.AddNewSyllabus(&syll)
	if err != nil {
		panic(err)
	}

	hasher := sha256.New()
	hasher.Write([]byte(syll.Title))
	token := models.MagicToken{Token: hasher.Sum(nil), SyllabusID: syll.ID}
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

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
