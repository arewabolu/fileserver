package main

import "github.com/golang-jwt/jwt/v4"

const (
	bucketName = "localhost-fileserver"
	region     = "eu-west-2"
)

// uuid provides a way to pass a users id
// through the request context.
var uuid = struct{}{}

type JWTClaims struct {
	Email  string
	UserId string
	jwt.RegisteredClaims
}

func DerefString(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

// usrPath creates a user path for s3
// bucket using the userId
func usrPath(userId string) string {
	return "user" + userId + "files/"
}
