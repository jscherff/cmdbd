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
	`github.com/jscherff/cmdbd/api`
	`github.com/jscherff/cmdbd/service`
	`github.com/jscherff/cmdbd/utils`
)

// Router is a Gorilla Mux router with additional methods.
type Router struct {
	*mux.Router
	RecoveryStack bool
	AuthSvc service.AuthSvc
	LoggerSvc service.LoggerSvc
}

// NewRouter creates and initializes a new Router instance.
func NewRouter(cf string, authSvc service.AuthSvc, logSvc service.LoggerSvc) (*Router, error) {

	this := &Router{
		Router: mux.NewRouter().StrictSlash(true),
		RecoveryStack: false,
		AuthSvc: authSvc,
		LoggerSvc: logSvc,
	}

	if err := utils.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	return this, nil
}

// AddRoutes adds a collection of one or more routes to the Router.
func (this *Router) AddRoutes(routes api.Routes) *Router {

	accessLog := this.LoggerSvc.AccessLog()
	recoveryLog := this.LoggerSvc.ErrorLog()

	for _, route := range routes {

		var handler http.Handler = route.HandlerFunc

		handler = handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(this.RecoveryStack),
			handlers.RecoveryLogger(recoveryLog))(handler)

		if route.Protected {
			handler = AuthTokenHandler(this.AuthSvc, handler)
		}

		handler = handlers.CombinedLoggingHandler(accessLog, handler)

		this.NewRoute().
			Name(route.Name).
			Path(route.Path).
			Methods(route.Method).
			Handler(handler)
	}

	return this
}
