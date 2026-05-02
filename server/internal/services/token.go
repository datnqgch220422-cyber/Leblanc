package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"leblanc/server/internal/models"
)

// GenerateVerificationToken creates a HMAC-SHA256 signed token that encodes the email and expiry.
func GenerateVerificationToken(email string) (token string, expiresAt time.Time) {
	expiresAt = time.Now().Add(verificationTTL())
	payload := email + "|" + strconv.FormatInt(expiresAt.Unix(), 10)
	sig := signPayload(payload)
	raw := payload + "|" + sig
	token = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(raw))
	return token, expiresAt
}

// VerifyToken validates signature and expiry, returning the embedded email.
func VerifyToken(token string) (string, error) {
	if token == "" {
		return "", errors.New("missing token")
	}
	decoded, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(token)
	if err != nil {
		return "", errors.New("invalid token encoding")
	}
	parts := strings.Split(string(decoded), "|")
	if len(parts) != 3 {
		return "", errors.New("invalid token structure")
	}
	email := parts[0]
	expUnix, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", errors.New("invalid expiry")
	}
	expiresAt := time.Unix(expUnix, 0)
	if time.Now().After(expiresAt) {
		return "", errors.New("token expired")
	}
	payload := parts[0] + "|" + parts[1]
	expectedSig := signPayload(payload)
	if !hmac.Equal([]byte(expectedSig), []byte(parts[2])) {
		return "", errors.New("invalid signature")
	}
	return email, nil
}

type RegistrationClaims struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
	Role         string `json:"role,omitempty"`
	Exp          int64  `json:"exp"`
}

// GenerateRegistrationToken signs registration claims (name/email/password hash/role) with expiry.
func GenerateRegistrationToken(claims RegistrationClaims) (token string, expiresAt time.Time, err error) {
	expiresAt = time.Now().Add(registrationTTL())
	claims.Exp = expiresAt.Unix()
	b, err := json.Marshal(claims)
	if err != nil {
		return "", time.Time{}, err
	}
	payload := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
	sig := signPayload(payload)
	raw := payload + "|" + sig
	token = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(raw))
	return token, expiresAt, nil
}

// VerifyRegistrationToken validates signature/expiry and returns claims.
func VerifyRegistrationToken(token string) (*RegistrationClaims, error) {
	if token == "" {
		return nil, errors.New("missing token")
	}
	decoded, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(token)
	if err != nil {
		return nil, errors.New("invalid token encoding")
	}
	parts := strings.Split(string(decoded), "|")
	if len(parts) != 2 {
		return nil, errors.New("invalid token structure")
	}
	payload := parts[0]
	sig := parts[1]
	expectedSig := signPayload(payload)
	if !hmac.Equal([]byte(expectedSig), []byte(sig)) {
		return nil, errors.New("invalid signature")
	}
	payloadBytes, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(payload)
	if err != nil {
		return nil, errors.New("invalid payload encoding")
	}
	var claims RegistrationClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("invalid claims")
	}
	if claims.Exp == 0 || time.Now().After(time.Unix(claims.Exp, 0)) {
		return nil, errors.New("token expired")
	}
	return &claims, nil
}

type SessionClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}

func GenerateSessionToken(user models.User) (token string, expiresAt time.Time, err error) {
	expiresAt = time.Now().Add(sessionTTL())
	claims := SessionClaims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		Role:   user.Role,
		Exp:    expiresAt.Unix(),
	}
	b, err := json.Marshal(claims)
	if err != nil {
		return "", time.Time{}, err
	}
	payload := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
	sig := signPayload(payload)
	raw := payload + "|" + sig
	token = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(raw))
	return token, expiresAt, nil
}

func VerifySessionToken(token string) (*SessionClaims, error) {
	if token == "" {
		return nil, errors.New("missing token")
	}
	decoded, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(token)
	if err != nil {
		return nil, errors.New("invalid token encoding")
	}
	parts := strings.Split(string(decoded), "|")
	if len(parts) != 2 {
		return nil, errors.New("invalid token structure")
	}
	payload := parts[0]
	sig := parts[1]
	expectedSig := signPayload(payload)
	if !hmac.Equal([]byte(expectedSig), []byte(sig)) {
		return nil, errors.New("invalid signature")
	}
	payloadBytes, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(payload)
	if err != nil {
		return nil, errors.New("invalid payload encoding")
	}
	var claims SessionClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("invalid claims")
	}
	if claims.UserID == "" || claims.Exp == 0 {
		return nil, errors.New("invalid claims")
	}
	if time.Now().After(time.Unix(claims.Exp, 0)) {
		return nil, errors.New("token expired")
	}
	return &claims, nil
}

func tokenSecret() []byte {
	secret := strings.TrimSpace(os.Getenv("TOKEN_SECRET"))
	if secret == "" {
		secret = "change-me-token-secret"
	}
	return []byte(secret)
}

func verificationTTL() time.Duration {
	return envDurationMinutes("VERIFICATION_TTL_MIN", 30)
}

func registrationTTL() time.Duration {
	return envDurationMinutes("REGISTRATION_TTL_MIN", 60)
}

func sessionTTL() time.Duration {
	if v := strings.TrimSpace(os.Getenv("SESSION_TTL_HOURS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return time.Duration(n) * time.Hour
		}
	}
	return envDurationMinutes("SESSION_TTL_MIN", 12*60)
}

func envDurationMinutes(key string, fallback int) time.Duration {
	minutes := fallback
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			minutes = n
		}
	}
	return time.Duration(minutes) * time.Minute
}

func signPayload(payload string) string {
	h := hmac.New(sha256.New, tokenSecret())
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}
