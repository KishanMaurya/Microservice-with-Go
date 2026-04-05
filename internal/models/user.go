package models

import "time"

type User struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"-" bson:"password"` // 🔥 MUST match
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}