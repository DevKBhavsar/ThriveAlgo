package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
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

func retrieve(id string) (Holiday, error) {
	collection := client.Database("ThriveAlgoData").Collection("holidays")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var holiday Holiday
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&holiday)
	if err == mongo.ErrNoDocuments {
		return Holiday{}, nil
	}
	return holiday, err
}

func (h *Holiday) create() error {
	collection := client.Database("ThriveAlgoData").Collection("holidays")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, h)
	return err
}

func (h *Holiday) update() error {
	collection := client.Database("ThriveAlgoData").Collection("holidays")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": h.ID}
	update := bson.M{"$set": bson.M{
		"date":        h.Date,
		"title":       h.Title,
		"description": h.Description,
	}}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("holiday not found")
	}
	return nil
}

func deleteHoliday(id string) error {
	collection := client.Database("ThriveAlgoData").Collection("holidays")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("holiday not found")
	}
	return nil
}
