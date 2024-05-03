package main

import (
	"net/http"
	"os"

	"github.com/LambdaaTeam/Emenu/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env")

	dbName := "emenu-dev"

	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
		dbName = "emenu"
	}

	database.DB = database.Connect(dbName)
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.POST("/push", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.GET("/ws", func(c *gin.Context) {
		// TODO: upgrade to websocket
		// TODO: handle websocket
		// TODO: create packet models
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	r.Run()
}
