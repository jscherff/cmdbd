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

package main

import (
	`net/http`
	`time`
	jwt `github.com/dgrijalva/jwt-go`
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func cmdbAuthSetTokenV1(w http.ResponseWriter, r *http.Request) {

	expireTime := time.Now().Add(time.Hour * 1)

	claims := Claims {
		Username: r.URL.User.Username(),
		StandardClaims: jwt.StandardClaims {
			Issuer: r.URL.Host,
			ExpiresAt: expireTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(`RS256`), claims)

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keys[`Private`])
	if err != nil {
		el.Panic(err)
	}
/*
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keys[`Public`])
	if err != nil {
		el.Panic(err)
	}
*/
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		el.Panic(err)
	}

	cookie := http.Cookie{
		Name: `Auth`,
		Value: signedToken,
		Expires: expireTime,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
}
