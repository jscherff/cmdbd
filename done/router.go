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

package main

import (
	`net/http`
	`github.com/gorilla/mux`
	`github.com/gorilla/handlers`
)

// Router is a Gorilla Mux router with additional methods.
type Router struct {
	*mux.Router
}

// NewRouter instantiates a new Router.
func NewRouter() *Router {
	return &Router{mux.NewRouter().StrictSlash(true)}
}

// AddRoutes adds a collection of one or more routes to the Router.
func (this *Router) AddRoutes(routes Routes) *Router {

	var (
		accessLog = conf.Logger.Logs[`access`]
		recoveryLog = conf.Logger.Logs[`error`]
		recoveryStack = conf.Logger.RecoveryStack
	)

	for _, route := range routes {

		var handler http.Handler

		handler = route.HandlerFunc

		handler = handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(recoveryStack),
			handlers.RecoveryLogger(recoveryLog))(handler)

		handler = handlers.CombinedLoggingHandler(
			accessLog, handler)

		this.	Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return this
}
