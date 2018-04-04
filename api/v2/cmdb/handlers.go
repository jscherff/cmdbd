
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
	`encoding/json`
	`fmt`
	`net/http`
	`time`
	`github.com/gorilla/mux`
	`github.com/jscherff/cmdbd/api`
	`github.com/jscherff/cmdbd/model/cmdb`
	`github.com/jscherff/cmdbd/service`
	`github.com/jscherff/gox/log`
)

// Templates for system and error messages.
const (
	fmtHostInfo = `host '%s' at '%s'`
	fmtAuthMissingCreds = `auth failure: missing credentials from ` + fmtHostInfo
	fmtAuthFailure = `auth failure for user '%s' on ` + fmtHostInfo + `: %v`
	fmtAuthSuccess = `auth success for user '%s' on ` + fmtHostInfo
	fmtEventSuccess = `event logged for ` + fmtHostInfo
	fmtHealthSuccess = `health check success for host at '%s'`
	fmtConnectSuccess = `connection id '%s' established for host at '%s'`
)

// Package variables required for operation.
var (
	authSvc service.AuthSvc
	systemLog, errorLog log.MLogger
)

// Init initializes the package variables required for operation.
func Init(as service.AuthSvc, sl, el log.MLogger) {
	authSvc, systemLog, errorLog = as, sl, el
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

	w.Header().Set(`Content-Type`, `text/plain; charset=UTF8`)

	if user.Username, passwd, ok = r.BasicAuth(); !ok {

		err := fmt.Errorf(fmtAuthMissingCreds, vars[`host`], r.RemoteAddr)
		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusUnauthorized)

	} else if err := user.Read(); err != nil {

		err = fmt.Errorf(fmtAuthFailure, user.Username, vars[`host`], r.RemoteAddr, err)
		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if err := user.VerifyPassword(passwd); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusUnauthorized)

	} else if err := user.VerifyAccess(); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusUnauthorized)

	} else if token, err := authSvc.CreateToken(user); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if tokenString, err := authSvc.CreateTokenString(token); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if cookie, err := authSvc.CreateCookie(tokenString); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		systemLog.Printf(fmtAuthSuccess, user.Username, vars[`host`], r.RemoteAddr)
		http.SetCookie(w, cookie)
	}
}

// CreateEvent writes an event to the datastore event log.
func CreateEvent(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	event := &cmdb.Event{}
	event.HostName, event.RemoteAddr = vars[`host`], r.RemoteAddr

	w.Header().Set(`Content-Type`, `text/plain; charset=UTF8`)

	if _, err := event.Create(); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		systemLog.Printf(fmtEventSuccess, event.HostName, event.RemoteAddr)
		w.WriteHeader(http.StatusCreated)
	}
}

// CheckHealth returns the health of the server.
func CheckHealth(w http.ResponseWriter, r *http.Request) {

	info := &cmdb.Info{
		Client:	r.RemoteAddr,
		Server:	r.URL.Host,
		Proto:	r.Proto,
		Method:	r.Method,
		Scheme:	r.URL.Scheme,
		Path:	r.URL.Path,
	}

	w.Header().Set(`Content-Type`, `text/plain; charset=UTF8`)

	if err := info.Read(); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		systemLog.Printf(fmtHealthSuccess, r.RemoteAddr)
		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode(info); err != nil {
			errorLog.Panic(err)
		}
	}
}


// CheckConcurrency allows tests against the server connection limit. It accepts
// an arbitrary connection ID from the client (presumably a counter), logs the
// connection information to the system log, and sleeps for 60 seconeds before
// returning an HTTP StatusOK to the client.
func CheckConcurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set(`Content-Type`, `text/plain; charset=UTF8`)
	systemLog.Printf(fmtConnectSuccess, vars[`id`], r.RemoteAddr)
	time.Sleep(60 * time.Second)
	w.WriteHeader(http.StatusNoContent)
}
