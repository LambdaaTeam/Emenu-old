package middlewares

import (
	"github.com/LambdaaTeam/Emenu/pkg/auth"
	"github.com/LambdaaTeam/Emenu/pkg/database"
	"github.com/LambdaaTeam/Emenu/pkg/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func JWTAuthClient() gin.HandlerFunc {
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
		clientCPF, err := auth.DecodeToken(token)

		if err != nil {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		ctx := c.Request.Context()

		var restaurant models.Restaurant
		err = database.GetCollection("restaurants").FindOne(ctx, bson.M{"tables.occupants.cpf": clientCPF}).Decode(&restaurant)
		if err != nil {
			c.JSON(404, gin.H{
				"error": "Client not found in any restaurant",
			})
			c.Abort()
			return
		}

		var client models.Client
		for _, table := range restaurant.Tables {
			for _, occupant := range table.Occupants {
				if occupant.CPF == clientCPF {
					client = occupant
					break
				}
			}
		}

		if client.CPF == "" {
			c.JSON(404, gin.H{
				"error": "Client not found in any restaurant",
			})
			c.Abort()
			return
		}

		c.Set("client", client)

		c.Next()
	}
}
