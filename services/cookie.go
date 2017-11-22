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
	Create(tokenString string, maxAge time.Duration) (*http.Cookie)
	Extract(request *http.Request) (*http.Cookie, error)
}

// authCookieService is a service that implements the AuthCookieService interface.
type authCookieService struct {}

// NewAuthCookieService returns an object that implements the AuthCookieService interface.
func NewAuthCookieService() AuthCookieService {
	return &authCookieService{}
}

// Create generates a new authentication http.Cookie from an auth token string.
func (this *authCookieService) Create(tokenString string, maxAge time.Duration) (*http.Cookie) {

	return &http.Cookie{
		Name: `Auth`,
		Value: tokenString,
		Expires: time.Now().Add(maxAge),
		HttpOnly: true,
	}
}

// Extract extracts the 'Auth' http.Cookie from an http.Request.
func (this *authCookieService) Extract(request *http.Request) (*http.Cookie, error) {

	if cookie, err := request.Cookie(`Auth`); err != nil {
		return nil, err
	} else {
		return cookie, nil
	}
}
