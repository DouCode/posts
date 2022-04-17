package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Form struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Telephone    string             `json:"telephone" bson:"telephone"`
	Password     string             `json:"password" bson:"password"`
	InviteCode   string             `json:"inviteCode" bson:"inviteCode"`
	Mail         string             `json:"mail" bson:"mail"`
	RegisteredAt time.Time          `json:"registeredAt" bson:"registeredAt"`
}
