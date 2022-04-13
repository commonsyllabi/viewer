package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
