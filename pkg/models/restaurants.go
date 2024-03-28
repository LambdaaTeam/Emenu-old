package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	TableStatusAvailable = "AVAILABLE"
	TableStatusReserved  = "RESERVED"
	TableStatusOccupied  = "OCCUPIED"
)

type Restaurant struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Name     string             `json:"name"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
	Address  struct {
		City     string `json:"city"`
		Country  string `json:"country"`
		PostCode string `json:"postCode"`
		Number   string `json:"number"`
		Street   string `json:"street"`
		Other    string `json:"other"`
	}
	Tables        []Table `json:"tables"`
	SchemaVersion int     `json:"schemaVersion"`
	CreatedAt     int     `json:"createdAt"`
	UpdatedAt     int     `json:"updatedAt"`
}

type PublicRestaurant struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	Name    string             `json:"name"`
	Address struct {
		City     string `json:"city"`
		Country  string `json:"country"`
		PostCode string `json:"postCode"`
		Number   string `json:"number"`
		Street   string `json:"street"`
		Other    string `json:"other"`
	}
	Tables    []Table `json:"tables"`
	CreatedAt int     `json:"createdAt"`
	UpdatedAt int     `json:"updatedAt"`
}

func (r *Restaurant) ToPublic() *PublicRestaurant {
	return &PublicRestaurant{
		ID:        r.ID,
		Name:      r.Name,
		Address:   r.Address,
		Tables:    r.Tables,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
