package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

	tableId := primitive.NewObjectID()

	shortnerBody, err := json.Marshal(map[string]string{
		"url": fmt.Sprintf("https://menu.emenu.psykka.xyz/?restaurantId=%s&table=%d&table_id=%s", id, number, tableId.Hex()),
	})

	if err != nil {
		return nil, fmt.Errorf("could not generate short URL")
	}

	shortnerResp, err := http.Post("https://short.emenu.psykka.xyz/", "application/json", bytes.NewBuffer(shortnerBody))

	if err != nil {
		return nil, fmt.Errorf("could not generate short URL")
	}
	defer shortnerResp.Body.Close()

	var shortnerRespBody map[string]string
	err = json.NewDecoder(shortnerResp.Body).Decode(&shortnerRespBody)
	if err != nil {
		return nil, fmt.Errorf("could not generate short URL")
	}

	table := models.Table{
		ID:        tableId,
		Number:    number,
		Url:       fmt.Sprintf("https://short.emenu.psykka.xyz/%s", shortnerRespBody["short"]),
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
			restaurant.Tables[i].Number = params.Number
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

func AddOrderItem(ctx context.Context, restaurantID string, orderID string, item models.OrderItem) (*models.PublicOrder, error) {
	var order models.Order

	if item.Quantity <= 0 {
		return nil, fmt.Errorf("invalid quantity")
	}

	orderObjID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID")
	}

	dbItem, err := GetItemFromMenu(ctx, restaurantID, item.ID.Hex())
	if err != nil {
		return nil, fmt.Errorf("item not found")
	}

	err = database.GetCollection("orders").FindOne(ctx, bson.M{"_id": orderObjID}).Decode(&order)
	if err != nil {
		return nil, fmt.Errorf("order not found")
	}

	order.Value += dbItem.Price * float64(item.Quantity)

	order.Items = append(order.Items, models.OrderItem{
		ID:          dbItem.ID,
		Quantity:    item.Quantity,
		Status:      models.ItemStatusToPrepare,
		Observation: item.Observation,
	})

	_, err = database.GetCollection("orders").UpdateOne(ctx, bson.M{"_id": orderObjID}, bson.M{"$set": bson.M{"items": order.Items}})
	if err != nil {
		return nil, fmt.Errorf("could not update order")
	}

	return order.ToPublic(), nil
}

func UpdateOrderItem(ctx context.Context, restaurantID string, orderID string, item models.OrderItem) (*models.PublicOrder, error) {
	var order models.Order

	orderObjID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID")
	}

	err = database.GetCollection("orders").FindOne(ctx, bson.M{"_id": orderObjID}).Decode(&order)
	if err != nil {
		return nil, fmt.Errorf("order not found")
	}

	for i, orderItem := range order.Items {
		if orderItem.ID == item.ID {
			if item.Quantity != 0 {
				order.Items[i].Quantity = item.Quantity
			}

			if item.Status != "" {
				if item.Status != models.ItemStatusToPrepare &&
					item.Status != models.ItemStatusPreparing &&
					item.Status != models.ItemStatusReady &&
					item.Status != models.ItemStatusDelivered {
					return nil, fmt.Errorf("invalid item status")
				}

				order.Items[i].Status = item.Status
			}

			order.Items[i].Observation = item.Observation

			_, err = database.GetCollection("orders").UpdateOne(ctx, bson.M{"_id": orderObjID}, bson.M{"$set": bson.M{"items": order.Items}})
			if err != nil {
				return nil, fmt.Errorf("could not update order")
			}

			return order.ToPublic(), nil
		}
	}

	return nil, fmt.Errorf("item not found")
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

func AddItemToMenu(ctx context.Context, restaurantID string, categoryID string, item models.Item) (*models.PublicMenu, error) {
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
			item.ID = primitive.NewObjectID()
			menu.Categories[i].Items = append(menu.Categories[i].Items, item)

			_, err = database.GetCollection("menus").UpdateOne(ctx, bson.M{"restaurant": objID}, bson.M{"$set": bson.M{"categories": menu.Categories}})
			if err != nil {
				return nil, fmt.Errorf("could not update menu")
			}

			return menu.ToPublic(), nil
		}
	}

	return nil, fmt.Errorf("category not found")
}

func GetItemFromMenu(ctx context.Context, restaurantID string, itemID string) (*models.Item, error) {
	var menu models.Menu

	objID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	itemObjID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return nil, fmt.Errorf("invalid item ID")
	}

	err = database.GetCollection("menus").FindOne(ctx, bson.M{"restaurant": objID}).Decode(&menu)
	if err != nil {
		return nil, fmt.Errorf("menu not found")
	}

	for _, category := range menu.Categories {
		for _, item := range category.Items {
			if item.ID == itemObjID {
				return &item, nil
			}
		}
	}

	for _, item := range menu.Highlights {
		if item.ID == itemObjID {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("item not found")
}

func UpdateItemInMenu(ctx context.Context, restaurantID string, categoryID string, itemID string, item models.Item) (*models.PublicMenu, error) {
	var menu models.Menu

	objID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	catID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID")
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
			for j, dbitem := range category.Items {
				if dbitem.ID == itemObjID {
					menu.Categories[i].Items[j] = models.Item{
						ID:          itemObjID,
						Name:        item.Name,
						Description: item.Description,
						Price:       item.Price,
						Image:       item.Image,
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

	return nil, fmt.Errorf("item not found")
}

func DeleteItemFromMenu(ctx context.Context, restaurantID string, categoryID string, itemID string) (*models.PublicMenu, error) {
	var menu models.Menu

	objID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID")
	}

	catID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID")
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
			for j, item := range category.Items {
				if item.ID == itemObjID {
					menu.Categories[i].Items = append(menu.Categories[i].Items[:j], menu.Categories[i].Items[j+1:]...)

					_, err = database.GetCollection("menus").UpdateOne(ctx, bson.M{"restaurant": objID}, bson.M{"$set": bson.M{"categories": menu.Categories}})
					if err != nil {
						return nil, fmt.Errorf("could not update menu")
					}

					return menu.ToPublic(), nil
				}
			}
		}
	}

	return nil, fmt.Errorf("item not found")
}

func AddClientToTable(ctx context.Context, restaurantID string, tableID string, client models.Client) (*models.PublicOrder, error) {
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

	for i, table := range restaurant.Tables {
		if table.ID == tableObjID {
			restaurant.Tables[i].Occupants = append(restaurant.Tables[i].Occupants, client)
			restaurant.Tables[i].Status = models.TableStatusOccupied

			_, err = database.GetCollection("restaurants").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"tables": restaurant.Tables}})
			if err != nil {
				return nil, fmt.Errorf("could not update restaurant tables")
			}

			orderId := primitive.NewObjectID()
			order := models.Order{
				ID:            orderId,
				RestaurantID:  objID,
				TableID:       tableObjID,
				Status:        models.OrderStatusOpen,
				Items:         []models.OrderItem{},
				Client:        client,
				SchemaVersion: models.OrderCurrentSchemaVersion,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			_, err = database.GetCollection("orders").InsertOne(ctx, order)
			if err != nil {
				return nil, fmt.Errorf("could not create order")
			}

			packet, err := json.Marshal(map[string]interface{}{
				"type":          "update_table_status",
				"restaurant_id": restaurantID,
				"order_id":      orderId.Hex(),
				"data":          models.TableStatusOccupied,
			})

			if err != nil {
				return nil, fmt.Errorf("could not notify ws")
			}

			response, err := http.Post("https://ws.emenu.psykka.xyz/notify", "application/json", bytes.NewBuffer(packet))
			if err != nil {
				return nil, fmt.Errorf("could not notify ws")
			}
			defer response.Body.Close()

			return order.ToPublic(), nil
		}
	}

	return nil, fmt.Errorf("table not found")
}
