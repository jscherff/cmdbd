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

// AuthTokenValidator is middleWare that validates a client authentication
// token prior to allowing access to protected pages.
func AuthTokenHandler(authSvc service.AuthSvc, next http.Handler) http.Handler {

	return http.HandlerFunc(

		func(w http.ResponseWriter, r *http.Request) {

			if tokenString, err := authSvc.ReadCookie(r); err != nil {

				http.Error(w, err.Error(), http.StatusUnauthorized)
				panic(err)

			} else if token, err := authSvc.ParseTokenString(tokenString); err != nil {

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

// ConnectionLimiter limits the number of concurrent connections using a
// counting semephore modeled with a buffered channel. At maxConnections
// active clients, new requests queue until a 'slot' becomes available.
// From https://pauladamsmith.com/blog/2016/04/max-clients-go-net-http.html.
func MaxConnectionHandler(maxConnections int, next http.Handler) http.Handler {

	sema := make(chan struct{}, maxConnections)

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			sema <- struct{}{}
			defer func() { <-sema }()
			next.ServeHTTP(w, r)
		},
	)
}
