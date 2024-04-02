package main

import (
	"net/http"
	"os"

	"github.com/LNA-DEV/Personal.API/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var logger = logrus.New()
var mongoConnectionString = ""
var apiKey = ""

func main() {
	mongoConnectionString = os.Getenv("MONGODB")
	apiKey = os.Getenv("API_KEY")

	if apiKey == "" {
		logger.Warn("ApiKey not set!!!")
		apiKey = "1234"
	}

	router := gin.Default()

	// Endpoints
	router.GET("/autouploader/pixelfed", getUploadedItemsRoot)
	router.POST("/autouploader/pixelfed", validateAPIKey(), addUploadedItem)

	router.Run("0.0.0.0:8080")
}

func getUploadedItemsRoot(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getUploadedItems().Value)
}

func getUploadedItems() AlreadyPublishedItems {
	alreadyUploaded, err := repository.ReadMongo[AlreadyPublishedItems]("Autouploader", "AlreadyUploaded", bson.D{{"key", "Pixelfed"}}, mongoConnectionString)

	if err != nil {
		logger.Error(err)

		return AlreadyPublishedItems{
			Key:        "Pixelfed",
			Value:      []string{},
			NotCreated: true,
		}
	}

	return alreadyUploaded
}

func addUploadedItem(c *gin.Context) {
	item := c.Query("item")

	items := getUploadedItems()

	items.Value = append(items.Value, item)

	var err error

	if items.NotCreated {
		items.NotCreated = false
		err = repository.WriteMongo("Autouploader", "AlreadyUploaded", items, mongoConnectionString)
	} else {
		err = repository.UpdateMongo("Autouploader", "AlreadyUploaded", bson.D{{"$set", items}}, bson.D{{"key", "Pixelfed"}}, mongoConnectionString)
	}

	if err != nil {
		logger.Error(err)

		c.Status(http.StatusInternalServerError)

		return
	}
}

func validateAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		authentication := c.Request.Header.Get("Authorization")
		expectedAuth := "ApiKey " + apiKey

		if authentication != expectedAuth {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Authentication failed"})
			c.Abort()
			return
		}

		c.Next()
	}
}
