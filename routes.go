package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

const apiVersion = "v1"

func initialiseRoutes(router *gin.Engine) {
	keysRoutes := router.Group(fmt.Sprintf("/api/%s/keys", apiVersion))
	keysRoutes.GET("/", getKeys)
	keysRoutes.GET("/:id", getKeys)
	keysRoutes.POST("/", createKey)
	keysRoutes.DELETE("/", deleteAllKeys)
	keysRoutes.DELETE("/kes:id", deleteKey)
}

func getKeys(c *gin.Context) {
	keystore := os.Getenv("KEYSTORE")
	keys, err := loadKeys(keystore)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	keyID := c.Param("id")
	if keyID != "" {
		keyID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		getSingleKey(c, keys, keyID)
	} else {
		c.JSON(http.StatusOK, keys)
	}
}

func getSingleKey(c *gin.Context, keys []key, keyID int) {
	for _, k := range keys {
		if k.ID == keyID {
			c.JSON(http.StatusOK, k)
			return
		}
	}
	c.JSON(http.StatusOK, key{})
}

func createKey(c *gin.Context) {
	keystore := os.Getenv("KEYSTORE")
	existingKeys, err := loadKeys(keystore)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	newKey := key{
		Type: c.PostForm("Type"),
	}
	switch newKey.Type {
	case "ssh":
		newKey, err = createSSHKey(newKey)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	default:
		err = errors.New("invalid key type. Type must be 'ssh'")
		c.AbortWithError(http.StatusBadRequest, err)
	}
	existingKeys = append(existingKeys, newKey)
	err = saveKeys(keystore, existingKeys)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, newKey)
}

func deleteAllKeys(c *gin.Context) {

}

func deleteKey(c *gin.Context) {

}
