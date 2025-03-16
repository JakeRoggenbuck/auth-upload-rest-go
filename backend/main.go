package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"log"
	"net/http"
	"os"
)

func getLogIn() gin.Accounts {
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		color.Error.Println("ADMIN_PASSWORD not set")
		log.Fatal("ADMIN_PASSWORD not set")
	}

	return gin.Accounts{
		"Admin": password,
	}
}

func setupLogging() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
}

func homePage(c *gin.Context) {
	c.String(http.StatusOK, "Home\n")
}

func uploadFile(c *gin.Context) {
	file, err := c.FormFile("myFile")
	if err != nil {
		fmt.Println(err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	log.Println(file.Filename)

	c.SaveUploadedFile(file, "out/uploaded-"+file.Filename)
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
	})
}

func listFiles(c *gin.Context) {
	var files []string

	f, err := os.Open("out/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Path 'out/' not created.",
		})
	}

	fileInfo, err := f.Readdir(-1)
	f.Close()

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}

	c.JSON(http.StatusOK, files)
}

func main() {
	setupLogging()
	router := gin.Default()

	routeStart := "/api/v1/"
	authAccount := getLogIn()

	authedSubRoute := router.Group(routeStart, gin.BasicAuth(authAccount))

	authedSubRoute.GET("/", homePage)
	authedSubRoute.POST("/upload", uploadFile)
	authedSubRoute.GET("/list", listFiles)

	listenPort := os.Getenv("PORT")
	if listenPort == "" {
		listenPort = "1357"
	}

	fmt.Print("\nHosted at ")
	color.Magenta.Println("http://localhost:" + listenPort + routeStart + "\n")

	router.Run(":" + listenPort)
}
