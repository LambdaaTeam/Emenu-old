package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const OrderCurrentSchemaVersion = 1

const (
	OrderStatusOpen   = "PENDING"
	OrderStatusClosed = "CLOSED"
)

const (
	ItemStatusToPrepare = "TO_PREPARE"
	ItemStatusPreparing = "PREPARING"
	ItemStatusReady     = "READY"
	ItemStatusDelivered = "DELIVERED"
)

type Order struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	RestaurantID  primitive.ObjectID `json:"restaurant" bson:"restaurant"`
	TableID       primitive.ObjectID `json:"table" bson:"table"`
	Status        string             `json:"status"`
	Value         float64            `json:"value"`
	Client        Client             `json:"client"`
	Items         []OrderItem        `json:"items"`
	SchemaVersion int                `json:"schema_version"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

type OrderItem struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Quantity    int                `json:"quantity"`
	Status      string             `json:"status"`
	Observation string             `json:"observation"`
}

type PublicOrder struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	RestaurantID primitive.ObjectID `json:"restaurant" bson:"restaurant"`
	TableID      primitive.ObjectID `json:"table" bson:"table"`
	Status       string             `json:"status"`
	Value        float64            `json:"value"`
	Client       Client             `json:"client"`
	Items        []OrderItem        `json:"items"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

type PublicOrderWithToken struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	RestaurantID primitive.ObjectID `json:"restaurant" bson:"restaurant"`
	TableID      primitive.ObjectID `json:"table" bson:"table"`
	Status       string             `json:"status"`
	Value        float64            `json:"value"`
	Client       Client             `json:"client"`
	Items        []OrderItem        `json:"items"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	Token        string             `json:"token"`
}

func (o *Order) ToPublic() *PublicOrder {
	return &PublicOrder{
		ID:           o.ID,
		RestaurantID: o.RestaurantID,
		TableID:      o.TableID,
		Status:       o.Status,
		Value:        o.Value,
		Client:       o.Client,
		Items:        o.Items,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}
}


func (r *PublicOrder) AddToken(token string) *PublicOrderWithToken {
	return &PublicOrderWithToken{
		ID: r.ID,
		RestaurantID: r.RestaurantID,
		TableID: r.TableID,
		Status: r.Status,
		Value: r.Value,
		Client: r.Client,
		Items: r.Items,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		Token:     token,
	}
}

