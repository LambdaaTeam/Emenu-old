package middlewares

import (
	"github.com/LambdaaTeam/Emenu/pkg/auth"
	"github.com/LambdaaTeam/Emenu/pkg/database"
	"github.com/LambdaaTeam/Emenu/pkg/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func JWTAuthRestaurant() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		if token == "" {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		token = token[7:]
		restaurantID, err := auth.DecodeToken(token)

		if err != nil {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		restaurantOID, err := primitive.ObjectIDFromHex(restaurantID)
		if err != nil {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		var restaurant models.Restaurant
		ctx := c.Request.Context()
		err = database.GetCollection("restaurants").FindOne(ctx, bson.M{"_id": restaurantOID}).Decode(&restaurant)

		if err != nil {
			c.JSON(404, gin.H{
				"error": "Restaurant not found",
			})
			c.Abort()
			return
		}

		c.Set("restaurant", restaurant.ID.Hex())

		c.Next()
	}
}
