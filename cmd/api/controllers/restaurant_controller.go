package controllers

import (
	"net/http"

	"github.com/LambdaaTeam/Emenu/cmd/api/services"
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


func GetAllTables(c *gin.Context) {
	restaurantID := c.Param("id")

	ctx := c.Request.Context()

	tables, err := services.GetAllTables(ctx, restaurantID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tables)
}