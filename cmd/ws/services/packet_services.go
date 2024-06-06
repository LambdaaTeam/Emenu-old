package services

import (
	"context"
	"fmt"

	"github.com/LambdaaTeam/Emenu/cmd/ws/shared"
	"github.com/LambdaaTeam/Emenu/pkg/database"
	"github.com/LambdaaTeam/Emenu/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandlePacket(ctx context.Context, packet *shared.Packet) shared.Packet {
	switch packet.Type {
	case shared.UpdateOrderStatus:
		return updateOrderStatus(ctx, packet.RestaurantID, packet.OrderID, packet.Data)
	case shared.UpdateItemStatus:
		return updateItemStatus(ctx, packet.RestaurantID, packet.OrderID, packet.ItemID, packet.Data)
	case shared.UpdateTableStatus:
		return updateTableStatus(ctx, packet.RestaurantID, packet.TableID, packet.Data)
	case shared.Heartbeat:
		return shared.Packet{
			Type: shared.Heartbeat,
			Data: "connected",
		}
	default:
		return shared.NewErrorPacket(packet.RestaurantID, "invalid packet type")
	}
}

func updateOrderStatus(ctx context.Context, restaurantID, orderID, status string) shared.Packet {
	if status != models.OrderStatusClosed && status != models.OrderStatusOpen {
		return shared.NewErrorPacket(restaurantID, "invalid status")
	}

	ordID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return shared.NewErrorPacket(restaurantID, err.Error())
	}

	restID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return shared.NewErrorPacket(restaurantID, err.Error())
	}

	var order models.Order

	if err := database.DB.Collection("orders").FindOne(ctx, bson.M{"_id": ordID}).Decode(&order); err != nil {
		return shared.NewErrorPacket(restaurantID, err.Error())
	}

	if order.RestaurantID != restID {
		return shared.NewErrorPacket(restaurantID, "restaurant id does not match")
	}

	if _, err := database.DB.Collection("orders").UpdateOne(ctx, bson.M{"_id": ordID}, bson.M{"$set": bson.M{"status": status}}); err != nil {
		return shared.NewErrorPacket(restaurantID, err.Error())
	}

	return shared.Packet{
		Type:         shared.UpdateOrderStatus,
		RestaurantID: restaurantID,
		OrderID:      orderID,
		Data:         status,
	}
}

func updateItemStatus(ctx context.Context, restaurantID string, orderID string, itemID string, status string) shared.Packet {
	if status != models.ItemStatusToPrepare && status != models.ItemStatusPreparing && status != models.ItemStatusReady && status != models.ItemStatusDelivered {
		return shared.NewErrorPacket(restaurantID, "invalid status")
	}

	resID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return shared.NewErrorPacket(restaurantID, err.Error())
	}

	ordID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return shared.NewErrorPacket(restaurantID, err.Error())
	}

	itmID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return shared.NewErrorPacket(restaurantID, err.Error())
	}

	var order models.Order

	if err := database.DB.Collection("orders").FindOne(ctx, bson.M{"_id": ordID}).Decode(&order); err != nil {
		return shared.NewErrorPacket(restaurantID, err.Error())
	}

	if order.RestaurantID != resID {
		return shared.NewErrorPacket(restaurantID, "restaurant id does not match")
	}

	var item models.OrderItem

	for _, i := range order.Items {
		if i.ID == itmID {
			item = i
			break
		}
	}

	if item.ID.IsZero() {
		return shared.NewErrorPacket(restaurantID, "item not found")
	}

	if _, err := database.DB.Collection("orders").UpdateOne(ctx, bson.M{"_id": ordID, "items._id": itmID}, bson.M{"$set": bson.M{"items.$.status": status}}); err != nil {
		return shared.NewErrorPacket(restaurantID, err.Error())
	}

	return shared.Packet{
		Type:         shared.UpdateItemStatus,
		RestaurantID: restaurantID,
		OrderID:      orderID,
		Data:         status,
	}
}

func updateTableStatus(ctx context.Context, restaurantID, tableID, data string) shared.Packet {
	if data != models.TableStatusOccupied && data != models.TableStatusAvailable && data != models.TableStatusReserved {
		return shared.NewErrorPacket(restaurantID, "invalid status")
	}

	resID, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return shared.NewErrorPacket(restaurantID, fmt.Errorf("invalid restaurant id: %w", err).Error())
	}

	tblID, err := primitive.ObjectIDFromHex(tableID)
	if err != nil {
		return shared.NewErrorPacket(restaurantID, fmt.Errorf("invalid table id: %w", err).Error())
	}

	if _, err := database.DB.Collection("restaurants").UpdateOne(ctx, bson.M{"_id": resID, "tables._id": tblID}, bson.M{"$set": bson.M{"tables.$.status": data}}); err != nil {
		return shared.NewErrorPacket(restaurantID, err.Error())
	}

	return shared.Packet{
		Type:         shared.UpdateTableStatus,
		RestaurantID: restaurantID,
		TableID:      tableID,
		Data:         data,
	}
}
