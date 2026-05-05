package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItem struct {
    DrinkID primitive.ObjectID `bson:"drinkId" json:"drinkId"`
    Qty     int                `bson:"qty" json:"qty"`
}

type Cart struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
    Email     string             `bson:"email" json:"email"`
    Items     []CartItem         `bson:"items" json:"items"`
    UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
