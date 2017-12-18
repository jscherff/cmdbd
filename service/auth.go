// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	`crypto/rsa`
	`fmt`
	`net/http`
	`time`
	jwt `github.com/dgrijalva/jwt-go`

	`github.com/jscherff/cmdbd/model/cmdb`
)

// AuthClaims is a custom claims object that contains user authentication
// and authorization infomration.
type AuthClaims struct {
	*cmdb.User
}

// Claims is a custom Claims object that extends jwt.StandardClaims.
type Claims struct {
	jwt.StandardClaims
	AuthClaims
}

// Token is a custom Token object that extends jwt.Token.
type Token struct {
	*jwt.Token
}

// AuthClaims extracts the AuthClaim claim from the token.
func (this *Token) AuthClaims() (AuthClaims) {
	return this.Claims.(Claims).AuthClaims
}

// AuthSvc is an interface that creates, parses, and validates Tokens.
type AuthSvc interface {
	CreateToken(user *cmdb.User) (token *Token, err error)
	CreateTokenString(token *Token) (tokenString string, err error)
	ParseTokenString(tokenString string) (token *Token, err error)
	CreateCookie(tokenString string) (cookie *http.Cookie, err error)
	ReadCookie(request *http.Request) (tokenString string, err error)
}

// authSvc is a service that implements the AuthSvc interface.
type authSvc struct {
	PrivateKey *rsa.PrivateKey
	PublicKey *rsa.PublicKey
	AuthMaxAge time.Duration
}

// NewAuthSvc returns an object that implements the AuthSvc interface.
func NewAuthSvc(pubKey, priKey []byte, authMaxAge time.Duration) (AuthSvc, error) {

	this := &authSvc{AuthMaxAge: authMaxAge}

	// Process RSA public key.

	if rsaKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKey); err != nil {
		return nil, err
	} else {
		this.PublicKey = rsaKey
	}

	// Process RSA private key.

	if rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(priKey); err != nil {
		return nil, err
	} else {
		this.PrivateKey = rsaKey
	}

	return this, nil
}

// CreateToken generates a new Token.
func (this *authSvc) CreateToken(user *cmdb.User) (*Token, error) {

	claims := &Claims {

		StandardClaims: jwt.StandardClaims {
			IssuedAt: time.Now().Unix(),
			ExpiresAt: time.Now().Add(this.AuthMaxAge).Unix(),
		},

		AuthClaims: AuthClaims{user},
	}

	return &Token{jwt.NewWithClaims(jwt.GetSigningMethod(`RS256`), claims)}, nil
}

// ParseTokenString parses a token string and returns a Token.
func (this *authSvc) ParseTokenString(tokenString string) (*Token, error) {

	token, err := jwt.ParseWithClaims(

		tokenString, &Claims{},

		func(t *jwt.Token) (interface{}, error) {

			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf(`unexpected signing method: %v`, t.Header[`alg`])
			}

			return this.PublicKey, nil
		},
	)

	if err != nil {
		return nil, err
	} else if token == nil {
		return nil, fmt.Errorf(`empty token`)
	} else if !token.Valid {
		return nil, fmt.Errorf(`invalid token`)
	}

	return &Token{token}, nil
}

// CreateTokenString returns a token string suitable for cookies.
func (this *authSvc) CreateTokenString(token *Token) (string, error) {
	return token.SignedString(this.PrivateKey)
}


// CreateCookie generates a new authentication http.Cookie from an auth token string.
func (this *authSvc) CreateCookie(tokenString string) (*http.Cookie, error) {

	return &http.Cookie{
		Name: `Auth`,
		Value: tokenString,
		Expires: time.Now().Add(this.AuthMaxAge),
		HttpOnly: true,
	}, nil
}

// ReadCookie extracts the 'Auth' http.Cookie from an http.Request.
func (this *authSvc) ReadCookie(request *http.Request) (string, error) {

	if cookie, err := request.Cookie(`Auth`); err != nil {
		return ``, err
	} else {
		return cookie.Value, nil
	}
}
