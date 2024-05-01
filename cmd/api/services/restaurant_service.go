package services

import (
	"context"
	"fmt"

	"github.com/LambdaaTeam/Emenu/pkg/database"
	"github.com/LambdaaTeam/Emenu/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func CreateTable(ctx context.Context, id string, number int) (*models.Table, error) {
	var restaurant models.Restaurant

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	err = database.GetCollection("restaurants").FindOne(ctx, bson.M{"_id": objID}).Decode(&restaurant)
	if err != nil {
		return nil, fmt.Errorf("restaurant not found")
	}

	// TODO: generate a unique URL
	table := models.Table{
		ID:        primitive.NewObjectID(),
		Number:    number,
		Url:       "",
		Status:    models.TableStatusAvailable,
		Occupants: []models.Client{},
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

			if restaurant.Tables[i].Occupants == nil {
				restaurant.Tables[i].Occupants = []models.Client{}
			}

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

func GetOrderByID(ctx context.Context, restaurantID string, orderID string) (*models.PublicOrder, error) {
	var order models.Order

	orderObjID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID")
	}

	err = database.GetCollection("orders").FindOne(ctx, bson.M{"_id": orderObjID}).Decode(&order)

	if err != nil {
		return nil, fmt.Errorf("order not found")
	}

	return order.ToPublic(), nil
}

func GetOrders(ctx context.Context, restaurantID string, page int64) (*[]models.PublicOrder, error) {
	var orders []models.Order

	objID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	cursor, err := database.GetCollection("orders").Find(
		ctx,
		bson.M{"restaurant": objID},
		options.
			Find().
			SetLimit(100).
			SetSkip((page-1)*100).
			SetSort(bson.M{"created_at": -1}),
	)

	if err != nil {
		return nil, fmt.Errorf("could not get orders")
	}

	err = cursor.All(ctx, &orders)
	if err != nil {
		return nil, fmt.Errorf("could not get orders")
	}

	publicOrders := make([]models.PublicOrder, len(orders))
	for i, order := range orders {
		publicOrders[i] = *order.ToPublic()
	}

	return &publicOrders, nil
}

func GetMenu(ctx context.Context, restaurantID string) (*models.PublicMenu, error) {
	var menu models.Menu

	objID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	err = database.GetCollection("menus").FindOne(ctx, bson.M{"restaurant": objID}).Decode(&menu)
	if err != nil {
		return nil, fmt.Errorf("menu not found")
	}

	return menu.ToPublic(), nil
}

func AddCategoryToMenu(ctx context.Context, restaurantID string, name string) (*models.PublicMenu, error) {
	var menu models.Menu

	objID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	err = database.GetCollection("menus").FindOne(ctx, bson.M{"restaurant": objID}).Decode(&menu)
	if err != nil {
		return nil, fmt.Errorf("menu not found")
	}

	menu.Categories = append(menu.Categories, models.Category{
		ID:   primitive.NewObjectID(),
		Name: name,
		Sub:  []models.Subcategory{},
	})

	_, err = database.GetCollection("menus").UpdateOne(ctx, bson.M{"restaurant": objID}, bson.M{"$set": bson.M{"categories": menu.Categories}})
	if err != nil {
		return nil, fmt.Errorf("could not update menu")
	}

	return menu.ToPublic(), nil
}

func UpdateCategoryInMenu(ctx context.Context, restaurantID string, categoryID string, name string) (*models.PublicMenu, error) {
	var menu models.Menu

	objID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	catID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID")
	}

	err = database.GetCollection("menus").FindOne(ctx, bson.M{"restaurant": objID}).Decode(&menu)
	if err != nil {
		return nil, fmt.Errorf("menu not found")
	}

	for i, cat := range menu.Categories {
		if cat.ID == catID {
			menu.Categories[i].Name = name

			_, err = database.GetCollection("menus").UpdateOne(ctx, bson.M{"restaurant": objID}, bson.M{"$set": bson.M{"categories": menu.Categories}})
			if err != nil {
				return nil, fmt.Errorf("could not update menu")
			}

			return menu.ToPublic(), nil
		}
	}

	return nil, fmt.Errorf("category not found")
}

func DeleteCategoryFromMenu(ctx context.Context, restaurantID string, categoryID string) (*models.PublicMenu, error) {
	var menu models.Menu

	objID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	catID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID")
	}

	err = database.GetCollection("menus").FindOne(ctx, bson.M{"restaurant": objID}).Decode(&menu)
	if err != nil {
		return nil, fmt.Errorf("menu not found")
	}

	for i, cat := range menu.Categories {
		if cat.ID == catID {
			menu.Categories = append(menu.Categories[:i], menu.Categories[i+1:]...)

			_, err = database.GetCollection("menus").UpdateOne(ctx, bson.M{"restaurant": objID}, bson.M{"$set": bson.M{"categories": menu.Categories}})
			if err != nil {
				return nil, fmt.Errorf("could not update menu")
			}

			return menu.ToPublic(), nil
		}
	}

	return nil, fmt.Errorf("category not found")
}

func AddSubcategoryToMenu(ctx context.Context, restaurantID string, categoryID string, name string) (*models.PublicMenu, error) {
	var menu models.Menu

	objID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	catID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID")
	}

	err = database.GetCollection("menus").FindOne(ctx, bson.M{"restaurant": objID}).Decode(&menu)
	if err != nil {
		return nil, fmt.Errorf("menu not found")
	}

	for i, category := range menu.Categories {
		if category.ID == catID {
			menu.Categories[i].Sub = append(menu.Categories[i].Sub, models.Subcategory{
				ID:    primitive.NewObjectID(),
				Name:  name,
				Items: []models.Item{},
			})

			_, err = database.GetCollection("menus").UpdateOne(ctx, bson.M{"restaurant": objID}, bson.M{"$set": bson.M{"categories": menu.Categories}})
			if err != nil {
				return nil, fmt.Errorf("could not update menu")
			}

			return menu.ToPublic(), nil
		}
	}

	return nil, fmt.Errorf("category not found")
}

func UpdateSubcategoryInMenu(ctx context.Context, restaurantID string, categoryID string, subcategoryID string, name string) (*models.PublicMenu, error) {
	var menu models.Menu

	objID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	catID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID")
	}

	subcatID, err := primitive.ObjectIDFromHex(subcategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid subcategory ID")
	}

	err = database.GetCollection("menus").FindOne(ctx, bson.M{"restaurant": objID}).Decode(&menu)
	if err != nil {
		return nil, fmt.Errorf("menu not found")
	}

	for i, category := range menu.Categories {
		if category.ID == catID {
			for j, subcategory := range category.Sub {
				if subcategory.ID == subcatID {
					menu.Categories[i].Sub[j].Name = name

					_, err = database.GetCollection("menus").UpdateOne(ctx, bson.M{"restaurant": objID}, bson.M{"$set": bson.M{"categories": menu.Categories}})
					if err != nil {
						return nil, fmt.Errorf("could not update menu")
					}

					return menu.ToPublic(), nil
				}
			}
		}
	}

	return nil, fmt.Errorf("subcategory not found")
}

func DeleteSubcategoryFromMenu(ctx context.Context, restaurantID string, categoryID string, subcategoryID string) (*models.PublicMenu, error) {
	var menu models.Menu

	objID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	catID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID")
	}

	subcatID, err := primitive.ObjectIDFromHex(subcategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid subcategory ID")
	}

	err = database.GetCollection("menus").FindOne(ctx, bson.M{"restaurant": objID}).Decode(&menu)
	if err != nil {
		return nil, fmt.Errorf("menu not found")
	}

	for i, category := range menu.Categories {
		if category.ID == catID {
			for j, subcategory := range category.Sub {
				if subcategory.ID == subcatID {
					menu.Categories[i].Sub = append(menu.Categories[i].Sub[:j], menu.Categories[i].Sub[j+1:]...)

					_, err = database.GetCollection("menus").UpdateOne(ctx, bson.M{"restaurant": objID}, bson.M{"$set": bson.M{"categories": menu.Categories}})
					if err != nil {
						return nil, fmt.Errorf("could not update menu")
					}

					return menu.ToPublic(), nil
				}
			}
		}
	}

	return nil, fmt.Errorf("subcategory not found")
}

func AddItemToMenu(ctx context.Context, restaurantID string, categoryID string, subcategoryID string, item models.Item) (*models.PublicMenu, error) {
	var menu models.Menu

	objID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	catID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID")
	}

	subcatID, err := primitive.ObjectIDFromHex(subcategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid subcategory ID")
	}

	err = database.GetCollection("menus").FindOne(ctx, bson.M{"restaurant": objID}).Decode(&menu)
	if err != nil {
		return nil, fmt.Errorf("menu not found")
	}

	for i, category := range menu.Categories {
		if category.ID == catID {
			for j, subcategory := range category.Sub {
				if subcategory.ID == subcatID {
					item.ID = primitive.NewObjectID()
					menu.Categories[i].Sub[j].Items = append(menu.Categories[i].Sub[j].Items, item)

					_, err = database.GetCollection("menus").UpdateOne(ctx, bson.M{"restaurant": objID}, bson.M{"$set": bson.M{"categories": menu.Categories}})
					if err != nil {
						return nil, fmt.Errorf("could not update menu")
					}

					return menu.ToPublic(), nil
				}
			}
		}
	}

	return nil, fmt.Errorf("subcategory not found")
}

func UpdateItemInMenu(ctx context.Context, restaurantID string, categoryID string, subcategoryID string, itemID string, item models.Item) (*models.PublicMenu, error) {
	var menu models.Menu

	objID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	catID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID")
	}

	subcatID, err := primitive.ObjectIDFromHex(subcategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid subcategory ID")
	}

	itemObjID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return nil, fmt.Errorf("invalid item ID")
	}

	err = database.GetCollection("menus").FindOne(ctx, bson.M{"restaurant": objID}).Decode(&menu)
	if err != nil {
		return nil, fmt.Errorf("menu not found")
	}

	for i, category := range menu.Categories {
		if category.ID == catID {
			for j, subcategory := range category.Sub {
				if subcategory.ID == subcatID {
					for k, dbitem := range subcategory.Items {
						if dbitem.ID == itemObjID {
							menu.Categories[i].Sub[j].Items[k] = models.Item{
								ID:          dbitem.ID,
								Name:        item.Name,
								Description: item.Description,
								Image:       item.Image,
								Price:       item.Price,
							}

							_, err = database.GetCollection("menus").UpdateOne(ctx, bson.M{"restaurant": objID}, bson.M{"$set": bson.M{"categories": menu.Categories}})
							if err != nil {
								return nil, fmt.Errorf("could not update menu")
							}

							return menu.ToPublic(), nil
						}
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("item not found")
}

func DeleteItemFromMenu(ctx context.Context, restaurantID string, categoryID string, subcategoryID string, itemID string) (*models.PublicMenu, error) {
	var menu models.Menu

	objID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	catID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID")
	}

	subcatID, err := primitive.ObjectIDFromHex(subcategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid subcategory ID")
	}

	itemObjID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return nil, fmt.Errorf("invalid item ID")
	}

	err = database.GetCollection("menus").FindOne(ctx, bson.M{"restaurant": objID}).Decode(&menu)
	if err != nil {
		return nil, fmt.Errorf("menu not found")
	}

	for i, category := range menu.Categories {
		if category.ID == catID {
			for j, subcategory := range category.Sub {
				if subcategory.ID == subcatID {
					for k, item := range subcategory.Items {
						if item.ID == itemObjID {
							menu.Categories[i].Sub[j].Items = append(menu.Categories[i].Sub[j].Items[:k], menu.Categories[i].Sub[j].Items[k+1:]...)

							_, err = database.GetCollection("menus").UpdateOne(ctx, bson.M{"restaurant": objID}, bson.M{"$set": bson.M{"categories": menu.Categories}})
							if err != nil {
								return nil, fmt.Errorf("could not update menu")
							}

							return menu.ToPublic(), nil
						}
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("item not found")
}
