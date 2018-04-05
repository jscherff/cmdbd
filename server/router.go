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

package server

import (
	`net/http`
	`github.com/gorilla/mux`
	`github.com/jscherff/cmdbd/api`
	`github.com/jscherff/cmdbd/service`
)

// Router is a Gorilla Mux router with additional methods.
type Router struct {
	*mux.Router
	AuthSvc service.AuthSvc
}

// NewRouter creates and initializes a new Router instance.
func NewRouter(as service.AuthSvc) (*Router, error) {
	mr := mux.NewRouter().StrictSlash(true)
	return &Router{Router: mr, AuthSvc: as}, nil
}

// AddRoutes adds a collection of one or more routes to the Router.
func (this *Router) AddRoutes(routes api.Routes) *Router {

	for _, route := range routes {

		var handler http.Handler = route.HandlerFunc

		if route.Protected {
			handler = AuthTokenHandler(handler, this.AuthSvc)
		}

		this.NewRoute().
			Name(route.Name).
			Path(route.Path).
			Methods(route.Method).
			Handler(handler)
	}

	return this
}
