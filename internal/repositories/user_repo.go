package repositories

import (
	"context"
	"time"

	"go-backend/internal/database"
	"go-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository struct{}

func (r *UserRepository) CreateUser(user models.User) error {
	collection := database.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		println("❌ Insert error:", err.Error())
		return err
	}

	// println("✅ User inserted")
	println("✅ Inserted ID:", result.InsertedID.(string))
	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	collection := database.DB.Collection("users")

	var user models.User

	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
