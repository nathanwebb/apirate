package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetKeysWithoutKeystoreFails(t *testing.T) {
	req, r := GetKeysWithKeystore(t)
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusBad := w.Code != http.StatusOK
		return statusBad
	})
}

func TestGetKeysWithKeystore(t *testing.T) {
	os.Setenv("KEYSTORE", "file:///keystore_test.json")
	req, r := GetKeysWithKeystore(t)
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		return statusOK
	})
}

func GetKeysWithKeystore(t *testing.T) (*http.Request, *gin.Engine) {
	r := gin.Default()
	r.GET("/api/v1/keys", getKeys)
	req, _ := http.NewRequest("GET", "/api/v1/keys", nil)
	return req, r
}
