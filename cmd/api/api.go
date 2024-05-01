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
		// Login and Register
		v1.POST("/login", controllers.Login)
		v1.POST("/register", controllers.Register)

		// Restaurants
		v1.GET("/restaurants/:id", controllers.GetOneRestaurant)

		// Tables
		v1.GET("/restaurants/:id/tables", controllers.GetAllTables)
		v1.POST("/restaurants/:id/tables", controllers.CreateTable)
		v1.GET("/restaurants/:id/tables/:tableId", controllers.GetTableById)
		v1.POST("/restaurants/:id/tables/:tableId", controllers.AddClientToTable)
		v1.PATCH("/restaurants/:id/tables/:tableId", controllers.UpdateTable)
		v1.DELETE("/restaurants/:id/tables/:tableId", controllers.DeleteTable)

		// Orders
		v1.GET("/restaurants/:id/orders", controllers.GetOrders)
		v1.GET("/restaurants/:id/orders/:orderId", controllers.GetOrderByID)
		v1.POST("/restaurants/:id/orders/:orderId", controllers.AddOrderItem)
		v1.PATCH("/restaurants/:id/orders/:orderId", controllers.UpdateOrderItem)

		// Menu
		v1.GET("/restaurants/:id/menu", controllers.GetMenu)

		// Menu Categories
		v1.POST("/restaurants/:id/menu/categories", controllers.AddCategoryToMenu)
		v1.PATCH("/restaurants/:id/menu/categories/:categoryId", controllers.UpdateCategory)
		v1.DELETE("/restaurants/:id/menu/categories/:categoryId", controllers.DeleteCategory)

		// Menu Subcategories
		v1.POST("/restaurants/:id/menu/categories/:categoryId/subcategories", controllers.AddSubcategoryToCategory)
		v1.PATCH("/restaurants/:id/menu/categories/:categoryId/subcategories/:subcategoryId", controllers.UpdateSubcategory)
		v1.DELETE("/restaurants/:id/menu/categories/:categoryId/subcategories/:subcategoryId", controllers.DeleteSubcategory)

		// Menu Items
		v1.POST("/restaurants/:id/menu/categories/:categoryId/subcategories/:subcategoryId/items", controllers.AddItemToMenu)
		v1.PATCH("/restaurants/:id/menu/categories/:categoryId/subcategories/:subcategoryId/items/:itemId", controllers.UpdateItem)
		v1.DELETE("/restaurants/:id/menu/categories/:categoryId/subcategories/:subcategoryId/items/:itemId", controllers.DeleteItem)
	}

	r.Run()
}
