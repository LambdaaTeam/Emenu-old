package main

import (
	"os"

	"github.com/LambdaaTeam/Emenu/cmd/api/controllers"
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

	r.GET("/", controllers.HealthCheck)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/login", controllers.Login)
		v1.POST("/register", controllers.Register)
	}

	r.Run()
}
