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

import `github.com/jscherff/cmdbd/server`

// Routes is a collection of REST API enpoints supporting USB CIs.
var Routes = server.Routes {

	server.Route {
		Name:		`USBCI Checkin Handler`,
		Method:		`POST`,
		Pattern:	`/v1/usbci/checkin/{host}/{vid}/{pid}`,
		HandlerFunc:	CheckinV1,
		Protected:	true,
	},

	server.Route {
		Name:		`USBCI Checkout Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbci/checkout/{host}/{vid}/{pid}/{sn}`,
		HandlerFunc:	CheckoutV1,
		Protected:	true,
	},

	server.Route {
		Name:		`USBCI NewSn Handler`,
		Method:		`POST`,
		Pattern:	`/v1/usbci/newsn/{host}/{vid}/{pid}`,
		HandlerFunc:	NewSnV1,
		Protected:	true,
	},

	server.Route {
		Name:		`USBCI Audit Handler`,
		Method:		`POST`,
		Pattern:	`/v1/usbci/audit/{host}/{vid}/{pid}/{sn}`,
		HandlerFunc:	AuditV1,
		Protected:	true,
	},
}
