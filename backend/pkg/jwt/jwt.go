package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	jwtlib "github.com/golang-jwt/jwt/v5" // 使用別名
)

var (
	SignKey   interface{}
	VerifyKey interface{}
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func InitKeys() error {
	privBytes, err := os.ReadFile("private.pem")
	if err != nil {
		return err
	}
	SignKey, err = jwtlib.ParseECPrivateKeyFromPEM(privBytes)
	if err != nil {
		return err
	}

	pubBytes, err := os.ReadFile("public.pem")
	if err != nil {
		return err
	}
	VerifyKey, err = jwtlib.ParseECPublicKeyFromPEM(pubBytes)
	return err
}

func GenerateToken(userID uint, role string) (string, error) {
	claims := jwtlib.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodES256, claims)
	return token.SignedString(SignKey)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return VerifyKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
