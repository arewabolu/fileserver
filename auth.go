package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type JWTClaims struct {
	Email string
	jwt.RegisteredClaims
}

func passwordHash(password string) (hashedPassword []byte, err error) {
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}
	return
}

func validateHash(givenPasswd, dbPasswd string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(givenPasswd), []byte(dbPasswd))
	return
}

func keyRetriever() []byte {
	data, err := os.ReadFile("./trivialfiles.txt")
	if err != nil {
		return nil
	}
	return data
}

func GenerateJWT(email string) (tokenString string, err error) {
	expirationTime := time.Now().Add(2 * time.Hour)
	claims := &JWTClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(keyRetriever())
	return
}

func ValidateJWT(givenToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		givenToken,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(keyRetriever()), nil
		},
	)
	if err != nil {
		return
	}
	fmt.Println(token.Signature)
	Assert, ok := token.Claims.(*JWTClaims)
	if !ok {
		err = errors.New("unable to verify claims")
		return
	}
	if Assert.ExpiresAt.Unix() < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}
