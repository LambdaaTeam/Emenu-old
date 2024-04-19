package services

import (
	"context"
	"fmt"

	"github.com/LambdaaTeam/Emenu/pkg/database"
	"github.com/LambdaaTeam/Emenu/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOneRestaurant(ctx context.Context, id string) (*models.PublicRestaurant, error) {
	var restaurant models.Restaurant

	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err

	}

	err = database.GetCollection("restaurants").FindOne(ctx, bson.M{"_id": objID}).Decode(&restaurant)

	if err != nil {
		return nil, err
	}

	return restaurant.ToPublic(), nil
}

func GetAllTables(ctx context.Context, id string) (*[]models.Table, error) {
	var restaurant models.Restaurant

	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, fmt.Errorf("invalid ID")
	}

	err = database.GetCollection("restaurants").FindOne(ctx, bson.M{"_id": objID}).Decode(&restaurant)

	if err != nil {
		return nil, fmt.Errorf("restaurant not found")
	}

	return &restaurant.Tables, nil
}
