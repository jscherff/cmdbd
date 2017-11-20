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
	jwt.StandardClaims
	CustomClaims map[string]string
}

// Token is a custom Token object that extends jwt.Token.
type Token struct {
	*jwt.Token
}

// Claim extracts the named claim from the token.
func (this *Token) Claim(name string) (string, error) {
	if value, ok := this.Claims.(Claims).CustomClaims[name]; !ok {
		return ``, fmt.Errorf(`claim %q does not exist`)
	} else {
		return value, nil
	}
}

// AuthTokenService is an interface that creates, parses, and validates Tokens.
type AuthTokenService interface {
	Create(map[string]string, time.Duration) (string, error)
	String(*Token) (string, error)
	Parse(string) (*Token, error)
	Valid(string) bool
}

// authTokenService is a service that implements the AuthTokenService interface.
type authTokenService struct {
	privateKey *rsa.PrivateKey
	publicKey *rsa.PublicKey
}

// NewAuthTokenService returns an object that implements the AuthTokenService interface.
func NewAuthTokenService(priKey, pubKey []byte) (AuthTokenService, error) {

	this := &authTokenService{}

	if rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(priKey); err != nil {
		return nil, err
	} else {
		this.privateKey = rsaKey
	}

	if rsaKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKey); err != nil {
		return nil, err
	} else {
		this.publicKey = rsaKey
	}

	return this, nil
}

// Create generates a new Token.
func (this *authTokenService) Create(custClaims map[string]string, maxAge time.Duration) (string, error) {

	claims := &Claims {

		StandardClaims: jwt.StandardClaims {
			IssuedAt: time.Now().Unix(),
			ExpiresAt: time.Now().Add(maxAge).Unix(),
		},

		CustomClaims: custClaims,
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(`RS256`), claims)
	return token.SignedString(this.privateKey)
}

// Parse parses a token string and returns a Token.
func (this *authTokenService) Parse(tokenString string) (*Token, error) {

	token, err := jwt.ParseWithClaims(

		tokenString, &Claims{},

		func(t *jwt.Token) (interface{}, error) {

			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf(`unexpected signing method: %v`, t.Header[`alg`])
			}

			return this.publicKey, nil
		},
	)

	if err != nil {
		return nil, err
	} else if token == nil {
		return nil, fmt.Errorf(`nil token`)
	}

	return &Token{token}, nil
}

// Valid validates a Token.
func (this *authTokenService) Valid(tokenString string) (bool) {

	if token, err := this.Parse(tokenString); err != nil {
		return false
	} else if !token.Valid {
		return false
	}

	return true
}

// String returns a token string suitable for cookies.
func (this *authTokenService) String(token *Token) (string, error) {
	return token.SignedString(this.privateKey)
}
