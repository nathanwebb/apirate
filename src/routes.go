package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const apiVersion = "v1"

func initialiseRoutes(router *gin.Engine) {
	basePath := "/api/" + apiVersion
	keysRoutes := router.Group(basePath + "/keys")
	keysRoutes.GET("/", getKeys)
	keysRoutes.GET("/:id", getKeys)
	keysRoutes.POST("/", createKey)
	keysRoutes.DELETE("/", deleteKeys)
	keysRoutes.DELETE("/:id", deleteKeys)
	resultsRoutes := router.Group(basePath + "/results")
	resultsRoutes.GET("/", getResults)
}

func getKeys(c *gin.Context) {
	keys, err := loadKeys()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	keyID := c.Param("id")
	if keyID != "" {
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		getSingleKey(c, keys, keyID)
	} else {
		c.JSON(http.StatusOK, keys)
	}
}

func getSingleKey(c *gin.Context, keys []key, keyID string) {
	for _, k := range keys {
		if k.ID == keyID {
			c.JSON(http.StatusOK, k)
			return
		}
	}
	c.JSON(http.StatusOK, key{})
}

func createKey(c *gin.Context) {
	existingKeys, err := loadKeys()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var newKey key
	if c.ShouldBindJSON(&newKey) == nil {
		log.Println(newKey.Type)
	}
	switch newKey.Type {
	case "ssh", "":
		newKey, err = createSSHKey(newKey)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	default:
		err = errors.New("invalid key type. Type must be 'ssh'")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	log.Println("running this section")
	existingKeys = append(existingKeys, newKey)
	err = saveKeys(existingKeys)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, newKey)
}

func deleteKeys(c *gin.Context) {
	keyID := c.Param("id")
	var err error
	if keyID != "" {
		err = deleteKey(keyID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	} else {
		err = deleteAllKeys()
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
		log.Printf("ERROR: loading commands failed: %s", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	command, err := getCommandForRequest(c, commands)
	if err != nil {
		log.Printf("ERROR: getting commands failed: %s", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	results, err := execCommand(command, queryArgs)
	if err != nil {
		log.Println(err.Error())
		results.Error = err.Error()
	}
	c.JSON(http.StatusOK, results)
}
