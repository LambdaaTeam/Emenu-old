package main

import (
	"net/http"
	"os"

	"github.com/LambdaaTeam/Emenu/cmd/ws/handlers"
	"github.com/LambdaaTeam/Emenu/pkg/database"
	"github.com/gin-contrib/cors"
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

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.POST("/notify", handlers.Notify)

	r.GET("/ws", handlers.UpgradeConnection)

	r.Run()
}
