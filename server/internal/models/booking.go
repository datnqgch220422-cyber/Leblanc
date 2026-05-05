package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingItem struct {
	DrinkID primitive.ObjectID `bson:"drinkId" json:"drinkId"`
	Qty     int                `bson:"qty" json:"qty"`
	Options map[string]any     `bson:"options" json:"options"`
}

type Booking struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Email          string             `bson:"email" json:"email"`
	Name           string             `bson:"name" json:"name"`
	Phone          string             `bson:"phone" json:"phone"`
	Time           time.Time          `bson:"time" json:"time"`
	Guests         int                `bson:"guests,omitempty" json:"guests,omitempty"`
	Items          []BookingItem      `bson:"items" json:"items"`
	Status         string             `bson:"status" json:"status"`
	Channel        string             `bson:"channel" json:"channel"`
	PaymentMethod  string             `bson:"paymentMethod,omitempty" json:"paymentMethod,omitempty"`
	PaymentStatus  string             `bson:"paymentStatus,omitempty" json:"paymentStatus,omitempty"`
	PaymentAmount  int64              `bson:"paymentAmount,omitempty" json:"paymentAmount,omitempty"`
	PaymentOrderID string             `bson:"paymentOrderId,omitempty" json:"paymentOrderId,omitempty"`
	PaymentTransID string             `bson:"paymentTransId,omitempty" json:"paymentTransId,omitempty"`
	PaymentMessage string             `bson:"paymentMessage,omitempty" json:"paymentMessage,omitempty"`
}
