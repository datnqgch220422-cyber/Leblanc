package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EmotionFit struct {
	Calm        float64 `bson:"calm" json:"calm"`
	Happy       float64 `bson:"happy" json:"happy"`
	Stressed    float64 `bson:"stressed" json:"stressed"`
	Sad         float64 `bson:"sad" json:"sad"`
	Adventurous float64 `bson:"adventurous" json:"adventurous"`
}

type Drink struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name       string             `bson:"name" json:"name"`
	Price      int                `bson:"price" json:"price"`
	Stock      int                `bson:"stock" json:"stock"`
	Available  bool               `bson:"available" json:"available"`
	Tags       []string           `bson:"tags" json:"tags"`
	Caffeine   string             `bson:"caffeine" json:"caffeine"` // low|med|high
	Temp       string             `bson:"temp" json:"temp"`         // hot|iced|either
	Sweetness  int                `bson:"sweetness" json:"sweetness"`
	ColorTone  string             `bson:"colorTone" json:"colorTone"` // warm|cool|neutral
	EmotionFit EmotionFit         `bson:"emotionFit" json:"emotionFit"`
	Image      string             `bson:"image" json:"image"`
	Desc       string             `bson:"desc" json:"desc"`
}
