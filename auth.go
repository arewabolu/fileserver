package main

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

//type JWTClaims struct {
//	Email  string
//	UserId int
//	jwt.RegisteredClaims
//}

//type adapter func(http.Handler) http.Handler

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

func GenerateJWT(ID int) (tokenString string, err error) {
	expirationTime := time.Now().Add(2 * time.Hour)
	claims := jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(ID),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(keyRetriever())
	return
}

func ValidateJWT(givenToken string) (string, error) {
	token, err := jwt.ParseWithClaims(
		givenToken,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return keyRetriever(), nil
		},
	)
	if err != nil {
		return "", err
	}
	//Converts if token.Claims is of type JWY
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		err = errors.New("unable to verify claims")
		return "", err
	}
	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return "", err
	}
	return claims.Issuer, err
}

func middlewareAUth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		parts := strings.Split(header, "Bearer")
		token := strings.TrimSpace(parts[2])

		issuerID, validationErr := ValidateJWT(token)
		if validationErr != nil {
			w.Write([]byte(validationErr.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		uri := strings.SplitN(r.RequestURI, "/", 4)

		if issuerID != uri[2] {
			w.WriteHeader(http.StatusUnauthorized)
			http.NotFound(w, r)
			return
		}
		//If the method is post and the uri has /user/userid/ then check if userid is equal to issuerid
		//if the method is get return a users file list in bucket
		//if the method is get with download then check if userid is equal to issuerid
		//fmt.Println(uri[2]) // get item 2 from split uri which should be id

		next.ServeHTTP(w, r)
	},
	)
}

// splits the request uri string
// and appends it to be S3 ready.
func reqSplit(r *http.Request) string {
	uri := strings.SplitN(r.RequestURI, "/", 4)
	return strAppender(uri[2])
}

// appends strings for S3 call
func strAppender(s string) string {
	return "user" + s + "files"
}
