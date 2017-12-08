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

// HandlersV2 contains http.HandleFunc signatures of CMDBd APIv2.
type HandlersV2 interface {
	SetAuthToken(http.ResponseWriter, *http.Request)
}

// handlersV2 implements the HandlersV2 interface.
type handlersV2 struct {
	errorLog log.MLogger
	systemLog log.MLogger
}

// NewHandlersV2 returns a new handlersV2 instance.
func NewHandlersV2(errLog, sysLog log.MLogger) HandlersV2 {
	return &handlersV2{
		errorLog: errLog,
		systemLog: sysLog,
	}
}

// SetAuthToken authenticates client using basic authentication and
// issues a token for API authentication if successful.
func (this *handlersV2) SetAuthToken(w http.ResponseWriter, r *http.Request) {

	if user, pass, ok := r.BasicAuth(); !ok {

		errorLog.Print(`missing credentials`)
		http.Error(w, `missing credentials`, http.StatusUnauthorized)

	} else if token, err := createAuthToken(user, pass, r.URL.Host); err != nil {

		errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)

	} else if cookie, err := createAuthCookie(token); err != nil {

		errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		systemLog.Printf(`issuing auth token to %q at %q`, user, r.RemoteAddr)
		http.SetCookie(w, cookie)

	}
}
