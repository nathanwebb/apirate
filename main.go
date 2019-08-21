package main

import (
	"net/http"
	"os"
)

const apiVersion = "v1"

func main() {
	keystore := os.Getenv("KEYSTORE")
	loadKeys(keystore)
}

func handleEvent(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
	case "GET":
	}
}
