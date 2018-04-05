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

package usbci

import `github.com/jscherff/cmdbd/api`

// Routes is a collection of HTTP verb/path-to-handler-function mappings.
var Routes = api.Routes {

	api.Route {
		Name:		`USB CI CheckIn Handler`,
		Path:		`/api/v3/cmdb/ci/usb/checkin/{host}/{vid}/{pid}`,
		Method:		`POST`,
		HandlerFunc:	CheckIn,
		Protected:	false,
	},

	api.Route {
		Name:		`USB CI Checkout Handler`,
		Path:		`/api/v3/cmdb/ci/usb/checkout/{host}/{vid}/{pid}/{sn}`,
		Method:		`GET`,
		HandlerFunc:	CheckOut,
		Protected:	false,
	},

	api.Route {
		Name:		`USB CI NewSn Handler`,
		Path:		`/api/v3/cmdb/ci/usb/newsn/{host}/{vid}/{pid}`,
		Method:		`POST`,
		HandlerFunc:	NewSn,
		Protected:	false,
	},

	api.Route {
		Name:		`USB CI Audit Handler`,
		Path:		`/api/v3/cmdb/ci/usb/audit/{host}/{vid}/{pid}/{sn}`,
		Method:		`POST`,
		HandlerFunc:	Audit,
		Protected:	false,
	},
}
