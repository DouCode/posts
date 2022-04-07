package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// swagger:parameters auth signIn
type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Username     string             `json:"username" bson:"username"`
	Telephone    string             `json:"telephone" bson:"telephone"`
	Password     string             `json:"password" bson:"password"`
	RegisteredAt time.Time          `json:"registeredAt" bson:"registeredAt"`
}
