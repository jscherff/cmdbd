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

package legacy

import (
	`fmt`
	`net/http`
	`time`
	`golang.org/x/crypto/bcrypt`
	jwt `github.com/dgrijalva/jwt-go`
)

type Claims struct {
	Username string `json:"username"`
	Role string `json:"role"`
	jwt.StandardClaims
}

func createAuthToken(user, pass, host string) (*jwt.Token, error) {

	if hash, err := GetUserPassword(user); err != nil {
		return nil, fmt.Errorf(`invalid username`)
	} else if err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)); err != nil {
		return nil, fmt.Errorf(`invalid password`)
	}

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

func createAuthCookie(token *jwt.Token) (*http.Cookie, error) {

	tokenString, err := token.SignedString(privateKey)

	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name: `Auth`,
		Value: tokenString,
		Expires: time.Now().Add(time.Hour * 1),
		HttpOnly: true,
	}

	return cookie, nil
}

func getAuthCookie(r *http.Request) (*http.Cookie, error) {

	if cookie, err := r.Cookie(`Auth`); err != nil {
		return nil, err
	} else {
		return cookie, nil
	}
}

func parseAuthToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(

		tokenString, &Claims{},

		func(t *jwt.Token) (interface{}, error) {

			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf(`unexpected signing method: %v`, t.Header[`alg`])
			}

			return publicKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func validateAuthToken (token *jwt.Token) (error) {

	if token == nil || !token.Valid {
		return fmt.Errorf(`nil or invalid auth token`)
	}

	return nil
}
