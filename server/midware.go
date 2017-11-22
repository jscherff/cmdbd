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
)

// AuthTokenValidator is middleware that can be used to validate
// a client authentication JWT prior to allowing access to protected
// pages.
func AuthTokenValidator(next http.Handler) (http.Handler) {

	return http.HandlerFunc(

		func(w http.ResponseWriter, r *http.Request) {

			fail := func(err error) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				panic(err)
			}

			if cookie, err := getAuthCookie(r); err != nil {
				fail(err)
			} else if token, err := parseAuthToken(cookie.Value); err != nil {
				fail(err)
			} else if err := validateAuthToken(token); err != nil {
				fail(err)
			} else if claims, ok := token.Claims.(*Claims); ok {
				ctx := context.WithValue(r.Context(), `Claims`, *claims)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				next.ServeHTTP(w, r)
			}
		},
	)
}
