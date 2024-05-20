package jwtHelper

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/naman2607/netflixClone/config"
)

type CustomClaim struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type TokenPayload struct {
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
	Issuer    string `json:"iss"`
}

func CreateNewToken(claim jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	conf := config.GetInstance()
	secret := []byte(conf.GetSecretKey())
	tokenString, _ := token.SignedString(secret)
	return tokenString
}

func ValidateToken(token string) bool {
	payload := strings.Split(token, ".")
	data, err := base64.StdEncoding.DecodeString(payload[1])
	if err != nil {
		log.Fatal("error in decoding token payload :", err)
	}
	var jwtPayload TokenPayload
	json.Unmarshal(data, &jwtPayload)
	claim := CustomClaim{
		Email: jwtPayload.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwtPayload.ExpiresAt,
			Issuer:    "netflix",
		},
	}
	secretKey, err := jwt.ParseWithClaims(token, claim, func(t *jwt.Token) (interface{}, error) {
		return []byte("hello"), nil
	})
	if err != nil {
		log.Println("error in validate token ", err)
	}
	log.Println("secret key", secretKey)
	return true
}
