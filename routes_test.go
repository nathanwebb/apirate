package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetKeysWithoutKeystoreFails(t *testing.T) {
	req, r := GetKeysWithKeystore(t, "")
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusBad := w.Code != http.StatusOK
		return statusBad
	})
}

func TestGetKeysWithKeystore(t *testing.T) {
	os.Setenv("KEYSTORE", "file:///keystore_test.json")
	req, r := GetKeysWithKeystore(t, "")
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		p, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		keys := []key{}
		err = json.Unmarshal(p, &keys)
		if err != nil {
			t.Error(err.Error())
		}
		if len(keys) == 0 {
			t.Error("no keys returned")
		}
		if keys[0].Type != "ssh" {
			t.Errorf("incorrect keys returned")
		}

		statusOK := w.Code == http.StatusOK
		return statusOK
	})
}

func TestGetSingleKeyWithKeystore(t *testing.T) {
	os.Setenv("KEYSTORE", "file:///keystore_test.json")
	req, r := GetKeysWithKeystore(t, "2")
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		p, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		keys := []key{}
		err = json.Unmarshal(p, &keys)
		if err != nil {
			t.Error(err.Error())
		}
		if len(keys) != 1 {
			t.Error("incorrect number of keys returned")
		}
		if keys[0].Type != "ssh" {
			t.Errorf("incorrect keys returned")
		}

		statusOK := w.Code == http.StatusOK
		return statusOK
	})
}

func TestCreateKeys(t *testing.T) {
	keyfile := "keystore_test_dummy.json"
	os.Setenv("KEYSTORE", "file:///"+keyfile)
	apiurl := "/api/v1/keys"
	w := httptest.NewRecorder()

	router := gin.Default()
	router.POST(apiurl, createKey)

	newKey := url.Values{}
	newKey.Set("Type", "ssh")
	req, _ := http.NewRequest("POST", apiurl, strings.NewReader(newKey.Encode()))
	req.Header.Add("Content-Length", strconv.Itoa(len(newKey.Encode())))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}
	_, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fail()
	}
	if err = os.Remove(keyfile); err != nil {
		t.Errorf("failed to remove keystore: %s", err)
	}
}

func GetKeysWithKeystore(t *testing.T, keyID string) (*http.Request, *gin.Engine) {
	r := gin.Default()
	apiurl := "/api/v1/keys"
	if keyID != "" {
		apiurl = fmt.Sprintf("%s?%s", apiurl, keyID)
	}
	r.GET(apiurl, getKeys)
	req, _ := http.NewRequest("GET", apiurl, nil)
	return req, r
}
