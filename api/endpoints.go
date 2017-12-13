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
	`github.com/gorilla/mux`
)

type Endpoint struct {
	Name string
	Path string
	Method string
	Handler http.Handler
	Protected bool
}

type Endpoints []Endpoint

type endpoint struct {


func (this *handlers) AddRoutes(router *mux.Router, mware http.Handler) *mux.Router {

	router.NewRoute().
		Name(`CMDB Authenticator`).
		Path(`/v2/cmdb/authenticate`).
		HandlerFunc(SetAuthToken).
		Methods(`GET`)
}
