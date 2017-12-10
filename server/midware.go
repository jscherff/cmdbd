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

package server

import (
	`context`
	`net/http`

	`github.com/jscherff/cmdbd/service`
)

// MiddleWare provides various middleWare handlers used by the HTTP Server.
type MiddleWare interface {
	AuthTokenValidator(next http.Handler) (http.Handler)
}

// middleWare is an implementation of the MiddleWare interface.
type middleWare struct {
	AuthTokenSvc service.AuthTokenService
	AuthCookieSvc service.AuthCookieService
}

// NewMiddleWare implements a new MiddleWare instance.
func NewMiddleWare(conf *Config) (MiddleWare, error) {

	this := &middleWare{
		AuthTokenSvc: confAuthTokenSvc,
		AuthCookieSvc: confAuthCookieSvc,
	}

	return this, nil
}

// AuthTokenValidator is middleWare that validates a client authentication
// token prior to allowing access to protected pages.
func (this *middleWare) AuthTokenValidator(next http.Handler) (http.Handler) {

	return http.HandlerFunc(

		func(w http.ResponseWriter, r *http.Request) {

			if tokenString, err := this.AuthCookieSvc.Read(r); err != nil {

				http.Error(w, err.Error(), http.StatusUnauthorized)
				panic(err)

			} else if token, err := this.AuthTokenSvc.Parse(tokenString); err != nil {

				http.Error(w, err.Error(), http.StatusUnauthorized)
				panic(err)

			} else {

				claims := token.AuthClaims()
				context := context.WithValue(r.Context(), `AuthClaims`, claims)
				next.ServeHTTP(w, r.WithContext(context))
			}
		},
	)
}
