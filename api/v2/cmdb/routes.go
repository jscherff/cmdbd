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

import `github.com/jscherff/cmdbd/api`

// Routes is a collection of HTTP verb/path-to-handler-function mappings.
var Routes = api.Routes {

	api.Route {
		Name:		`CMDB Authenticator`,
		Path:		`/v2/cmdb/authenticate/{host}`,
		Method:		`GET`,
		HandlerFunc:	SetAuthToken,
		Protected:	false,
	},

	api.Route {
		Name:		`CMDB Event Logger`,
		Path:		`/v2/cmdb/event/create/{host}`,
		Method:		`POST`,
		HandlerFunc:	CreateEvent,
		Protected:	false,
	},

	api.Route {
		Name:		`CMDB Health Check`,
		Path:		`/v2/cmdb/health/check`,
		Method:		`GET`,
		HandlerFunc:	CheckHealth,
		Protected:	false,
	},
}
