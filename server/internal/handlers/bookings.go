package handlers

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"time"

	"leblanc/server/internal/db"
	"leblanc/server/internal/models"
	"leblanc/server/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetBookingPaymentStatus(c *gin.Context) {
	paymentOrderID := strings.TrimSpace(c.Query("paymentOrderId"))
	if paymentOrderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "paymentOrderId is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var booking models.Booking
	if err := db.DB.Collection("bookings").FindOne(ctx, bson.M{"paymentOrderId": paymentOrderID}).Decode(&booking); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	items := make([]gin.H, 0, len(booking.Items))
	for _, item := range booking.Items {
		itemName := item.DrinkID.Hex()
		var drink models.Drink
		if err := db.DB.Collection("drinks").FindOne(ctx, bson.M{"_id": item.DrinkID}).Decode(&drink); err == nil {
			itemName = drink.Name
		}
		items = append(items, gin.H{
			"drinkId": item.DrinkID.Hex(),
			"name":    itemName,
			"qty":     item.Qty,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":             true,
		"paymentStatus":  booking.PaymentStatus,
		"status":         booking.Status,
		"paymentOrderId": booking.PaymentOrderID,
		"paymentTransId":  booking.PaymentTransID,
		"paymentMessage":  booking.PaymentMessage,
		"paymentAmount":   booking.PaymentAmount,
		"bookingId":      booking.ID.Hex(),
		"booking": gin.H{
			"bookingId": booking.ID.Hex(),
			"email":     booking.Email,
			"name":      booking.Name,
			"phone":     booking.Phone,
			"time":      booking.Time,
			"guests":    booking.Guests,
			"channel":   booking.Channel,
			"items":     items,
		},
	})
}

func GetBookings(c *gin.Context) {
	email := strings.TrimSpace(c.Query("email"))
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.Find().SetSort(bson.D{{Key: "time", Value: -1}})
	cur, err := db.DB.Collection("bookings").Find(ctx, bson.M{"email": email}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)

	var list []models.Booking
	if err := cur.All(ctx, &list); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if list == nil {
		list = []models.Booking{}
	}

	if len(list) > 1 {
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].Time.After(list[j].Time)
		})
	}

	c.JSON(http.StatusOK, list)
}

func CreateBooking(c *gin.Context) {
	var b models.Booking
	if err := c.BindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	b.Email = strings.TrimSpace(b.Email)
	if b.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}
	b.Status = models.NormalizeBookingStatus(b.Status)
	if b.Channel == "" {
		b.Channel = "web"
	}
	b.PaymentMethod = strings.TrimSpace(strings.ToLower(b.PaymentMethod))
	b.PaymentStatus = models.PaymentStatusPending
	b.PaymentMessage = "pending payment"
	b.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	totalAmount := int64(0)
	if b.PaymentMethod == "vnpay" {
		for _, item := range b.Items {
			var drink models.Drink
			if err := db.DB.Collection("drinks").FindOne(ctx, bson.M{"_id": item.DrinkID}).Decode(&drink); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "drink not found"})
				return
			}
			if !drink.Available || drink.Stock < item.Qty {
				c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient stock for " + drink.Name})
				return
			}
			totalAmount += int64(drink.Price * item.Qty)
		}
		b.PaymentAmount = totalAmount
	}

	res, err := db.DB.Collection("bookings").InsertOne(ctx, b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if b.PaymentMethod == "vnpay" {
		bookingIDHex := b.ID.Hex()
		payUrl, err := services.CreateVNPayURL(
			bookingIDHex,
			totalAmount,
			"Thanh toan don hang LeBlanc",
			c.ClientIP(),
		)
		if err != nil {
			_, _ = db.DB.Collection("bookings").DeleteOne(ctx, bson.M{"_id": b.ID})
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, _ = db.DB.Collection("bookings").UpdateOne(
			ctx,
			bson.M{"_id": b.ID},
			bson.M{"$set": bson.M{
				"paymentOrderId": bookingIDHex,
				"paymentMessage": "pending vnpay payment",
			}},
		)

		c.JSON(http.StatusOK, gin.H{
			"ok":             true,
			"id":             res.InsertedID,
			"paymentStatus":  models.PaymentStatusPending,
			"paymentOrderId": bookingIDHex,
			"payUrl":         payUrl,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "id": res.InsertedID})
}

// CancelBooking allows a user to cancel their booking by id.
// Expects JSON body { "email": "user@example.com" } to verify ownership.
func CancelBooking(c *gin.Context) {
	idParam := strings.TrimSpace(c.Param("id"))
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking id is required"})
		return
	}

	// attempt to parse as ObjectID
	oid, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	var payload struct{
		Email string `json:"email"`
	}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var booking models.Booking
	if err := db.DB.Collection("bookings").FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return
	}

	// simple ownership check: email in body must match booking email
	if strings.TrimSpace(strings.ToLower(payload.Email)) != strings.TrimSpace(strings.ToLower(booking.Email)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "email does not match booking owner"})
		return
	}

	// Only allow cancel if not already cancelled or completed
	status := strings.ToLower(booking.Status)
	if status == "cancelled" || status == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking cannot be cancelled"})
		return
	}

	// perform update
	_, err = db.DB.Collection("bookings").UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"status": "cancelled"}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "status": "cancelled"})
}
