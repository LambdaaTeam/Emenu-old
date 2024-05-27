package controllers

import (
	"net/http"
	"strconv"

	"github.com/LambdaaTeam/Emenu/cmd/api/services"
	"github.com/LambdaaTeam/Emenu/pkg/auth"
	"github.com/LambdaaTeam/Emenu/pkg/models"
	"github.com/gin-gonic/gin"
)

func GetOneRestaurant(c *gin.Context) {
	restaurantID := c.Param("id")

	ctx := c.Request.Context()

	restaurant, err := services.GetOneRestaurant(ctx, restaurantID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

func CreateTable(c *gin.Context) {
	restaurantID := c.MustGet("restaurant").(string)

	var tablePayload struct {
		Number int `json:"number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&tablePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input, please check your data"})
		return
	}

	ctx := c.Request.Context()
	table, err := services.CreateTable(ctx, restaurantID, tablePayload.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

func UpdateTable(c *gin.Context) {
	restaurantID := c.MustGet("restaurant").(string)
	tableID := c.Param("tableId")

	var tablePayload models.Table
	if err := c.ShouldBindJSON(&tablePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input, please check your data"})
		return
	}

	ctx := c.Request.Context()
	updatedTable, err := services.UpdateTable(ctx, restaurantID, tableID, tablePayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTable)
}

func DeleteTable(c *gin.Context) {
	restaurantID := c.MustGet("restaurant").(string)
	tableID := c.Param("tableId")

	ctx := c.Request.Context()
	err := services.DeleteTable(ctx, restaurantID, tableID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to delete table": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "table deleted successfully"})
}

func GetTableById(c *gin.Context) {
	restaurantID := c.MustGet("restaurant").(string)
	tableID := c.Param("tableId")

	ctx := c.Request.Context()

	table, err := services.GetTableById(ctx, restaurantID, tableID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

func GetAllTables(c *gin.Context) {
	restaurantID := c.MustGet("restaurant").(string)
	ctx := c.Request.Context()

	tables, err := services.GetAllTables(ctx, restaurantID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tables)
}

func GetOrders(c *gin.Context) {
	restaurantID := c.MustGet("restaurant").(string)
	page := c.Query("page")

	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.ParseInt(page, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return
	}

	ctx := c.Request.Context()

	orders, err := services.GetOrders(ctx, restaurantID, pageInt)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func AddOrderItem(c *gin.Context) {
	restaurantID := c.Param("id")
	orderID := c.Param("orderId")

	var itemPayload models.OrderItem

	if err := c.ShouldBindJSON(&itemPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input, please check your data"})
		return
	}

	ctx := c.Request.Context()
	order, err := services.AddOrderItem(ctx, restaurantID, orderID, itemPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func UpdateOrderItem(c *gin.Context) {
	restaurantID := c.MustGet("restaurant").(string)
	orderID := c.Param("orderId")

	var itemPayload models.OrderItem

	if err := c.ShouldBindJSON(&itemPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input, please check your data"})
		return
	}

	ctx := c.Request.Context()
	order, err := services.UpdateOrderItem(ctx, restaurantID, orderID, itemPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func GetOrderByID(c *gin.Context) {
	restaurantID := c.MustGet("restaurant").(string)
	orderID := c.Param("orderId")

	ctx := c.Request.Context()

	order, err := services.GetOrderByID(ctx, restaurantID, orderID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func GetMenu(c *gin.Context) {
	restaurantID := c.Param("id")

	ctx := c.Request.Context()

	menu, err := services.GetMenu(ctx, restaurantID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

func AddCategoryToMenu(c *gin.Context) {
	restaurantID := c.MustGet("restaurant").(string)

	var categoryPayload struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&categoryPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input, please check your data"})
		return
	}

	ctx := c.Request.Context()
	menu, err := services.AddCategoryToMenu(ctx, restaurantID, categoryPayload.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

func UpdateCategory(c *gin.Context) {
	restaurantID := c.MustGet("restaurant").(string)
	categoryID := c.Param("categoryId")

	var categoryPayload struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&categoryPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input, please check your data"})
		return
	}

	ctx := c.Request.Context()
	updatedCategory, err := services.UpdateCategoryInMenu(ctx, restaurantID, categoryID, categoryPayload.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCategory)
}

func DeleteCategory(c *gin.Context) {
	restaurantID := c.MustGet("restaurant").(string)
	categoryID := c.Param("categoryId")

	ctx := c.Request.Context()
	menu, err := services.DeleteCategoryFromMenu(ctx, restaurantID, categoryID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

func AddItemToMenu(c *gin.Context) {
	restaurantID := c.MustGet("restaurant").(string)
	categoryID := c.Param("categoryId")

	var itemPayload models.Item

	if err := c.ShouldBindJSON(&itemPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input, please check your data"})
		return
	}

	ctx := c.Request.Context()
	menu, err := services.AddItemToMenu(ctx, restaurantID, categoryID, itemPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

func UpdateItem(c *gin.Context) {
	restaurantID := c.MustGet("restaurant").(string)
	categoryID := c.Param("categoryId")
	itemID := c.Param("itemId")

	var itemPayload models.Item

	if err := c.ShouldBindJSON(&itemPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input, please check your data"})
		return
	}

	ctx := c.Request.Context()
	updatedItem, err := services.UpdateItemInMenu(ctx, restaurantID, categoryID, itemID, itemPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedItem)
}

func DeleteItem(c *gin.Context) {
	restaurantID := c.MustGet("restaurant").(string)
	categoryID := c.Param("categoryId")
	itemID := c.Param("itemId")

	ctx := c.Request.Context()
	menu, err := services.DeleteItemFromMenu(ctx, restaurantID, categoryID, itemID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

func AddClientToTable(c *gin.Context) {
	restaurantID := c.Param("id")
	tableID := c.Param("tableId")

	var clientPayload models.Client
	if err := c.ShouldBindJSON(&clientPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input, please check your data"})
		return
	}

	clientToken, err := auth.GenerateClientToken(clientPayload.CPF)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	table, err := services.AddClientToTable(ctx, restaurantID, tableID, clientPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table.AddToken(clientToken))
}
