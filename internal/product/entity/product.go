package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Price       string             `bson:"price,omitempty" json:"price,omitempty"`
	Tag         []string           `bson:"tag,omitempty" json:"tag,omitempty"`
	Discount    string             `bson:"discount,omitempty" json:"discount,omitempty"`
	Stock       int                `bson:"stock,omitempty" json:"stock,omitempty"`
	Image       []string           `bson:"image,omitempty" json:"image,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	CreatedBy   string             `bson:"created_by,omitempty" json:"created_by,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
