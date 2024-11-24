package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	// Create a context with a timeout of 15 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Access the logs collection from the database
	collection := client.Database("logs").Collection("logs")

	// Set options for sorting by creation date in descending order
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Execute the find operation with proper context
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Printf("Error finding all documents: %v", err)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("Error closing cursor: %v", err)
		}
	}()

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		// Decode the current document into item
		if err := cursor.Decode(&item); err != nil {
			log.Printf("Error decoding log entry: %v", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}
