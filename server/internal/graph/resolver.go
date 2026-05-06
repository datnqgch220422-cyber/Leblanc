package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"leblanc/server/internal/db"
	"leblanc/server/internal/models"
	"leblanc/server/internal/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Resolver struct{}

// Query resolvers
func (r *Resolver) Drinks(ctx context.Context) ([]*models.Drink, error) {
	cur, err := db.DB.Collection("drinks").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var drinks []*models.Drink
	if err := cur.All(ctx, &drinks); err != nil {
		return nil, err
	}
	return drinks, nil
}

func (r *Resolver) Drink(ctx context.Context, id string) (*models.Drink, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format")
	}

	var drink models.Drink
	err = db.DB.Collection("drinks").FindOne(ctx, bson.M{"_id": objID}).Decode(&drink)
	if err != nil {
		return nil, err
	}
	return &drink, nil
}

func (r *Resolver) Users(ctx context.Context) ([]*models.User, error) {
	cur, err := db.DB.Collection("users").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []*models.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Resolver) Bookings(ctx context.Context) ([]*models.Booking, error) {
	cur, err := db.DB.Collection("bookings").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var bookings []*models.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// Mutation resolvers
type CreateBookingInput struct {
	Email   string             `json:"email"`
	Name    string             `json:"name"`
	Phone   string             `json:"phone"`
	Time    string             `json:"time"`
	Guests  int                `json:"guests"`
	Items   []BookingItemInput `json:"items"`
	Channel string             `json:"channel"`
}

type BookingItemInput struct {
	DrinkID string `json:"drinkId"`
	Qty     int    `json:"qty"`
	Options string `json:"options"`
}

func (r *Resolver) CreateBooking(ctx context.Context, input CreateBookingInput) (*models.Booking, error) {
	timeVal, err := time.Parse(time.RFC3339, input.Time)
	if err != nil {
		return nil, fmt.Errorf("invalid time format")
	}
	email := strings.TrimSpace(input.Email)
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	items := make([]models.BookingItem, len(input.Items))
	for i, item := range input.Items {
		drinkID, err := primitive.ObjectIDFromHex(item.DrinkID)
		if err != nil {
			return nil, fmt.Errorf("invalid drink ID format")
		}

		var options map[string]any
		if item.Options != "" {
			if err := json.Unmarshal([]byte(item.Options), &options); err != nil {
				options = map[string]any{}
			}
		} else {
			options = map[string]any{}
		}

		items[i] = models.BookingItem{
			DrinkID: drinkID,
			Qty:     item.Qty,
			Options: options,
		}
	}

	booking := models.Booking{
		ID:      primitive.NewObjectID(),
		Email:   email,
		Name:    input.Name,
		Phone:   input.Phone,
		Time:    timeVal,
		Guests:  input.Guests,
		Items:   items,
		Channel: input.Channel,
	}

	_, err = db.DB.Collection("bookings").InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Ok   bool         `json:"ok"`
	User *models.User `json:"user"`
}

func (r *Resolver) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.TrimSpace(input.Email)
	password := strings.TrimSpace(input.Password)

	if name == "" || email == "" || password == "" {
		return nil, fmt.Errorf("name, email and password are required")
	}

	coll := db.DB.Collection("users")
	lowerName := strings.ToLower(name)
	lowerEmail := strings.ToLower(email)

	// Check if user exists
	filter := bson.M{"$or": []bson.M{
		{"nameLower": lowerName},
		{"emailLower": lowerEmail},
	}}
	if err := coll.FindOne(ctx, filter).Err(); err == nil {
		return nil, fmt.Errorf("user already exists")
	} else if err != mongo.ErrNoDocuments {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password")
	}

	user := models.User{
		ID:           primitive.NewObjectID(),
		Name:         name,
		NameLower:    lowerName,
		Email:        email,
		EmailLower:   lowerEmail,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	if _, err := coll.InsertOne(ctx, user); err != nil {
		return nil, err
	}

	return &AuthResponse{Ok: true, User: &user}, nil
}

func (r *Resolver) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	nameOrEmail := strings.TrimSpace(input.Name)
	password := strings.TrimSpace(input.Password)

	if nameOrEmail == "" || password == "" {
		return nil, fmt.Errorf("name and password are required")
	}

	coll := db.DB.Collection("users")
	lower := strings.ToLower(nameOrEmail)

	filter := bson.M{"$or": []bson.M{
		{"nameLower": lower},
		{"emailLower": lower},
	}}

	var user models.User
	if err := coll.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &AuthResponse{Ok: true, User: &user}, nil
}

type RecommendationScore struct {
	DrinkID string  `json:"drinkId"`
	Score   float64 `json:"score"`
}

func (r *Resolver) RecommendFromFeatures(ctx context.Context, caffeine *string, temp *string, sweetness *int) ([]*RecommendationScore, error) {
	// Get all drinks
	cur, err := db.DB.Collection("drinks").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var drinks []models.Drink
	if err := cur.All(ctx, &drinks); err != nil {
		return nil, err
	}

	var caffeineVal, tempVal string
	var sweetnessVal int

	if caffeine != nil {
		caffeineVal = *caffeine
	}
	if temp != nil {
		tempVal = *temp
	}
	if sweetness != nil {
		sweetnessVal = *sweetness
	}

	// Score each drink
	scores := services.ScoreDrinks(drinks, caffeineVal, tempVal, sweetnessVal)

	// Convert to response format
	result := make([]*RecommendationScore, len(scores))
	for i, score := range scores {
		result[i] = &RecommendationScore{
			DrinkID: score.DrinkID,
			Score:   score.Score,
		}
	}

	return result, nil
}
