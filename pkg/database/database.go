package database

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func Connect(dbname string) *mongo.Database {
	uri := os.Getenv("DATABASE_URL")

	if uri == "" {
		panic("DATABASE_URL is required")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		panic(err)
	}

	db := client.Database(dbname)

	fmt.Println("Connected to MongoDB!")

	return db
}

func GetCollection(collection string) *mongo.Collection {
	return DB.Collection(collection)
}
