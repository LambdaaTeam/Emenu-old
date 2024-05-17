package controllers

import (
	"net/http"

	"github.com/LambdaaTeam/Emenu/cmd/api/services"
	"github.com/LambdaaTeam/Emenu/pkg/auth"
	"github.com/LambdaaTeam/Emenu/pkg/models"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var payload models.RestaurantRegister

	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.Email == "" || payload.Name == "" || payload.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email, name and password are required"})
		return
	}

	restaurant, err := services.Register(ctx, payload)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

func Login(c *gin.Context) {
	var payload models.RestaurantLogin

	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurant, err := services.Login(ctx, payload)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurantToken, err := auth.GenerateRestaurantToken(restaurant.ID)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, restaurantToken)
}
