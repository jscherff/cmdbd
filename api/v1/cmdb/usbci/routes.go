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

import (
	`github.com/jscherff/cmdbd/api`
	v2 `github.com/jscherff/cmdbd/api/v2/cmdb/usbci`
)

// Routes is a collection of HTTP verb/path-to-handler-function mappings.
var Routes = api.Routes {

	api.Route {
		Name:		`USBCI CheckIn Handler`,
		Path:		`/v1/usbci/checkin/{host}/{vid}/{pid}`,
		Method:		`POST`,
		HandlerFunc:	v2.CheckIn,
		Protected:	true,
	},

	api.Route {
		Name:		`USBCI CheckOut Handler`,
		Path:		`/v1/usbci/checkout/{host}/{vid}/{pid}/{sn}`,
		Method:		`GET`,
		HandlerFunc:	v2.CheckOut,
		Protected:	true,
	},

	api.Route {
		Name:		`USBCI NewSn Handler`,
		Path:		`/v1/usbci/newsn/{host}/{vid}/{pid}`,
		Method:		`POST`,
		HandlerFunc:	v2.NewSn,
		Protected:	true,
	},

	api.Route {
		Name:		`USBCI Audit Handler`,
		Path:		`/v1/usbci/audit/{host}/{vid}/{pid}/{sn}`,
		Method:		`POST`,
		HandlerFunc:	v2.Audit,
		Protected:	true,
	},
}
