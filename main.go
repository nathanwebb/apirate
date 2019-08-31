package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	setDefaultEnv()
}

func main() {
	router := gin.Default()
	initialiseRoutes(router)
	router.Run()
}

func setDefaultEnv() {
	vars := []struct {
		key   string
		value string
	}{{"KEYSTORE", "file:///default_keystore.json"}}
	for _, v := range vars {
		if tv := os.Getenv(v.key); tv == "" {
			os.Setenv(v.key, v.value)
		}
	}
}
