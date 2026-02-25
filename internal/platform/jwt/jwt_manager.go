package jwt

import (
	"log"
	"time"

	"github.com/danilobml/workstream/internal/platform/errs"
	"github.com/danilobml/workstream/internal/platform/models"
	"github.com/golang-jwt/jwt/v5"
)

type JwtManager struct {
	SecretKey []byte
}

type Claims struct {
	Email string
	Roles []models.Role
	jwt.RegisteredClaims
}

type ResetClaims struct {
	Sub string `json:"sub"` // User.Id
	Exp int64  `json:"exp"`
}

const resetTTL = 15 * time.Minute

func NewJwtManager(secretKey []byte) *JwtManager {
	return &JwtManager{
		SecretKey: secretKey,
	}
}

func (j *JwtManager) CreateToken(email string, roles []models.Role) (string, error) {
	claims := Claims{
		Email: email,
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(j.SecretKey)
}

func (j *JwtManager) ParseAndValidateToken(tokenString string) (*Claims, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	token, err := parser.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		return j.SecretKey, nil
	})
	if err != nil {
		log.Println("error parsing token: ", err.Error())
		return nil, errs.ErrParsingToken
	}

	if !token.Valid {
		return nil, errs.ErrInvalidToken
	}

	return token.Claims.(*Claims), nil
}

func (m *JwtManager) CreateResetToken(userID string, userEmail string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"email": userEmail,
		"exp": time.Now().Add(resetTTL).Unix(),
		"prp": "reset",
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(m.SecretKey)
}

func (m *JwtManager) VerifyResetToken(tokenStr string) (string, string, error) {
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return m.SecretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil || !tok.Valid {
		if err != nil {
			log.Println("error parsing token:", err.Error())
		}
		return "", "", errs.ErrInvalidToken
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errs.ErrInvalidToken
	}

	if claims["prp"] != "reset" {
		return "", "", errs.ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", "", errs.ErrInvalidToken
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return "", "", errs.ErrInvalidToken
	}

	expF, ok := claims["exp"].(float64)
	if !ok {
		return "", "", errs.ErrInvalidToken
	}
	if time.Now().Unix() > int64(expF) {
		return "", "", errs.ErrInvalidToken
	}

	return sub, email, nil
}