package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetKeysWithoutKeystoreFails(t *testing.T) {
	os.Setenv("KEYSTORE", "")
	req, r := GetKeysWithKeystore(t, "")
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusBad := w.Code != http.StatusOK
		return statusBad
	})
}

func TestGetKeysWithKeystore(t *testing.T) {
	err := createTestKeyStore("keystore_test.json")
	if err != nil {
		t.Fatal(err.Error())
	}
	req, r := GetKeysWithKeystore(t, "")
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		p, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		keys := []key{}
		err = json.Unmarshal(p, &keys)
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(keys) == 0 {
			t.Fatal("no keys returned")
		}
		if keys[0].Type != "ssh" {
			t.Fatal("incorrect keys returned")
		}

		statusOK := w.Code == http.StatusOK
		return statusOK
	})
}

func TestGetSingleKeyWithKeystore(t *testing.T) {
	err := createTestKeyStore("keystore_test.json")
	if err != nil {
		t.Fatal(err.Error())
	}
	req, r := GetKeysWithKeystore(t, "2")
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		p, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		responseKey := key{}
		err = json.Unmarshal(p, &responseKey)
		if err != nil {
			t.Error(err.Error())
		}
		if responseKey.Type != "ssh" {
			t.Errorf("incorrect keys returned")
		}

		statusOK := w.Code == http.StatusOK
		return statusOK
	})
}

func TestCreateKeys(t *testing.T) {
	keyfile := "keystore_test_dummy.json"
	os.Setenv("KEYSTORE", keyfile)
	apiurl := "/api/v1/keys"

	router := gin.Default()
	router.POST(apiurl, createKey)

	newKey := []byte(`{"type": "ssh"}`)
	req, _ := http.NewRequest("POST", apiurl, bytes.NewBuffer(newKey))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
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

func TestGetResults(t *testing.T) {
	commandsfile := "commands_config_test.json"
	os.Setenv("COMMANDSTORE", "file:///"+commandsfile)
	apiURL := "/api/v1/results"
	r := gin.Default()
	r.GET(apiURL, getResults)

	tests := []struct {
		request    string
		statusCode int
	}{
		{"?name=local%20ping&ip=127.0.0.1", http.StatusOK},
		{"?name=local%20ping&ip=127.0.0.1%3Btouch%20test.txt", http.StatusOK},
		{"?name=cat&ip=/etc/passwd", http.StatusBadRequest},
	}
	for _, c := range tests {
		req, _ := http.NewRequest("GET", apiURL+c.request, nil)
		testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
			_, err := os.Open("malicious_test.txt")
			if err == nil || !os.IsNotExist(err) {
				t.Errorf("malicious attack was successful")
			}
			statusOK := w.Code == c.statusCode
			if !statusOK {
				t.Errorf("expected status code: %d, got %d", c.statusCode, w.Code)
			}
			log.Println(c.statusCode)
			if c.statusCode != http.StatusOK {
				return statusOK
			}
			p, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Errorf("%s", err.Error())
			}
			result := command{}
			err = json.Unmarshal(p, &result)
			if err != nil {
				t.Error(err.Error())
			}
			if result.Stdout == "" && result.Stderr == "" && result.Error == "" {
				t.Errorf("no output returned: %s", c.request)
			}
			t.Log(result.Stdout)
			return statusOK
		})
	}
}

func GetKeysWithKeystore(t *testing.T, keyID string) (*http.Request, *gin.Engine) {
	r := gin.Default()
	apiUrl := "/api/v1/keys"
	requestUrl := apiUrl
	if keyID != "" {
		requestUrl = requestUrl + "/" + keyID
		apiUrl = apiUrl + "/:id"
	}
	r.GET(apiUrl, getKeys)
	req, _ := http.NewRequest("GET", requestUrl, nil)
	return req, r
}
