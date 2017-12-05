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
	`github.com/gorilla/handlers`
	`github.com/jscherff/cmdbd/utils`
	`github.com/jscherff/gox/log`
)

// Router is a Gorilla Mux router with additional methods.
type Router struct {
	*mux.Router
	MiddleWare MiddleWare
	AccessLog log.MLogger
	RecoveryLog log.MLogger
	RecoveryStack bool
}

// NewRouter creates and initializes a new Router instance.
func NewRouter(cf string, mware MiddleWare, alog, rlog log.MLogger) (*Router, error) {

	this := &Router{
		Router: mux.NewRouter().StrictSlash(true),
		MiddleWare: mware,
		AccessLog: alog,
		RecoveryLog: rlog,
	}

	if err := utils.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	return this, nil
}

// AddRoutes adds a collection of one or more routes to the Router.
func (this *Router) AddRoutes(routes Routes) *Router {

	for _, route := range routes {

		var handler http.Handler

		handler = route.HandlerFunc

		if route.Protected {
			handler = this.MiddleWare.AuthTokenValidator(handler)
		}

		handler = handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(this.RecoveryStack),
			handlers.RecoveryLogger(this.RecoveryLog))(handler)

		handler = handlers.CombinedLoggingHandler(
			this.AccessLog, handler)

		this.	Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return this
}
