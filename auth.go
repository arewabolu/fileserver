package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func passwordHash(password string) (hashedPassword []byte, err error) {
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}
	return
}

func validateHash(givenPasswd, dbPasswd string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(dbPasswd), []byte(givenPasswd))
	return
}

func keyRetriever() []byte {
	data, err := os.ReadFile("./trivialfiles.txt")
	if err != nil {
		return nil
	}
	return data
}

func GenerateJWT(ID int, email string) (tokenString string, err error) {
	expirationTime := time.Now().Add(2 * time.Hour)
	claims := JWTClaims{
		Email:  email,
		UserId: strconv.Itoa(ID),
		RegisteredClaims: jwt.RegisteredClaims{
			//Issuer:    strconv.Itoa(ID),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(keyRetriever())
	return
}

func ValidateJWT(givenToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		givenToken,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return keyRetriever(), nil
		},
	)
	if err != nil {
		return nil, err
	}
	//Converts if token.Claims is of type JWY
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		err = errors.New("unable to verify claims")
		return nil, err
	}
	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return nil, err
	}

	return claims, err
}

func middlewareAUth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		parts := strings.Split(header, "Bearer")
		token := strings.TrimSpace(parts[2])

		claims, validationErr := ValidateJWT(token)
		if validationErr != nil {
			w.Write([]byte(validationErr.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//	fmt.Println(claims.ID)
		ctx := context.WithValue(r.Context(), uuid, claims.UserId)

		//If the method is post and the uri has /user/userid/ then check if userid is equal to issuerid
		//if the method is get return a users file list in bucket
		//if the method is get with download then check if userid is equal to issuerid
		//fmt.Println(uri[2]) // get item 2 from split uri which should be id

		next.ServeHTTP(w, r.WithContext(ctx))
	},
	)
}
