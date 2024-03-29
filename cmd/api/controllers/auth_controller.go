package controllers

import (
	"net/http"

	"github.com/LambdaaTeam/Emenu/cmd/api/services"
	"github.com/LambdaaTeam/Emenu/pkg/models"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var payload models.RestaurantRegiter

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurant, err := services.Register(payload)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

func Login(c *gin.Context) {
	var payload models.RestaurantLogin

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurant, err := services.Login(payload)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}
