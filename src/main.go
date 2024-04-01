package main

import (
    "github.com/gin-gonic/gin"	
)

func main(){
    router := gin.Default()
    router.GET("/autouploader/pixelfed", getUploadedItems)

    router.Run("localhost:8080")
}

func getUploadedItems(c *gin.Context){
	
}