package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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
	}{
		{"KEYSTORE", "file:///var/local/apirate/default_keystore.json"},
		{"COMMANDSTORE", "file:///var/local/apirate/default_commandstore.json"},
	}
	for _, v := range vars {
		if tv := os.Getenv(v.key); tv == "" {
			os.Setenv(v.key, v.value)
		}
		fmt.Printf("setting %s to %s\n", v.key, os.Getenv(v.key))
	}
}
