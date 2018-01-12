
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
	`fmt`
	`net/http`
	`github.com/jscherff/cmdbd/model/cmdb`
	`github.com/jscherff/cmdbd/service`
)

// Package variables required for operation.
var (
	authSvc service.AuthSvc
	loggerSvc service.LoggerSvc
)

// Init initializes the package variables required for operation.
func Init(as service.AuthSvc, ls service.LoggerSvc) {
	authSvc, loggerSvc = as, ls
}

// SetauthToken authenticates client using basic authentication and
// issues a token for API authentication if successful.
func SetAuthToken(w http.ResponseWriter, r *http.Request) {

	user := &cmdb.User{}

	username, password, ok := r.BasicAuth()

	if !ok {
		loggerSvc.ErrorLog().Print(`missing credentials`)
		http.Error(w, `missing credentials`, http.StatusUnauthorized)
	}

	user.Username = username

	if err := user.Read(); err != nil {

		err = fmt.Errorf(`user %q not found: %v`, username, err)
		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if err := user.VerifyPassword(password); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)

	} else if err := user.VerifyAccess(); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)

	} else if token, err := authSvc.CreateToken(user); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if tokenString, err := authSvc.CreateTokenString(token); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if cookie, err := authSvc.CreateCookie(tokenString); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		loggerSvc.SystemLog().Printf(`issuing auth token to %q at %q`, user.Username, r.RemoteAddr)
		http.SetCookie(w, cookie)
	}
}
