package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() (*mongo.Client, error) {
	MongoDB := "mongodb://localhost:27017"
	fmt.Println("connecting to " + MongoDB + "...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDB))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("connected to database")
	return client, nil
}

var Client, Error = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection =  client.Database("restaurant").Collection(collectionName)
	return collection
}