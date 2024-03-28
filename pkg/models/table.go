package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Table struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Number    int                `json:"number"`
	Url       string             `json:"url"`
	Status    string             `json:"status"`
	Occupants []Client           `json:"occupants"`
}
