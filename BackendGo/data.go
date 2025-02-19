// data.go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoClient() *mongo.Client {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not ping MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB")
	return client
}

func retrieveAll() ([]Holiday, error) {
	collection := client.Database("ThriveAlgoData").Collection("holidays")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var holidays []Holiday
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &holidays); err != nil {
		return nil, err
	}

	return holidays, nil
}

func (h *Holiday) create() error {
	collection := client.Database("ThriveAlgoData").Collection("holidays")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, h)
	return err
}

func deleteHoliday(id string) error {
	collection := client.Database("ThriveAlgoData").Collection("holidays")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("holiday not found")
	}
	return nil
}
