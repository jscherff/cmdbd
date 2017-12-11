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

package cmdb

import (
	`net/http`
	`github.com/jscherff/cmdbd/model/cmdb`
	`github.com/jscherff/cmdbd/service`
)

// Handlers contains http.HandleFunc signatures of CMDBd APIv2.
type Handlers interface {
	SetAuthToken(http.ResponseWriter, *http.Request)
}

// handlers implements the Handlers interface.
type handlers struct {
	ErrorLog service.Logger
	SystemLog service.Logger
	AuthTokenSvc service.AuthTokenService
	AuthCookieSvc service.AuthCookieService
}

// NewHandlers returns a new handlers instance.
func NewHandlers(errLog, sysLog service.Logger, ats service.AuthTokenService, acs service.AuthCookieService) Handlers {
	return &handlers{
		ErrorLog: errLog,
		SystemLog: sysLog,
		AuthTokenSvc: ats,
	}
}

// SetAuthToken authenticates client using basic authentication and
// issues a token for API authentication if successful.
func (this *handlers) SetAuthToken(w http.ResponseWriter, r *http.Request) {

	user := &cmdb.User{}

	username, password, ok := r.BasicAuth()

	if !ok {
		this.ErrorLog.Print(`missing credentials`)
		http.Error(w, `missing credentials`, http.StatusUnauthorized)
	}

	user.Username = username

	if err := user.Read(user); err != nil {
		this.ErrorLog.Printf(`user %q not found`, username)
		http.Error(w, `user not found`, http.StatusNotFound)
	} else if err := user.Verify(password); err != nil {
		this.ErrorLog.Print(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	if token, err := this.AuthTokenSvc.Create(user); err != nil {
		this.ErrorLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if tokenString, err := this.AuthTokenSvc.String(token); err != nil {
		this.ErrorLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if cookie, err := this.AuthCookieSvc.Create(tokenString); err != nil {
		this.ErrorLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		this.SystemLog.Printf(`issuing auth token to %q at %q`, user.Username, r.RemoteAddr)
		http.SetCookie(w, cookie)
	}
}
