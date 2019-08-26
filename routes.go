package main

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

const apiVersion = "v1"

func initialiseRoutes(router *gin.Engine) {
	basePath := "/api/" + apiVersion
	keysRoutes := router.Group("/keys")
	keysRoutes.GET("/", getKeys)
	keysRoutes.GET("/:id", getKeys)
	keysRoutes.POST("/", createKey)
	keysRoutes.DELETE("/", deleteKeys)
	keysRoutes.DELETE("/:id", deleteKeys)
	router.GET(basePath+"/results", getResults)
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

func deleteKeys(c *gin.Context) {
	keystore := os.Getenv("KEYSTORE")
	keyIDStr := c.Param("id")
	var err error
	if keyIDStr != "" {
		keyID, err := strconv.Atoi(keyIDStr)
		err = deleteKey(keystore, keyID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	} else {
		err = deleteAllKeys(keystore)
	}
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func getResults(c *gin.Context) {
	queryArgs := c.Request.URL.Query()
	commandstore := os.Getenv("COMMANDSTORE")
	commands, err := loadCommands(commandstore)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	command, err := getCommandForRequest(c, commands)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	results, err := execCommand(command, queryArgs)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, results)
}
