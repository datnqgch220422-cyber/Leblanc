package handlers

import (
	"context"
	"net/http"
	"time"

	"leblanc/server/internal/db"
	"leblanc/server/internal/models"
	"leblanc/server/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func RecoFromFeatures(c *gin.Context) {
	type RecoRequest struct {
		Caffeine  string `json:"caffeine"`
		Temp      string `json:"temp"`
		Sweetness int    `json:"sweetness"`
	}

	var payload RecoRequest
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := db.DB.Collection("drinks").Find(ctx, bson.D{})
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }

	var drinks []models.Drink
	if err := cur.All(ctx, &drinks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}

	scores := services.ScoreDrinks(drinks, payload.Caffeine, payload.Temp, payload.Sweetness)

	// Build response
	type ScoreItem struct {
		DrinkID string  `json:"drinkId"`
		Score   float64 `json:"score"`
	}
	resp := make([]ScoreItem, len(scores))
	for i, s := range scores {
		resp[i] = ScoreItem{DrinkID: s.DrinkID, Score: s.Score}
	}
	c.JSON(http.StatusOK, resp)
}
