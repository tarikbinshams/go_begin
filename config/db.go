package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("❌ Failed to connect to MongoDB:", err)
	}

	// Ping the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ Failed to ping MongoDB:", err)
	}

	fmt.Println("✅ Connected to MongoDB!")
	DB = client
}

// Function to get a collection
func GetCollection(collectionName string) *mongo.Collection {
	// ConnectDB()
	if DB == nil {
		log.Fatal("❌ Database is not initialized. Call ConnectDB() first.")

	}
	return DB.Database("go").Collection(collectionName)
}
