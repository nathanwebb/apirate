package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	url := fmt.Sprintf("/%s/health", apiVersion)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleEvent)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := `{"status":0`
	if !strings.HasPrefix(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetResults(t *testing.T) {

	for _, ts := range slackts {
		url := fmt.Sprintf("/%s/event?slackts=%s", apiVersion, ts)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handleEvent)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		expected := `{"status":0`
		if !strings.HasPrefix(rr.Body.String(), expected) {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
		fmt.Printf("%s\n", rr.Body.String())
	}

}
