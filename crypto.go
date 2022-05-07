package main

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateJWTES384(user_id string) string {
	private_key_bytes, err := os.ReadFile("ec-private.pem")
	if err != nil {
		log.Fatal(err)
	}
	private_key, err := jwt.ParseECPrivateKeyFromPEM(private_key_bytes)
	if err != nil {
		log.Fatal(err)
	}
	iat := time.Now()
	exp := iat.Add(time.Duration(15) * time.Minute)
	unsigned_token := jwt.NewWithClaims(jwt.SigningMethodES384, jwt.MapClaims{
		"iat":            iat,
		"exp":            exp,
		"user_global_id": user_id,
	})

	token, err := unsigned_token.SignedString(private_key)
	if err != nil {
		log.Fatal(err)
	}
	return token
}
