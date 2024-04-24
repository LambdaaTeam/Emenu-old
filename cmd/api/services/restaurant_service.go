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

func GetTableById(ctx context.Context, restaurantID string, tableID string) (*models.Table, error) {
	var restaurant models.Restaurant

	objID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	tableObjID, err := primitive.ObjectIDFromHex(tableID)
	if err != nil {
		return nil, fmt.Errorf("invalid table ID")
	}

	err = database.GetCollection("restaurants").FindOne(ctx, bson.M{"_id": objID}).Decode(&restaurant)
	if err != nil {
		return nil, fmt.Errorf("restaurant not found")
	}

	for _, table := range restaurant.Tables {
		if table.ID == tableObjID {
			return &table, nil
		}
	}

	return nil, fmt.Errorf("table not found")
}

func CreateTable(ctx context.Context, id string, params models.Table) (*models.Table, error) {
    var restaurant models.Restaurant

    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, fmt.Errorf("invalid restaurant ID")
    }

    err = database.GetCollection("restaurants").FindOne(ctx, bson.M{"_id": objID}).Decode(&restaurant)
    if err != nil {
        return nil, fmt.Errorf("restaurant not found")
    }

    table := models.Table{
        ID: primitive.NewObjectID(),
        Number: params.Number,
        Url: params.Url,
        Status: params.Status,
        Occupants: params.Occupants,
    }

    restaurant.Tables = append(restaurant.Tables, table)

    _, err = database.GetCollection("restaurants").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"tables": restaurant.Tables}})
    if err != nil {
        return nil, fmt.Errorf("could not update restaurant tables")
    }

    return &table, nil
}


func UpdateTable(ctx context.Context, restaurantID string, tableID string, params models.Table) (*models.Table, error) {
    var restaurant models.Restaurant

    objID, err := primitive.ObjectIDFromHex(restaurantID)
    if err != nil {
        return nil, fmt.Errorf("invalid restaurant ID: %w", err)
    }

    tableObjID, err := primitive.ObjectIDFromHex(tableID)
    if err != nil {
        return nil, fmt.Errorf("invalid table ID: %w", err)
    }

    err = database.GetCollection("restaurants").FindOne(ctx, bson.M{"_id": objID}).Decode(&restaurant)
    if err != nil {
        return nil, fmt.Errorf("restaurant not found: %w", err)
    }

    for i, table := range restaurant.Tables {
        if table.ID == tableObjID {
            restaurant.Tables[i] = params
            restaurant.Tables[i].ID = tableObjID 

            _, err = database.GetCollection("restaurants").UpdateOne(
                ctx,
                bson.M{"_id": objID},
                bson.M{"$set": bson.M{"tables": restaurant.Tables}},
            )
            if err != nil {
                return nil, fmt.Errorf("could not update restaurant tables: %w", err)
            }

            return &restaurant.Tables[i], nil 
        }
    }

    return nil, fmt.Errorf("table not found")
}


func DeleteTable(ctx context.Context, restaurantID string, tableID string) error {
	var restaurant models.Restaurant

	objID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return fmt.Errorf("invalid restaurant ID: %w", err)
	}

	tableObjID, err := primitive.ObjectIDFromHex(tableID)
	if err != nil {
		return fmt.Errorf("invalid table ID: %w", err)
	}

	err = database.GetCollection("restaurants").FindOne(ctx, bson.M{"_id": objID}).Decode(&restaurant)
	if err != nil {
		return fmt.Errorf("restaurant not found: %w", err)
	}

	for i, table := range restaurant.Tables {
		if table.ID == tableObjID {
			restaurant.Tables = append(restaurant.Tables[:i], restaurant.Tables[i+1:]...)

			_, err = database.GetCollection("restaurants").UpdateOne(
				ctx,
				bson.M{"_id": objID},
				bson.M{"$set": bson.M{"tables": restaurant.Tables}},
			)
			if err != nil {
				return fmt.Errorf("could not update restaurant tables: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("table not found")
}