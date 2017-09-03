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

import "net/http"

type Route struct {
	Name		string
	Method		string
	Pattern		string
	HandlerFunc	http.HandlerFunc
}

type Routes []Route

var routes = Routes {

	Route {
		Name:		"Serial",
		Method:		"POST",
		Pattern:	"/serial/{objectType}",
		HandlerFunc:	SerialHandler,
	},

	Route {
		Name:		"Checkin",
		Method:		"POST",
		Pattern:	"/checkin/{objectType}",
		HandlerFunc:	CheckinHandler,
	},

	Route {
		Name:		"Audit",
		Method:		"POST",
		Pattern:	"/audit/{serialNum}",
		HandlerFunc:	AuditHandler,
	},
}
