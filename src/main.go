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

func main() {
	mongoConnectionString = os.Getenv("MONGODB")

	router := gin.Default()
	router.GET("/autouploader/pixelfed", getUploadedItemsRoot)
	router.POST("/autouploader/pixelfed", addUploadedItem)

	router.Run("0.0.0.0:8080")
}

func getUploadedItemsRoot(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getUploadedItems().Value)
}

func getUploadedItems() *AlreadyPublishedItems {
	alreadyUploaded, err := repository.ReadMongo[AlreadyPublishedItems]("Autouploader", "AlreadyUploaded", bson.D{{"key", "Pixelfed"}}, mongoConnectionString)

	if err != nil {
		logger.Error(err)

		return nil
	}

	return &alreadyUploaded
}

func addUploadedItem(c *gin.Context) {
	item := c.Query("item")
	create := false

	items := getUploadedItems()

	if items == nil {
		items = &AlreadyPublishedItems{
			Key: "Pixelfed",
		}
		items.Value = []string{}

		create = true
	}

	items.Value = append(items.Value, item)

	var err error

	if create {
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
