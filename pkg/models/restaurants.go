package models

import (
	"time"

	"github.com/LambdaaTeam/Emenu/pkg/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const RestaurantCurrentSchemaVersion = 1

const (
	TableStatusAvailable = "AVAILABLE"
	TableStatusReserved  = "RESERVED"
	TableStatusOccupied  = "OCCUPIED"
)

type Table struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Number    int                `json:"number"`
	Url       string             `json:"url"`
	Status    string             `json:"status"`
	Occupants []Client           `json:"occupants"`
}


type Restaurant struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name"`
	Email    string             `json:"email"`
	Password []byte             `json:"password"`
	Address  struct {
		City     string `json:"city"`
		Country  string `json:"country"`
		PostCode string `json:"post_code"`
		Number   int    `json:"number"`
		Street   string `json:"street"`
		Other    string `json:"other"`
	} `json:"address"`
	Tables        []Table   `json:"tables"`
	SchemaVersion int       `json:"schema_version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PublicRestaurant struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Name    string             `json:"name"`
	Address struct {
		City     string `json:"city"`
		Country  string `json:"country"`
		PostCode string `json:"post_code"`
		Number   int    `json:"number"`
		Street   string `json:"street"`
		Other    string `json:"other"`
	} `json:"address"`
	Tables    []Table   `json:"tables"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RestaurantRegister struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  struct {
		City     string `json:"city"`
		Country  string `json:"country"`
		PostCode string `json:"post_code"`
		Number   int    `json:"number"`
		Street   string `json:"street"`
		Other    string `json:"other"`
	} `json:"address"`
}

type RestaurantLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

func (r *RestaurantRegister) ToRestaurant() *Restaurant {
	return &Restaurant{
		ID:            primitive.NewObjectID(),
		Name:          r.Name,
		Email:         r.Email,
		Password:      auth.HashPassword(r.Password),
		Address:       r.Address,
		Tables:        []Table{},
		SchemaVersion: RestaurantCurrentSchemaVersion,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
