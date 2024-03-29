package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Image       string             `json:"image"`
	Price       float64            `json:"price"`
}

type Category struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id"`
	Name  string             `json:"name"`
	Items []Item             `json:"items"`
}

type Menu struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	RestaurantID  primitive.ObjectID `json:"restaurant" bson:"restaurant"`
	Highlights    []Item             `json:"highlights"`
	Categories    []Category         `json:"categories"`
	SchemaVersion int                `json:"schema_version"`
	CreatedAt     int                `json:"created_at"`
	UpdatedAt     int                `json:"updated_at"`
}

type PublicMenu struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	RestaurantID primitive.ObjectID `json:"restaurant" bson:"restaurant"`
	Highlights   []Item             `json:"highlights"`
	Categories   []Category         `json:"categories"`
	CreatedAt    int                `json:"created_at"`
	UpdatedAt    int                `json:"updated_at"`
}

func (m *Menu) ToPublic() *PublicMenu {
	return &PublicMenu{
		ID:           m.ID,
		RestaurantID: m.RestaurantID,
		Highlights:   m.Highlights,
		Categories:   m.Categories,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
