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

// NewRoutes returns a collection of REST APIv2 endpoints providing CMDB CI data.
func NewRoutes(hf Handlers) server.Routes {

	return server.Routes {

		server.Route {
			Name:		`USB CI Checkin Handler`,
			Method:		`POST`,
			Pattern:	`/v2/cmdb/usbci/checkin/{host}/{vid}/{pid}`,
			HandlerFunc:	hf.Checkin,
			Protected:	true,
		},

		server.Route {
			Name:		`USB CI Checkout Handler`,
			Method:		`GET`,
			Pattern:	`/v2/cmdb/usbci/checkout/{host}/{vid}/{pid}/{sn}`,
			HandlerFunc:	hf.CheckOut,
			Protected:	true,
		},

		server.Route {
			Name:		`USB CI NewSn Handler`,
			Method:		`POST`,
			Pattern:	`/v2/cmdb/usbci/newsn/{host}/{vid}/{pid}`,
			HandlerFunc:	hf.NewSn,
			Protected:	true,
		},

		server.Route {
			Name:		`USB CI Audit Handler`,
			Method:		`POST`,
			Pattern:	`/v2/cmdb/usbci/audit/{host}/{vid}/{pid}/{sn}`,
			HandlerFunc:	hf.Audit,
			Protected:	true,
		},
	}
}
