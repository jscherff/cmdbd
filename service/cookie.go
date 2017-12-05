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

// AuthCookieService is an interface that create and extracts authentication cookies.
type AuthCookieService interface {
	Create(tokenString string) (cookie *http.Cookie)
	Read(request *http.Request) (tokenString string, err error)
}

// authCookieService is a service that implements the AuthCookieService interface.
type authCookieService struct {
	maxAge time.Duration
}

// NewAuthCookieService returns an object that implements the AuthCookieService interface.
func NewAuthCookieService(maxAge time.Duration) (AuthCookieService, error) {
	return &authCookieService{maxAge}, nil
}

// Create generates a new authentication http.Cookie from an auth token string.
func (this *authCookieService) Create(tokenString string) (*http.Cookie) {

	return &http.Cookie{
		Name: `Auth`,
		Value: tokenString,
		Expires: time.Now().Add(this.maxAge),
		HttpOnly: true,
	}
}

// Read extracts the 'Auth' http.Cookie from an http.Request.
func (this *authCookieService) Read(request *http.Request) (string, error) {

	if cookie, err := request.Cookie(`Auth`); err != nil {
		return ``, err
	} else {
		return cookie.Value, nil
	}
}
