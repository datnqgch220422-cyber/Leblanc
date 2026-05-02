package handlers

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	"leblanc/server/internal/db"
	"leblanc/server/internal/models"
	"leblanc/server/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	NameOrEmail string `json:"name"`
	Password    string `json:"password"`
}

type verifyRequest struct {
	Email string `json:"email"`
}

type verifyTokenRequest struct {
	Token string `json:"token"`
}

func GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coll := db.DB.Collection("users")

	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)

	var users []models.User
	if err := cur.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	publicUsers := make([]models.PublicUser, len(users))
	for i, u := range users {
		publicUsers[i] = u.Public()
	}

	c.JSON(http.StatusOK, publicUsers)
}

func RegisterUser(c *gin.Context) {
	var req registerRequest
	if err := c.BindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Name == "" || req.Email == "" || req.Password == "" {
		respondError(c, http.StatusBadRequest, "name, email and password are required")
		return
	}
	if !isValidEmail(req.Email) {
		respondError(c, http.StatusBadRequest, "invalid email")
		return
	}
	// MX lookup: if no MX -> invalid; if lookup errors -> log and allow (to avoid blocking due to DNS issues).
	if emailRequiresMXCheck() {
		if ok := hasMXRecord(req.Email); !ok {
			respondError(c, http.StatusBadRequest, "invalid email")
			return
		}
	}
	if reservedAdminEmailLower() != "" && strings.ToLower(req.Email) == reservedAdminEmailLower() {
		respondError(c, http.StatusBadRequest, "this email is reserved")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coll := db.DB.Collection("users")

	lowerName := strings.ToLower(req.Name)
	lowerEmail := strings.ToLower(req.Email)

	// ensure user doesn't already exist
	filter := bson.M{"$or": []bson.M{
		{"nameLower": lowerName},
		{"emailLower": lowerEmail},
	}}
	if err := coll.FindOne(ctx, filter).Err(); err == nil {
		respondError(c, http.StatusConflict, "user already exists")
		return
	} else if err != mongo.ErrNoDocuments {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not hash password")
		return
	}

	// Insert a pending (unverified) user immediately so you can see it in Mongo,
	// then return the verification token/URL to complete activation.
	user := models.User{
		ID:           primitive.NewObjectID(),
		Name:         req.Name,
		NameLower:    lowerName,
		Email:        req.Email,
		EmailLower:   lowerEmail,
		Role:         "user",
		Verified:     false,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}
	if _, err := coll.InsertOne(ctx, user); err != nil {
		respondError(c, http.StatusInternalServerError, "could not create user")
		return
	}

	claims := services.RegistrationClaims{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         "user",
	}
	token, expiresAt, err := services.GenerateRegistrationToken(claims)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not generate token")
		return
	}
	verifyURLBase := strings.TrimSpace(os.Getenv("FRONTEND_VERIFY_URL"))
	verifyURL := ""
	if verifyURLBase != "" {
		verifyURL = fmt.Sprintf("%s?token=%s", strings.TrimRight(verifyURLBase, "/"), token)
	}
	log.Printf("register-init: email=%s verifyUrlGenerated=%t", req.Email, verifyURL != "")

	c.JSON(http.StatusCreated, gin.H{
		"ok":        true,
		"token":     token,
		"expiresAt": expiresAt,
		"verifyUrl": verifyURL,
		"user":      user.Public(),
	})
}

func LoginUser(c *gin.Context) {
	var req loginRequest
	if err := c.BindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	req.NameOrEmail = strings.TrimSpace(req.NameOrEmail)
	req.Password = strings.TrimSpace(req.Password)
	if req.NameOrEmail == "" || req.Password == "" {
		respondError(c, http.StatusBadRequest, "name and password are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coll := db.DB.Collection("users")

	lookup := strings.ToLower(req.NameOrEmail)
	filter := bson.M{
		"$or": []bson.M{
			{"nameLower": lookup},
			{"emailLower": lookup},
		},
	}

	var user models.User
	if err := coll.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			respondError(c, http.StatusUnauthorized, "user not found")
		} else {
			respondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		respondError(c, http.StatusUnauthorized, "incorrect password")
		return
	}

	if !user.Verified {
		respondError(c, http.StatusUnauthorized, "email not verified")
		return
	}

	token, expiresAt, err := services.GenerateSessionToken(user)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "could not create session")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":        true,
		"user":      user.Public(),
		"token":     token,
		"expiresAt": expiresAt,
	})
}

// Request verification token (e.g., to send via email).
func RequestVerify(c *gin.Context) {
	var req verifyRequest
	if err := c.BindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		respondError(c, http.StatusBadRequest, "email is required")
		return
	}
	token, expiresAt := services.GenerateVerificationToken(req.Email)
	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"expiresAt": expiresAt,
	})
}

// Verify token and return embedded email.
func VerifyToken(c *gin.Context) {
	var req verifyTokenRequest
	if err := c.BindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	// Try registration token first (contains user data).
	if claims, err := services.VerifyRegistrationToken(req.Token); err == nil && claims != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		coll := db.DB.Collection("users")
		lowerEmail := strings.ToLower(claims.Email)
		var existing models.User
		err := coll.FindOne(ctx, bson.M{"emailLower": lowerEmail}).Decode(&existing)
		if err == mongo.ErrNoDocuments {
			user := models.User{
				ID:           primitive.NewObjectID(),
				Name:         claims.Name,
				NameLower:    strings.ToLower(claims.Name),
				Email:        claims.Email,
				EmailLower:   lowerEmail,
				Role:         claims.Role,
				Verified:     true,
				PasswordHash: claims.PasswordHash,
				CreatedAt:    time.Now(),
			}
			if user.Role == "" {
				user.Role = "user"
			}
			if _, err := coll.InsertOne(ctx, user); err != nil {
				respondError(c, http.StatusInternalServerError, err.Error())
				return
			}
			c.JSON(http.StatusOK, gin.H{"ok": true, "email": claims.Email, "user": user.Public()})
			return
		} else if err == nil {
			update := bson.M{"verified": true}
			if claims.PasswordHash != "" {
				update["passwordHash"] = claims.PasswordHash
			}
			if claims.Name != "" {
				update["name"] = claims.Name
				update["nameLower"] = strings.ToLower(claims.Name)
			}
			if claims.Role != "" {
				update["role"] = claims.Role
			}
			_, _ = coll.UpdateOne(ctx, bson.M{"_id": existing.ID}, bson.M{"$set": update})
			// Refresh user fields.
			updated := existing
			if claims.Name != "" {
				updated.Name = claims.Name
				updated.NameLower = strings.ToLower(claims.Name)
			}
			if claims.Role != "" {
				updated.Role = claims.Role
			}
			updated.Verified = true
			if claims.PasswordHash != "" {
				updated.PasswordHash = claims.PasswordHash
			}
			c.JSON(http.StatusOK, gin.H{"ok": true, "email": claims.Email, "user": updated.Public()})
			return
		} else {
			respondError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Fallback: legacy simple email token.
	email, err := services.VerifyToken(req.Token)
	if err != nil {
		respondError(c, http.StatusUnauthorized, err.Error())
		return
	}
	filter := bson.M{"emailLower": strings.ToLower(email)}
	_, _ = db.DB.Collection("users").UpdateOne(context.Background(), filter, bson.M{"$set": bson.M{"verified": true}})
	var user models.User
	_ = db.DB.Collection("users").FindOne(context.Background(), filter).Decode(&user)

	c.JSON(http.StatusOK, gin.H{"ok": true, "email": email, "user": user.Public()})
}

func isValidEmail(addr string) bool {
	_, err := mail.ParseAddress(addr)
	return err == nil
}

func hasMXRecord(addr string) bool {
	parts := strings.Split(addr, "@")
	if len(parts) != 2 {
		return false
	}
	mx, err := net.LookupMX(parts[1])
	if err != nil {
		return false
	}
	return len(mx) > 0
}

func reservedAdminEmailLower() string {
	return strings.ToLower(strings.TrimSpace(os.Getenv("ADMIN_EMAIL")))
}

func emailRequiresMXCheck() bool {
	return !strings.EqualFold(strings.TrimSpace(os.Getenv("EMAIL_REQUIRE_MX")), "false")
}
