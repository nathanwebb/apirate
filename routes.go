package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

const apiVersion = "v1"

var keystore = os.Getenv("KEYSTORE")

func initialiseRoutes(router *gin.Engine) {
	keysRoutes := router.Group(fmt.Sprintf("/api/%s/keys", apiVersion))
	keysRoutes.GET("/", getKeys)
	keysRoutes.GET("/:id", getSingleKey)
	keysRoutes.POST("/", createKey)
	keysRoutes.DELETE("/", deleteAllKeys)
	keysRoutes.DELETE("/kes:id", deleteKey)
}

func getKeys(c *gin.Context) {
	loadKeys(keystore)

}

func createKey(c *gin.Context) {

}

func getSingleKey(c *gin.Context) {

}

func deleteAllKeys(c *gin.Context) {

}

func deleteKey(c *gin.Context) {

}
