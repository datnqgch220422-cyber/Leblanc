package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"leblanc/server/internal/db"
	"leblanc/server/internal/models"
	"leblanc/server/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCart(c *gin.Context) {
	email := strings.TrimSpace(c.Query("email"))
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cart models.Cart
	err := db.DB.Collection("carts").FindOne(ctx, bson.M{"email": email}).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			cart = models.Cart{Email: email, Items: []models.CartItem{}}
			c.JSON(http.StatusOK, cart)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func AddCartItem(c *gin.Context) {
	var p struct {
		Email   string `json:"email"`
		DrinkID string `json:"drinkId"`
		Qty     int    `json:"qty"`
	}
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.Email = strings.TrimSpace(p.Email)
	if p.Email == "" || p.DrinkID == "" || p.Qty <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email, drinkId and positive qty are required"})
		return
	}

	did, err := primitive.ObjectIDFromHex(p.DrinkID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid drinkId"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verify drink exists and has enough stock
	var drink models.Drink
	if err := db.DB.Collection("drinks").FindOne(ctx, bson.M{"_id": did}).Decode(&drink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "drink not found"})
		return
	}
	if !drink.Available || drink.Stock < p.Qty {
		c.JSON(http.StatusBadRequest, gin.H{"error": "drink not available or insufficient stock"})
		return
	}

	// Load existing cart (if any)
	var cart models.Cart
	err = db.DB.Collection("carts").FindOne(ctx, bson.M{"email": p.Email}).Decode(&cart)
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cart.Items == nil {
		cart.Items = []models.CartItem{}
	}

	found := false
	for i := range cart.Items {
		if cart.Items[i].DrinkID == did {
			cart.Items[i].Qty += p.Qty
			found = true
			break
		}
	}
	if !found {
		cart.Items = append(cart.Items, models.CartItem{DrinkID: did, Qty: p.Qty})
	}
	cart.Email = p.Email
	cart.UpdatedAt = time.Now()

	// replace/upsert cart
	opts := options.Replace().SetUpsert(true)
	_, err = db.DB.Collection("carts").ReplaceOne(ctx, bson.M{"email": p.Email}, cart, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// helper: update takes qty, removes if qty <=0
func UpdateCartItem(c *gin.Context) {
	var p struct {
		Email   string `json:"email"`
		DrinkID string `json:"drinkId"`
		Qty     int    `json:"qty"`
	}
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.Email = strings.TrimSpace(p.Email)
	if p.Email == "" || p.DrinkID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and drinkId are required"})
		return
	}
	did, err := primitive.ObjectIDFromHex(p.DrinkID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid drinkId"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cart models.Cart
	err = db.DB.Collection("carts").FindOne(ctx, bson.M{"email": p.Email}).Decode(&cart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart not found"})
		return
	}

	// modify
	changed := false
	newItems := []models.CartItem{}
	for _, it := range cart.Items {
		if it.DrinkID == did {
			if p.Qty > 0 {
				it.Qty = p.Qty
				newItems = append(newItems, it)
				changed = true
			} else {
				changed = true
			}
		} else {
			newItems = append(newItems, it)
		}
	}
	if !changed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "item not in cart"})
		return
	}

	cart.Items = newItems
	cart.UpdatedAt = time.Now()
	_, err = db.DB.Collection("carts").ReplaceOne(ctx, bson.M{"email": p.Email}, cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func RemoveCartItem(c *gin.Context) {
	var p struct {
		Email   string `json:"email"`
		DrinkID string `json:"drinkId"`
	}
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.Email = strings.TrimSpace(p.Email)
	if p.Email == "" || p.DrinkID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and drinkId are required"})
		return
	}
	did, err := primitive.ObjectIDFromHex(p.DrinkID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid drinkId"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cart models.Cart
	err = db.DB.Collection("carts").FindOne(ctx, bson.M{"email": p.Email}).Decode(&cart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart not found"})
		return
	}
	newItems := []models.CartItem{}
	removed := false
	for _, it := range cart.Items {
		if it.DrinkID == did {
			removed = true
			continue
		}
		newItems = append(newItems, it)
	}
	if !removed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "item not in cart"})
		return
	}
	cart.Items = newItems
	cart.UpdatedAt = time.Now()
	_, err = db.DB.Collection("carts").ReplaceOne(ctx, bson.M{"email": p.Email}, cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func Checkout(c *gin.Context) {
	var p struct {
		Email         string `json:"email"`
		Name          string `json:"name"`
		Phone         string `json:"phone"`
		Channel       string `json:"channel"`
		PaymentMethod string `json:"paymentMethod"`
	}
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.Email = strings.TrimSpace(p.Email)
	if p.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cart models.Cart
	if err := db.DB.Collection("carts").FindOne(ctx, bson.M{"email": p.Email}).Decode(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart not found"})
		return
	}
	if len(cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no drinks selected for this booking"})
		return
	}

	// Validate stock and prepare booking items
	var booking models.Booking
	booking.Email = p.Email
	booking.Name = p.Name
	booking.Phone = p.Phone
	booking.Time = time.Now()
	booking.Status = models.BookingStatusPending
	booking.PaymentStatus = models.PaymentStatusPending
	booking.Channel = strings.TrimSpace(p.Channel)
	if booking.Channel == "" {
		booking.Channel = "web-booking"
	}
	booking.PaymentMethod = strings.TrimSpace(strings.ToLower(p.PaymentMethod))
	if booking.PaymentMethod == "" {
		booking.PaymentMethod = "counter"
	}
	booking.Items = []models.BookingItem{}
	var totalAmount int64 = 0

	for _, it := range cart.Items {
		var drink models.Drink
		if err := db.DB.Collection("drinks").FindOne(ctx, bson.M{"_id": it.DrinkID}).Decode(&drink); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "drink not found"})
			return
		}
		if !drink.Available || drink.Stock < it.Qty {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient stock for " + drink.Name})
			return
		}
		totalAmount += int64(drink.Price * it.Qty)
		booking.Items = append(booking.Items, models.BookingItem{DrinkID: it.DrinkID, Qty: it.Qty})
	}

	booking.PaymentAmount = totalAmount

	// Insert booking
	booking.ID = primitive.NewObjectID()
	booking.PaymentOrderID = booking.ID.Hex()
	booking.PaymentMessage = "pending vnpay payment"
	res, err := db.DB.Collection("bookings").InsertOne(ctx, booking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bookingID, _ := res.InsertedID.(primitive.ObjectID)
	bookingIDHex := bookingID.Hex()

	if booking.PaymentMethod == "vnpay" {
		ipAddr := c.ClientIP()
		payUrl, err := services.CreateVNPayURL(
			bookingIDHex,
			totalAmount,
			"Thanh toan don hang LeBlanc",
			ipAddr,
		)
		if err != nil {
			_, _ = db.DB.Collection("bookings").DeleteOne(ctx, bson.M{"_id": bookingID})
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, _ = db.DB.Collection("bookings").UpdateOne(
			ctx,
			bson.M{"_id": bookingID},
			bson.M{"$set": bson.M{
				"paymentOrderId": booking.PaymentOrderID,
				"paymentMessage": "pending vnpay payment",
			}},
		)

		c.JSON(http.StatusOK, gin.H{
			"ok":             true,
			"message":        "Booking created, waiting for VNPay payment",
			"id":             bookingIDHex,
			"paymentMethod":  "vnpay",
			"paymentStatus":  models.PaymentStatusPending,
			"paymentOrderId": booking.PaymentOrderID,
			"payUrl":         payUrl,
		})
		return
	}

	// Decrement stock
	for _, it := range cart.Items {
		_, _ = db.DB.Collection("drinks").UpdateOne(ctx, bson.M{"_id": it.DrinkID}, bson.M{"$inc": bson.M{"stock": -it.Qty}})
	}

	// Clear cart
	_, _ = db.DB.Collection("carts").DeleteOne(ctx, bson.M{"email": p.Email})

	_, _ = db.DB.Collection("bookings").UpdateOne(
		ctx,
		bson.M{"_id": bookingID},
		bson.M{"$set": bson.M{"paymentStatus": models.PaymentStatusPaid}},
	)

	c.JSON(http.StatusOK, gin.H{"ok": true, "id": res.InsertedID, "paymentStatus": models.PaymentStatusPaid})
}
