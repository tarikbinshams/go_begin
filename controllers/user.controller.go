package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"hello/config"
	"hello/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Email string             `json:"email" bson:"email"`
}

// Create User
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Get collection
	var userCollection *mongo.Collection = config.GetCollection("users")
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Get collection
	var userCollection *mongo.Collection = config.GetCollection("users")
	w.Header().Set("Content-Type", "application/json")

	var users []models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	var usersResponse []UserResponse
	for _, user := range users {
		usersResponse = append(usersResponse, UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	json.NewEncoder(w).Encode(usersResponse)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Get collection
	var userCollection *mongo.Collection = config.GetCollection("users")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	userResponse := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	json.NewEncoder(w).Encode(userResponse)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get collection
	var userCollection *mongo.Collection = config.GetCollection("users")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": user}
	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User updated successfully")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get collection
	var userCollection *mongo.Collection = config.GetCollection("users")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = userCollection.DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User deleted successfully")
}
