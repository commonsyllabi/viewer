package api

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const singleTestFile = "../../pkg/commoncartridge/test_files/test_01.imscc"

func TestHandlePing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	res := httptest.NewRecorder()
	handlePing(res, req)
	result := res.Result()

	defer result.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error in reading response body: %v", err)
	}

	if string(data) != "pong" {
		t.Errorf("expected response to be pong, got %v", string(data))
	}
}

func TestHandleUpload(t *testing.T) {
	body, writer := createFormData("cartridge", singleTestFile, t)

	req := httptest.NewRequest(http.MethodPost, "/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res := httptest.NewRecorder()

	handleUpload(res, req)
	result := res.Result()
	defer result.Body.Close()

	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("error in reading response body: %v", err)
	}

	if res.Code != 200 {
		t.Errorf("expected 200 response code, got %d", res.Code)
	}
}

func TestHandleFile(t *testing.T) {
	// uploadTestCartridge(t)
	TestHandleUpload(t)

	req := httptest.NewRequest(http.MethodGet, "/file/i3755487a331b36c76cec8bbbcdb7cc66?cartridge=test_01.imscc", nil)
	res := httptest.NewRecorder()
	handleFile(res, req)
	result := res.Result()

	if res.Code != 200 {
		t.Errorf("expected 200 response code, got %d", res.Code)
	}

	defer result.Body.Close()
}

func TestHandleResource(t *testing.T) {
	TestHandleUpload(t)

	req := httptest.NewRequest(http.MethodGet, "/resource/i3755487a331b36c76cec8bbbcdb7cc66?cartridge=test_01.imscc", nil)
	res := httptest.NewRecorder()
	handleResource(res, req)
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
