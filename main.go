package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
)

func homePage(c *gin.Context) {
	c.String(http.StatusOK, "Home\n")
}

func uploadFile(c *gin.Context) {
	file, err := c.FormFile("myFile")
	if err != nil {
		fmt.Println(err)
	}

	log.Println(file.Filename)

	c.SaveUploadedFile(file, "out/uploaded-" + file.Filename)
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func main() {
	router := gin.Default()

	authedSubRoute := router.Group("/api/v1/", gin.BasicAuth(gin.Accounts{
		"admin": "adminpass",
	}))

	authedSubRoute.GET("/", homePage)
	authedSubRoute.POST("/upload", uploadFile)
	
	listenPort := "1357"
	router.Run(":" + listenPort)
}
