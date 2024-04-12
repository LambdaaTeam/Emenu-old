package services

import (
	"context"
	"errors"

	pkg "github.com/LambdaaTeam/Emenu/pkg/auth"
	"github.com/LambdaaTeam/Emenu/pkg/database"
	"github.com/LambdaaTeam/Emenu/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(ctx context.Context, payload models.RestaurantRegister) (*models.PublicRestaurant, error) {
	restaurant := payload.ToRestaurant()

	_, err := database.GetCollection("restaurants").InsertOne(ctx, restaurant)

	if err != nil {
		return nil, err
	}

	publicRestaurant := restaurant.ToPublic()

	return publicRestaurant, nil
}

func Login(ctx context.Context, payload models.RestaurantLogin) (*models.PublicRestaurant, error) {
	var restaurant models.Restaurant
	err := database.GetCollection("restaurants").FindOne(ctx, bson.M{"email": payload.Email}).Decode(&restaurant)

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
