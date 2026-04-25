package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *jwt.Token
	publicKey  *jwt.Token
	signKey    interface{}
	verifyKey  interface{}
)

func InitKeys() error {
	privBytes, err := os.ReadFile("private.pem")
	if err != nil {
		return fmt.Errorf("無法讀取 private.pem: %w", err)
	}
	signKey, err = jwt.ParseECPrivateKeyFromPEM(privBytes)
	if err != nil {
		return fmt.Errorf("解析私鑰失敗: %w", err)
	}

	pubBytes, err := os.ReadFile("public.pem")
	if err != nil {
		return fmt.Errorf("無法讀取 public.pem: %w", err)
	}
	verifyKey, err = jwt.ParseECPublicKeyFromPEM(pubBytes)
	if err != nil {
		return fmt.Errorf("解析公鑰失敗: %w", err)
	}

	return nil
}

func GenerateToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(signKey)
}
