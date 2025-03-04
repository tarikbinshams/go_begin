package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"hello/config"
	"hello/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("secret-key")

type AuthClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func createToken(email string) (string, error) {
	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		fmt.Printf("Error generating token string: %v", err)
		return "", err
	}
	return tokenString, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var userCollection *mongo.Collection = config.GetCollection("users")
	w.Header().Set("Content-Type", "application/json")

	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = userCollection.FindOne(ctx, bson.M{"email": credentials.Email}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := createToken(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Invalid credentials")
	}

	// Return token
	response := map[string]string{
		"token": token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var userCollection *mongo.Collection = config.GetCollection("users")
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User created successfully")
}
