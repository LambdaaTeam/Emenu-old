package main

import (
	"github.com/LambdaaTeam/Emenu/cmd/api/controllers"
	"github.com/LambdaaTeam/Emenu/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../../.env")

	database.DB = database.Connect()
}

func main() {
	r := gin.Default()

	r.GET("/", controllers.HealthCheck)

	r.Run()
}
