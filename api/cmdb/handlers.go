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
	`github.com/jscherff/gox/log`
)

// V1 is an interface that contains V1 http.HandleFunc signatures of the cmdb API.
type HandlersV1 interface {
	SetAuthToken(http.ResponseWriter, *http.Request)
}

// v1 is a http.HandleFunc object that implements the cmdb.V1 interface.
type handlersV1 struct {
	errorLog log.MLogger
	systemLog log.MLogger
}

// NewV1 returns a new instance of an object implementing the cmdb.V1 interface.
func NewHandlersV1(errLog, sysLog log.MLogger) HandlersV1 {
	return &handlersV1{
		errorLog: errLog,
		systemLog: sysLog,
	}
}

// cmdbAuthSetTokenV1 authenticates client using basic authentication and
// issues a JWT for API authentication if successful.
func (this *handlersV1) SetAuthToken(w http.ResponseWriter, r *http.Request) {
	/*
	if user, pass, ok := r.BasicAuth(); !ok {
		el.Print(`missing credentials`)
		http.Error(w, `missing credentials`, http.StatusUnauthorized)
	} else if token, err := createAuthToken(user, pass, r.URL.Host); err != nil {
		el.Print(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
	} else if cookie, err := createAuthCookie(token); err != nil {
		el.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		sl.Printf(`issuing auth token to %q at %q`, user, r.RemoteAddr)
		http.SetCookie(w, cookie)
	}
	*/
}
