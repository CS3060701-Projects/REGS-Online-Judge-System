package jwt

import (
	"os"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5" // 使用別名
)

var (
	SignKey   interface{}
	VerifyKey interface{}
)

func InitKeys() error {
	privBytes, _ := os.ReadFile("private.pem")
	var err error
	SignKey, err = jwtlib.ParseECPrivateKeyFromPEM(privBytes)
	if err != nil {
		return err
	}

	pubBytes, _ := os.ReadFile("public.pem")
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
