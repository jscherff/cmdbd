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
	`net/http`
	`time`
)

// AuthCookieSvc is an interface that create and extracts authentication cookies.
type AuthCookieSvc interface {
	Create(tokenString string) (cookie *http.Cookie, err error)
	Read(request *http.Request) (tokenString string, err error)
}

// authCookieSvc is a service that implements the AuthCookieSvc interface.
type authCookieSvc struct {
	MaxAge time.Duration
}

// NewAuthCookieSvc returns an object that implements the AuthCookieSvc interface.
func NewAuthCookieSvc(maxAge time.Duration) (AuthCookieSvc, error) {
	return &authCookieSvc{MaxAge: maxAge}, nil
}

// Create generates a new authentication http.Cookie from an auth token string.
func (this *authCookieSvc) Create(tokenString string) (*http.Cookie, error) {

	return &http.Cookie{
		Name: `Auth`,
		Value: tokenString,
		Expires: time.Now().Add(this.MaxAge),
		HttpOnly: true,
	}, nil
}

// Read extracts the 'Auth' http.Cookie from an http.Request.
func (this *authCookieSvc) Read(request *http.Request) (string, error) {

	if cookie, err := request.Cookie(`Auth`); err != nil {
		return ``, err
	} else {
		return cookie.Value, nil
	}
}
