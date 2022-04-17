package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	UserId      int                `json:"user_id" gorm:"not null" bson:"userid"`
	CategoryId  int                `json:"category_id" gorm:"not null" bson:"categoryid"`
	Category    *Category
	Title       string    `json:"title" gorm:"type:varchar(50);not null" bson:"title"`
	HeadImg     string    `json:"head_img" bson:"headImg"`
	Content     string    `json:"content" gorm:"tyep:text;not null" bson:"content"`
	PublishedAt time.Time `json:"publishedAt" bson:"publishedAt"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamp" bson:"updatedAt"`
	Tags        string    `json:"tags" bson:"tags"`
	UserName    string    `json:"username" bson:"username"`
}
