package server

import (
	`net/http`
	jwt `github.com/dgrijalva/jwt-go`
)

type Claims struct {
	jwt.StandardClaims
	CustomClaims map[string]string
}

func getAuthCookie(request *http.Request) (*http.Cookie, error) {
	return nil, nil
}

func parseAuthToken(cookieValue string) (*jwt.Token, error) {
	return nil, nil
}

func validateAuthToken(token *jwt.Token) (error) {
	return nil
}
