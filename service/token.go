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
	`time`
	jwt `github.com/dgrijalva/jwt-go`
)

// Claims is a custom Claims object that extends jwt.StandardClaims.
type Claims struct {
	Username string `json:"username"`
	Role string `json:"role"`
	jwt.StandardClaims
}

// AuthTokenService is an interface that creates, parses, and validates jwt.Tokens.
type AuthTokenService interface {
	Create(user, pass, host string) (*jwt.Token, error)
	Parse(tokenString string) (*jwt.Token, error)
	Validate(token *jwt.Token) (error)
}

// authTokenService is a service that implements the AuthTokenService interface.
type authTokenService struct {
	AuthTokenService
	PrivateKey *rsa.PrivateKey
	PublicKey *rsa.PublicKey
}

// NewAuthTokenService returns an object that implements the AuthTokenService interface.
func NewAuthTokenService(priKey, pubKey []byte) (*authTokenService, error) {

	this := &authTokenService{}

	if rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(priKey); err != nil {
		return nil, err
	} else {
		this.PrivateKey = rsaKey
	}

	if rsaKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKey); err != nil {
		return nil, err
	} else {
		this.PublicKey = rsaKey
	}

	return this, nil
}

// Create generates a new jwt.Token.
func (this *authTokenService) Create(user, host string) (*jwt.Token, error) {

	claims := &Claims {
		Username: user,
		StandardClaims: jwt.StandardClaims {
			Issuer: host,
			IssuedAt: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	return jwt.NewWithClaims(jwt.GetSigningMethod(`RS256`), claims), nil
}

// Parse parses a token string and returns a jwt.Token.
func (this *authTokenService) Parse(tokenString string) (*jwt.Token, error) {

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
	}

	return token, nil
}

// Validate validates a jwt.Token.
func (this *authTokenService) Validate(token *jwt.Token) (error) {

	if token == nil || !token.Valid {
		return fmt.Errorf(`nil or invalid auth token`)
	}

	return nil
}
