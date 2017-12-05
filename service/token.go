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
	`io/ioutil`
	`time`
	jwt `github.com/dgrijalva/jwt-go`

	`github.com/jscherff/cmdbd/model/cmdb`
)

// AuthClaims is a custom Claims object that extends jwt.StandardClaims.
type AuthClaims struct {
	jwt.StandardClaims
	cmdb.User
}

// Token is a custom Token object that extends jwt.Token.
type Token struct {
	*jwt.Token
}

// Username extracts the Username AuthClaim claim from the token.
func (this *Token) User() (cmdb.User) {
	return this.Claims.(AuthClaims).User
}

// AuthTokenService is an interface that creates, parses, and validates Tokens.
type AuthTokenService interface {
	Create(user, role string) (token *Token)
	String(token *Token) (tokenString string, err error)
	Parse(tokenString string) (token *Token, err error)
	Valid(tokenString string) (ok bool)
}

// authTokenService is a service that implements the AuthTokenService interface.
type authTokenService struct {
	privateKey *rsa.PrivateKey
	publicKey *rsa.PublicKey
	maxAge time.Duration
}

// NewAuthTokenService returns an object that implements the AuthTokenService interface.
func NewAuthTokenService(keyFiles map[string]string, maxAge time.Duration) (AuthTokenService, error) {

	this := &authTokenService{maxAge: maxAge}
	priKeyName, pubKeyName := `PrivateRSA`, `PublicRSA`

	// Read and store RSA private key.

	if priKeyFile, ok := keyFiles[priKeyName]; !ok {
		return nil, fmt.Errorf(`key name %q not found`, priKeyName)
	} else if pemKey, err := ioutil.ReadFile(priKeyFile); err != nil {
		return nil, err
	} else if rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemKey); err != nil {
		return nil, err
	} else {
		this.privateKey = rsaKey
	}

	// Read and store RSA public key.

	if pubKeyFile, ok := keyFiles[pubKeyName]; !ok {
		return nil, fmt.Errorf(`key name %q not found`, pubKeyName)
	} else if pemKey, err := ioutil.ReadFile(pubKeyFile); err != nil {
		return nil, err
	} else if rsaKey, err := jwt.ParseRSAPublicKeyFromPEM(pemKey); err != nil {
		return nil, err
	} else {
		this.publicKey = rsaKey
	}

	return this, nil
}

// Create generates a new Token.
func (this *authTokenService) Create(username, role string) (*Token) {

	claims := &AuthClaims {

		StandardClaims: jwt.StandardClaims {
			IssuedAt: time.Now().Unix(),
			ExpiresAt: time.Now().Add(this.maxAge).Unix(),
		},

		User: cmdb.User {
			Username: username,
			Role: role,
		},
	}

	return &Token{jwt.NewWithClaims(jwt.GetSigningMethod(`RS256`), claims)}
}

// Parse parses a token string and returns a Token.
func (this *authTokenService) Parse(tokenString string) (*Token, error) {

	token, err := jwt.ParseWithClaims(

		tokenString, &AuthClaims{},

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
