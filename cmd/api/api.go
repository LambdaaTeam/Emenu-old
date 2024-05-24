package main

import (
	"os"

	"github.com/LambdaaTeam/Emenu/cmd/api/controllers"
	"github.com/LambdaaTeam/Emenu/cmd/api/middlewares"
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
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE"}
	r.Use(cors.New(config))

	r.GET("/", controllers.HealthCheck)

	v1 := r.Group("/api/v1")
	{
		// Auth
		v1.POST("/login", controllers.Login)
		v1.POST("/register", controllers.Register)

		// Restaurants
		v1.GET("/restaurants/:id", controllers.GetOneRestaurant)

		// Tables
		v1.GET("/restaurants/@me/tables", middlewares.JWTAuthRestaurant(), controllers.GetAllTables)
		v1.POST("/restaurants/@me/tables", middlewares.JWTAuthRestaurant(), controllers.CreateTable)
		v1.GET("/restaurants/@me/tables/:tableId", middlewares.JWTAuthRestaurant(), controllers.GetTableById)
		v1.POST("/restaurants/:id/tables/:tableId", controllers.AddClientToTable) // Add client to table (no auth)
		v1.PATCH("/restaurants/@me/tables/:tableId", middlewares.JWTAuthRestaurant(), controllers.UpdateTable)
		v1.DELETE("/restaurants/@me/tables/:tableId", middlewares.JWTAuthRestaurant(), controllers.DeleteTable)

		// Orders
		v1.GET("/restaurants/@me/orders", middlewares.JWTAuthRestaurant(), controllers.GetOrders)
		v1.GET("/restaurants/@me/orders/:orderId", middlewares.JWTAuthRestaurant(), controllers.GetOrderByID)
		v1.POST("/restaurants/:id/orders/:orderId", middlewares.JWTAuthClient(), controllers.AddOrderItem)
		v1.PATCH("/restaurants/@me/orders/:orderId", middlewares.JWTAuthRestaurant(), controllers.UpdateOrderItem)

		// Menu
		v1.GET("/restaurants/:id/menu", controllers.GetMenu)

		// Menu Categories
		v1.POST("/restaurants/@me/menu/categories", middlewares.JWTAuthRestaurant(), controllers.AddCategoryToMenu)
		v1.PATCH("/restaurants/@me/menu/categories/:categoryId", middlewares.JWTAuthRestaurant(), controllers.UpdateCategory)
		v1.DELETE("/restaurants/@me/menu/categories/:categoryId", middlewares.JWTAuthRestaurant(), controllers.DeleteCategory)

		// Menu Subcategories
		v1.POST("/restaurants/@me/menu/categories/:categoryId/subcategories", middlewares.JWTAuthRestaurant(), controllers.AddSubcategoryToCategory)
		v1.PATCH("/restaurants/@me/menu/categories/:categoryId/subcategories/:subcategoryId", middlewares.JWTAuthRestaurant(), controllers.UpdateSubcategory)
		v1.DELETE("/restaurants/@me/menu/categories/:categoryId/subcategories/:subcategoryId", middlewares.JWTAuthRestaurant(), controllers.DeleteSubcategory)

		// Menu Items
		v1.POST("/restaurants/@me/menu/categories/:categoryId/subcategories/:subcategoryId/items", middlewares.JWTAuthRestaurant(), controllers.AddItemToMenu)
		v1.PATCH("/restaurants/@me/menu/categories/:categoryId/subcategories/:subcategoryId/items/:itemId", middlewares.JWTAuthRestaurant(), controllers.UpdateItem)
		v1.DELETE("/restaurants/@me/menu/categories/:categoryId/subcategories/:subcategoryId/items/:itemId", middlewares.JWTAuthRestaurant(), controllers.DeleteItem)
	}

	r.Run()
}
