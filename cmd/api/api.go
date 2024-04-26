package main

import (
	"os"

	"github.com/LambdaaTeam/Emenu/cmd/api/controllers"
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

	r.GET("/", controllers.HealthCheck)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/login", controllers.Login)
		v1.POST("/register", controllers.Register)
		v1.GET("/restaurants/:id", controllers.GetOneRestaurant)
		v1.GET("/restaurants/:id/tables", controllers.GetAllTables)
		v1.GET("/restaurants/:id/tables/:tableId", controllers.GetTableById)
		v1.POST("/restaurants/:id/tables", controllers.CreateTable)
		v1.PATCH("/restaurants/:id/tables/:tableId", controllers.UpdateTable)
		v1.DELETE("/restaurants/:id/tables/:tableId", controllers.DeleteTable)
		v1.GET("/restaurants/:id/orders", controllers.GetOrders)
		v1.GET("/restaurants/:id/orders/:orderId", controllers.GetOrderByID)
	}

	r.Run()
}
