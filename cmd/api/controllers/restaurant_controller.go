package controllers

import (
	"net/http"

	"github.com/LambdaaTeam/Emenu/cmd/api/services"
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
    restaurantID := c.Param("id")

    var tablePayload models.Table
    if err := c.ShouldBindJSON(&tablePayload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input, please check your data"})
        return
    }

    ctx := c.Request.Context()
    table, err := services.CreateTable(ctx, restaurantID, tablePayload)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"failed to create a new table": err.Error()})
        return
    }

    c.JSON(http.StatusOK, table)
}

func UpdateTable(c *gin.Context) {
	restaurantID := c.Param("id")
	tableID := c.Param("tableId")

	var tablePayload models.Table
	if err := c.ShouldBindJSON(&tablePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input, please check your data"})
		return
	}

	ctx := c.Request.Context()
	updatedTable, err := services.UpdateTable(ctx, restaurantID, tableID, tablePayload) 
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to update table": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTable) 
}


func DeleteTable(c *gin.Context) {
	restaurantID := c.Param("id")
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
	restaurantID := c.Param("id")
	tableID := c.Param("tableId")

	ctx := c.Request.Context()

	table, err := services.GetTableById(ctx, restaurantID, tableID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to get table": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

func GetAllTables(c *gin.Context) {
	restaurantID := c.Param("id")

	ctx := c.Request.Context()

	tables, err := services.GetAllTables(ctx, restaurantID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to get all tables": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tables)
}


