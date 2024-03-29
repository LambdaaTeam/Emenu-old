package services

import (
	"context"
	"errors"

	pkg "github.com/LambdaaTeam/Emenu/pkg/auth"
	"github.com/LambdaaTeam/Emenu/pkg/database"
	"github.com/LambdaaTeam/Emenu/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(payload models.RestaurantRegiter) (*models.PublicRestaurant, error) {
	restaurant := payload.ToRestaurant()

	_, err := database.GetCollection("restaurants").InsertOne(context.TODO(), restaurant)

	if err != nil {
		return nil, err
	}

	publicRestaurant := restaurant.ToPublic()

	return publicRestaurant, nil
}

func Login(payload models.RestaurantLogin) (*models.PublicRestaurant, error) {
	var restaurant models.Restaurant
	err := database.GetCollection("restaurants").FindOne(context.TODO(), bson.M{"email": payload.Email}).Decode(&restaurant)

	if err != nil {
		return nil, err
	}

	if !pkg.IsPasswordValid(restaurant.Password, payload.Password) {
		return nil, errors.New("invalid password")
	}

	// token := pkg.GenerateToken(user.ID.Hex())
	publicUserWithToken := restaurant.ToPublic()

	return publicUserWithToken, nil
}
