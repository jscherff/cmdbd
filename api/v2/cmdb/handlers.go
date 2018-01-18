
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
	`github.com/gorilla/mux`
	`github.com/jscherff/cmdbd/model/cmdb`
	`github.com/jscherff/cmdbd/service`
)

// Templates for system and error messages.
const (
	fmtAuthMissingCreds = `auth failure on host '%s' at '%s' - missing credentials`
	fmtAuthUserNotFound = `auth failure on host '%s' at '%s' - user '%s' not found: %v`
	fmtAuthSuccess = `auth success on host '%s' at '%s' - issuing token for user '%s'`
	fmtEventSuccess = `event logged for host '%s' at '%s'`
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

	var (
		ok bool
		passwd string
	)

	vars := mux.Vars(r)
	user := &cmdb.User{}

	if user.Username, passwd, ok = r.BasicAuth(); !ok {

		err := fmt.Errorf(fmtAuthMissingCreds, vars[`host`], r.RemoteAddr)
		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)

	} else if err := user.Read(); err != nil {

		err = fmt.Errorf(fmtAuthUserNotFound, vars[`host`], r.RemoteAddr, user.Username, err)
		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if err := user.VerifyPassword(passwd); err != nil {

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

		loggerSvc.SystemLog().Printf(fmtAuthSuccess, vars[`host`], r.RemoteAddr, user.Username)
		http.SetCookie(w, cookie)
	}
}

// CreateEvent writes an event to the datastore event log.
func CreateEvent(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	event := &cmdb.Event{}
	event.HostName, event.RemoteAddr = vars[`host`], r.RemoteAddr

	if _, err := event.Create(); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		loggerSvc.SystemLog().Printf(fmtEventSuccess, event.HostName, event.RemoteAddr)
		w.WriteHeader(http.StatusCreated)
	}
}
