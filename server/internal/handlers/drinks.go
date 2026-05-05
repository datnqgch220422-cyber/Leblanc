package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"leblanc/server/internal/db"
	"leblanc/server/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetDrinks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter, sortOptions := buildDrinkListQuery(c)
	findOptions := options.Find()
	if sortOptions != nil {
		findOptions.SetSort(sortOptions)
	}

	cur, err := db.DB.Collection("drinks").Find(ctx, filter, findOptions)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }

	var list []models.Drink
	if err := cur.All(ctx, &list); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, list)
}

func buildDrinkListQuery(c *gin.Context) (bson.M, bson.D) {
	filter := bson.M{}

	if q := strings.TrimSpace(c.Query("q")); q != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": q, "$options": "i"}},
			{"tags": bson.M{"$elemMatch": bson.M{"$regex": q, "$options": "i"}}},
		}
	}

	if available := strings.TrimSpace(c.Query("available")); available != "" {
		if val, err := strconv.ParseBool(available); err == nil {
			filter["available"] = val
		}
	}

	if stock := strings.TrimSpace(c.Query("stock")); stock != "" {
		switch strings.ToLower(stock) {
		case "positive", "in", "instock":
			filter["stock"] = bson.M{"$gt": 0}
		case "zero", "soldout", "out":
			filter["stock"] = bson.M{"$lte": 0}
		}
	}

	var sort bson.D
	switch strings.ToLower(strings.TrimSpace(c.Query("sort"))) {
	case "price":
		sort = bson.D{{Key: "price", Value: 1}}
	case "stock":
		sort = bson.D{{Key: "stock", Value: -1}}
	case "name":
		fallthrough
	default:
		sort = bson.D{{Key: "name", Value: 1}}
	}

	if order := strings.ToLower(strings.TrimSpace(c.Query("order"))); order == "desc" {
		for i := range sort {
			sort[i].Value = -1
		}
	}

	return filter, sort
}
