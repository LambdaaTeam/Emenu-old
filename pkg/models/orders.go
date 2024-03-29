package models

import "go.mongodb.org/mongo-driver/bson/primitive"

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
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	RestaurantID  primitive.ObjectID `json:"restaurant" bson:"restaurant"`
	TableID       string             `json:"table"`
	Status        string             `json:"status"`
	Value         float64            `json:"value"`
	Client        Client             `json:"client"`
	Items         []OrderItem        `json:"items"`
	SchemaVersion int                `json:"schema_version"`
	CreatedAt     int                `json:"created_at"`
	UpdatedAt     int                `json:"updated_at"`
}

type OrderItem struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name"`
	Price       float64            `json:"price"`
	Quantity    int                `json:"quantity"`
	Status      string             `json:"status"`
	Observation string             `json:"observation"`
}

type PublicOrder struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	RestaurantID primitive.ObjectID `json:"restaurant" bson:"restaurant"`
	TableID      string             `json:"table"`
	Status       string             `json:"status"`
	Value        float64            `json:"value"`
	Client       Client             `json:"client"`
	Items        []OrderItem        `json:"items"`
	CreatedAt    int                `json:"created_at"`
	UpdatedAt    int                `json:"updated_at"`
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
