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

// AuthClaims is a custom claims object that contains user authentication
// and authorization infomration.
type AuthClaims struct {
	cmdb.User
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

// Username extracts the Username AuthClaim claim from the token.
func (this *Token) AuthClaims() (AuthClaims) {
	return this.Claims.(Claims).AuthClaims
}

// AuthTokenService is an interface that creates, parses, and validates Tokens.
type AuthTokenService interface {
	Create(user cmdb.User) (token *Token)
	String(token *Token) (tokenString string, err error)
	Parse(tokenString string) (token *Token, err error)
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
func (this *authTokenService) Create(user cmdb.User) (*Token) {

	claims := &Claims {

		StandardClaims: jwt.StandardClaims {
			IssuedAt: time.Now().Unix(),
			ExpiresAt: time.Now().Add(this.maxAge).Unix(),
		},

		AuthClaims: AuthClaims{user},
	}

	return &Token{jwt.NewWithClaims(jwt.GetSigningMethod(`RS256`), claims)}
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
		return nil, fmt.Errorf(`empty token`)
	} else if !token.Valid {
		return nil, fmt.Errorf(`invalid token`)
	}

	return &Token{token}, nil
}

// String returns a token string suitable for cookies.
func (this *authTokenService) String(token *Token) (string, error) {
	return token.SignedString(this.privateKey)
}
