package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"leblanc/server/internal/db"
	"leblanc/server/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type adminUserPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Verified bool   `json:"verified"`
}

type adminBookingItemPayload struct {
	DrinkID string         `json:"drinkId"`
	Qty     int            `json:"qty"`
	Options map[string]any `json:"options"`
}

type adminBookingPayload struct {
	Email   string                    `json:"email"`
	Name    string                    `json:"name"`
	Phone   string                    `json:"phone"`
	Time    string                    `json:"time"`
	Guests  int                       `json:"guests"`
	Status  string                    `json:"status"`
	Items   []adminBookingItemPayload `json:"items"`
	Channel string                    `json:"channel"`
}

type adminDrinkPayload struct {
	Name      string   `json:"name"`
	Price     int      `json:"price"`
	Stock     int      `json:"stock"`
	Available bool     `json:"available"`
	Tags      []string `json:"tags"`
	Caffeine  string   `json:"caffeine"`
	Temp      string   `json:"temp"`
	Sweetness int      `json:"sweetness"`
	ColorTone string   `json:"colorTone"`
	Image     string   `json:"image"`
	Desc      string   `json:"desc"`
}

func AdminListUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := optionsFindByCreatedAtDesc()
	cur, err := db.DB.Collection("users").Find(ctx, bson.D{}, opts.toFindOptions())
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer cur.Close(ctx)

	var users []models.User
	if err := cur.All(ctx, &users); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	publicUsers := make([]models.PublicUser, len(users))
	for i, u := range users {
		publicUsers[i] = u.Public()
	}

	respondData(c, http.StatusOK, publicUsers)
}

func AdminCreateUser(c *gin.Context) {
	var payload adminUserPayload
	if err := c.BindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	payload.Name = strings.TrimSpace(payload.Name)
	payload.Email = strings.TrimSpace(payload.Email)
	payload.Password = strings.TrimSpace(payload.Password)

	if payload.Name == "" || payload.Email == "" || payload.Password == "" {
		respondError(c, http.StatusBadRequest, "name, email and password are required")
		return
	}
	if !isValidEmail(payload.Email) {
		respondError(c, http.StatusBadRequest, "invalid email")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not hash password")
		return
	}

	user := models.User{
		ID:           primitive.NewObjectID(),
		Name:         payload.Name,
		NameLower:    strings.ToLower(payload.Name),
		Email:        payload.Email,
		EmailLower:   strings.ToLower(payload.Email),
		Role:         normalizeUserRole(payload.Role),
		Verified:     payload.Verified,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := db.DB.Collection("users").InsertOne(ctx, user); err != nil {
		status, msg := normalizeMongoWriteError(err, "could not create user")
		respondError(c, status, msg)
		return
	}

	respondMessage(c, http.StatusCreated, "user created", user.Public())
}

func AdminUpdateUser(c *gin.Context) {
	userID, ok := parseHexID(c, "id")
	if !ok {
		return
	}

	var payload adminUserPayload
	if err := c.BindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	update := bson.M{}
	if name := strings.TrimSpace(payload.Name); name != "" {
		update["name"] = name
		update["nameLower"] = strings.ToLower(name)
	}
	if email := strings.TrimSpace(payload.Email); email != "" {
		if !isValidEmail(email) {
			respondError(c, http.StatusBadRequest, "invalid email")
			return
		}
		update["email"] = email
		update["emailLower"] = strings.ToLower(email)
	}
	update["role"] = normalizeUserRole(payload.Role)
	update["verified"] = payload.Verified

	if password := strings.TrimSpace(payload.Password); password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			respondError(c, http.StatusInternalServerError, "could not hash password")
			return
		}
		update["passwordHash"] = string(hash)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := db.DB.Collection("users")
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": update}); err != nil {
		status, msg := normalizeMongoWriteError(err, "could not update user")
		respondError(c, status, msg)
		return
	}

	var user models.User
	if err := coll.FindOne(ctx, bson.M{"_id": userID}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			respondError(c, http.StatusNotFound, "user not found")
			return
		}
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondMessage(c, http.StatusOK, "user updated", user.Public())
}

func AdminDeleteUser(c *gin.Context) {
	userID, ok := parseHexID(c, "id")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := db.DB.Collection("users").DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if res.DeletedCount == 0 {
		respondError(c, http.StatusNotFound, "user not found")
		return
	}

	respondMessage(c, http.StatusOK, "user deleted", nil)
}

func AdminListBookings(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := db.DB.Collection("bookings").Find(ctx, bson.D{}, optionsFindByTimeDesc().toFindOptions())
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer cur.Close(ctx)

	var bookings []models.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondData(c, http.StatusOK, bookings)
}

func AdminCreateBooking(c *gin.Context) {
	var payload adminBookingPayload
	if err := c.BindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	booking, err := buildBookingFromPayload(payload)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	booking.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := db.DB.Collection("bookings").InsertOne(ctx, booking); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondMessage(c, http.StatusCreated, "booking created", booking)
}

func AdminUpdateBooking(c *gin.Context) {
	bookingID, ok := parseHexID(c, "id")
	if !ok {
		return
	}

	var payload adminBookingPayload
	if err := c.BindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	booking, err := buildBookingFromPayload(payload)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	update := bson.M{
		"email":   booking.Email,
		"name":    booking.Name,
		"phone":   booking.Phone,
		"time":    booking.Time,
		"guests":  booking.Guests,
		"items":   booking.Items,
		"status":  booking.Status,
		"channel": booking.Channel,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := db.DB.Collection("bookings")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": bookingID}, bson.M{"$set": update})
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if res.MatchedCount == 0 {
		respondError(c, http.StatusNotFound, "booking not found")
		return
	}

	var updated models.Booking
	if err := coll.FindOne(ctx, bson.M{"_id": bookingID}).Decode(&updated); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondMessage(c, http.StatusOK, "booking updated", updated)
}

func AdminDeleteBooking(c *gin.Context) {
	bookingID, ok := parseHexID(c, "id")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := db.DB.Collection("bookings").DeleteOne(ctx, bson.M{"_id": bookingID})
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if res.DeletedCount == 0 {
		respondError(c, http.StatusNotFound, "booking not found")
		return
	}

	respondMessage(c, http.StatusOK, "booking deleted", nil)
}

func AdminListDrinks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter, sortOptions := buildDrinkListQuery(c)
	findOptions := options.Find()
	findOptions.SetSort(sortOptions)

	cur, err := db.DB.Collection("drinks").Find(ctx, filter, findOptions)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer cur.Close(ctx)

	var drinks []models.Drink
	if err := cur.All(ctx, &drinks); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondData(c, http.StatusOK, drinks)
}

func AdminCreateDrink(c *gin.Context) {
	var payload adminDrinkPayload
	if err := c.BindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	drink, err := buildDrinkFromPayload(payload)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	drink.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := db.DB.Collection("drinks").InsertOne(ctx, drink); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondMessage(c, http.StatusCreated, "product created", drink)
}

func AdminUpdateDrink(c *gin.Context) {
	drinkID, ok := parseHexID(c, "id")
	if !ok {
		return
	}

	var payload adminDrinkPayload
	if err := c.BindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	drink, err := buildDrinkFromPayload(payload)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	update := bson.M{
		"name":      drink.Name,
		"price":     drink.Price,
		"stock":     drink.Stock,
		"available": drink.Available,
		"tags":      drink.Tags,
		"caffeine":  drink.Caffeine,
		"temp":      drink.Temp,
		"sweetness": drink.Sweetness,
		"colorTone": drink.ColorTone,
		"image":     drink.Image,
		"desc":      drink.Desc,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := db.DB.Collection("drinks")
	res, err := coll.UpdateOne(ctx, bson.M{"_id": drinkID}, bson.M{"$set": update})
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if res.MatchedCount == 0 {
		respondError(c, http.StatusNotFound, "drink not found")
		return
	}

	var updated models.Drink
	if err := coll.FindOne(ctx, bson.M{"_id": drinkID}).Decode(&updated); err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	respondMessage(c, http.StatusOK, "product updated", updated)
}

func AdminDeleteDrink(c *gin.Context) {
	drinkID, ok := parseHexID(c, "id")
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := db.DB.Collection("drinks").DeleteOne(ctx, bson.M{"_id": drinkID})
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if res.DeletedCount == 0 {
		respondError(c, http.StatusNotFound, "drink not found")
		return
	}

	respondMessage(c, http.StatusOK, "product deleted", nil)
}

func buildBookingFromPayload(payload adminBookingPayload) (models.Booking, error) {
	payload.Email = strings.TrimSpace(payload.Email)
	payload.Name = strings.TrimSpace(payload.Name)
	payload.Phone = strings.TrimSpace(payload.Phone)
	payload.Channel = strings.TrimSpace(payload.Channel)

	if payload.Email == "" || payload.Name == "" || payload.Phone == "" || payload.Time == "" {
		return models.Booking{}, errors.New("email, name, phone and time are required")
	}

	timeVal, err := time.Parse(time.RFC3339, payload.Time)
	if err != nil {
		return models.Booking{}, errors.New("invalid time format")
	}

	items, err := mapBookingItems(payload.Items)
	if err != nil {
		return models.Booking{}, err
	}

	channel := payload.Channel
	if channel == "" {
		channel = "admin"
	}

	return models.Booking{
		Email:   payload.Email,
		Name:    payload.Name,
		Phone:   payload.Phone,
		Time:    timeVal,
		Guests:  payload.Guests,
		Items:   items,
		Status:  models.NormalizeBookingStatus(payload.Status),
		Channel: channel,
	}, nil
}

func mapBookingItems(items []adminBookingItemPayload) ([]models.BookingItem, error) {
	mapped := make([]models.BookingItem, 0, len(items))
	for _, item := range items {
		drinkID, err := primitive.ObjectIDFromHex(strings.TrimSpace(item.DrinkID))
		if err != nil {
			return nil, errors.New("invalid drink ID format")
		}
		qty := item.Qty
		if qty <= 0 {
			qty = 1
		}

		options := item.Options
		if options == nil {
			options = map[string]any{}
		}

		mapped = append(mapped, models.BookingItem{
			DrinkID: drinkID,
			Qty:     qty,
			Options: options,
		})
	}
	return mapped, nil
}

func buildDrinkFromPayload(payload adminDrinkPayload) (models.Drink, error) {
	payload.Name = strings.TrimSpace(payload.Name)
	payload.Image = strings.TrimSpace(payload.Image)
	payload.Desc = strings.TrimSpace(payload.Desc)
	payload.Caffeine = strings.TrimSpace(payload.Caffeine)
	payload.Temp = strings.TrimSpace(payload.Temp)
	payload.ColorTone = strings.TrimSpace(payload.ColorTone)

	if payload.Name == "" {
		return models.Drink{}, errors.New("name is required")
	}
	if payload.Price < 0 {
		return models.Drink{}, errors.New("price must be zero or greater")
	}

	return models.Drink{
		Name:      payload.Name,
		Price:     payload.Price,
		Stock:     payload.Stock,
		Available: payload.Available || payload.Stock > 0,
		Tags:      normalizeTags(payload.Tags),
		Caffeine:  normalizeOrDefault(payload.Caffeine, "none"),
		Temp:      normalizeOrDefault(payload.Temp, "iced"),
		Sweetness: payload.Sweetness,
		ColorTone: normalizeOrDefault(payload.ColorTone, "neutral"),
		Image:     payload.Image,
		Desc:      payload.Desc,
	}, nil
}

func parseHexID(c *gin.Context, param string) (primitive.ObjectID, bool) {
	id, err := primitive.ObjectIDFromHex(c.Param(param))
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid id")
		return primitive.NilObjectID, false
	}
	return id, true
}

func normalizeUserRole(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "admin":
		return "admin"
	default:
		return "user"
	}
}

func normalizeTags(tags []string) []string {
	normalized := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		normalized = append(normalized, tag)
	}
	return normalized
}

func normalizeOrDefault(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func normalizeMongoWriteError(err error, fallback string) (int, string) {
	if mongo.IsDuplicateKeyError(err) {
		return http.StatusConflict, "name or email already exists"
	}
	return http.StatusInternalServerError, fallback
}

func optionsFindByCreatedAtDesc() *mongoOptionsFind {
	return &mongoOptionsFind{
		Sort: bson.D{{Key: "createdAt", Value: -1}},
	}
}

func optionsFindByTimeDesc() *mongoOptionsFind {
	return &mongoOptionsFind{
		Sort: bson.D{{Key: "time", Value: -1}},
	}
}

func optionsFindByNameAsc() *mongoOptionsFind {
	return &mongoOptionsFind{
		Sort: bson.D{{Key: "name", Value: 1}},
	}
}

type mongoOptionsFind struct {
	Sort any
}

func (o *mongoOptionsFind) toFindOptions() *options.FindOptions {
	if o == nil {
		return nil
	}
	findOptions := options.Find()
	findOptions.SetSort(o.Sort)
	return findOptions
}
