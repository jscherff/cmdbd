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

package usb

import (
	`github.com/jscherff/cmdbd/server`
)

// NewRoutesV1 returns a collection of REST API endpoints providing CMDB CI data.
func NewRoutesV1(hf HandlerFuncsV1) server.Routes {

	return server.Routes {

		server.Route {
			Name:		`USB CI Checkin Handler`,
			Method:		`POST`,
			Pattern:	`/v1/usbci/checkin/{host}/{vid}/{pid}`,
			HandlerFunc:	hf.CheckIn,
			Protected:	true,
		},

		server.Route {
			Name:		`USBCI Checkout Handler`,
			Method:		`GET`,
			Pattern:	`/v1/usbci/checkout/{host}/{vid}/{pid}/{sn}`,
			HandlerFunc:	hf.CheckOut,
			Protected:	true,
		},

		server.Route {
			Name:		`USBCI NewSn Handler`,
			Method:		`POST`,
			Pattern:	`/v1/usbci/newsn/{host}/{vid}/{pid}`,
			HandlerFunc:	hf.NewSn,
			Protected:	true,
		},

		server.Route {
			Name:		`USBCI Audit Handler`,
			Method:		`POST`,
			Pattern:	`/v1/usbci/audit/{host}/{vid}/{pid}/{sn}`,
			HandlerFunc:	hf.Audit,
			Protected:	true,
		},
	}
}
