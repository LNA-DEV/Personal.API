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

	// Endpoints for Pixelfed
	router.GET("/autouploader/pixelfed", func(c *gin.Context) {
		getUploadedItems(c, "Pixelfed")
	})
	router.POST("/autouploader/pixelfed", validateAPIKey(), func(c *gin.Context) {
		addUploadedItem(c, "Pixelfed")
	})

	// Endpoints for Bluesky
	router.GET("/autouploader/bluesky", func(c *gin.Context) {
		getUploadedItems(c, "Bluesky")
	})
	router.POST("/autouploader/bluesky", validateAPIKey(), func(c *gin.Context) {
		addUploadedItem(c, "Bluesky")
	})

	router.Run("0.0.0.0:8080")
}

func getUploadedItems(c *gin.Context, platform string) {
	items := fetchUploadedItems(platform)
	c.IndentedJSON(http.StatusOK, items.Value)
}

func fetchUploadedItems(platform string) AlreadyPublishedItems {
	alreadyUploaded, err := repository.ReadMongo[AlreadyPublishedItems](
		"Autouploader",
		"AlreadyUploaded",
		bson.D{{"key", platform}},
		mongoConnectionString,
	)

	if err != nil {
		logger.Error(err)
		return AlreadyPublishedItems{
			Key:        platform,
			Value:      []string{},
			NotCreated: true,
		}
	}

	return alreadyUploaded
}

func addUploadedItem(c *gin.Context, platform string) {
	item := c.Query("item")
	if item == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Item is required"})
		return
	}

	items := fetchUploadedItems(platform)
	items.Value = append(items.Value, item)

	var err error
	if items.NotCreated {
		items.NotCreated = false
		err = repository.WriteMongo("Autouploader", "AlreadyUploaded", items, mongoConnectionString)
	} else {
		err = repository.UpdateMongo(
			"Autouploader",
			"AlreadyUploaded",
			bson.D{{"$set", bson.D{{"value", items.Value}}}},
			bson.D{{"key", platform}},
			mongoConnectionString,
		)
	}

	if err != nil {
		logger.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
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
