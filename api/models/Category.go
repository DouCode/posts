package models

type Category struct {
	ID        uint   `json:"id" gorm:"primary_key" bson:"id"`
	Name      string `json:"name" gorm:"type:varchar(50);not null;unique" bson:"name"`
	CreatedAt Time   `json:"created_at" gorm:"type:timestamp" bson:"createdAt"`
	UpdatedAt Time   `json:"updated_at" gorm:"type:timestamp" bson:"updatedAt"`
}
