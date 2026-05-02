package handlers

import (
	"context"
	"errors"
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
)

const authUserContextKey = "authUser"

var loadUserByID = func(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	if err := db.DB.Collection("users").FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func RequireAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c.GetHeader("Authorization"))
		if token == "" {
			respondError(c, http.StatusUnauthorized, "missing bearer token")
			c.Abort()
			return
		}

		claims, err := services.VerifySessionToken(token)
		if err != nil {
			respondError(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		userID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			respondError(c, http.StatusUnauthorized, "invalid token subject")
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		user, err := loadUserByID(ctx, userID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				respondError(c, http.StatusUnauthorized, "account not found")
				c.Abort()
				return
			}
			respondError(c, http.StatusInternalServerError, "could not validate account")
			c.Abort()
			return
		}

		if !user.Verified {
			respondError(c, http.StatusForbidden, "verified account required")
			c.Abort()
			return
		}

		if normalizeUserRole(user.Role) != "admin" {
			respondError(c, http.StatusForbidden, "admin access required")
			c.Abort()
			return
		}

		c.Set(authUserContextKey, user.Public())
		c.Next()
	}
}

func extractBearerToken(header string) string {
	parts := strings.Fields(strings.TrimSpace(header))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return parts[1]
}
