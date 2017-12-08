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

package usbmeta

import (
	`github.com/jscherff/cmdbd/server`
)

// NewRoutesV2 returns a collection of REST API endpoints providing USB metadata.
func NewRoutesV2(hf HandlersV2) server.Routes {

	return server.Routes {

		server.Route {
			Name:		`USB Metadata Vendor Handler`,
			Method:		`GET`,
			Pattern:	`/v2/cmdb/usbmeta/vendor/{vid}`,
			HandlerFunc:	hf.Vendor,
			Protected:	false,
		},

		server.Route {
			Name:		`USB Metadata Product Handler`,
			Method:		`GET`,
			Pattern:	`/v2/cmdb/usbmeta/vendor/{vid}/{pid}`,
			HandlerFunc:	hf.Product,
			Protected:	false,
		},

		server.Route {
			Name:		`USB Metadata Class Handler`,
			Method:		`GET`,
			Pattern:	`/v2/cmdb/usbmeta/class/{cid}`,
			HandlerFunc:	hf.Class,
			Protected:	false,
		},

		server.Route {
			Name:		`USB Metadata SubClass Handler`,
			Method:		`GET`,
			Pattern:	`/v2/cmdb/usbmeta/subclass/{cid}/{sid}`,
			HandlerFunc:	hf.SubClass,
			Protected:	false,
		},

		server.Route {
			Name:		`USB Metadata Protocol Handler`,
			Method:		`GET`,
			Pattern:	`/v2/cmdb/usbmeta/protocol/{cid}/{sid}/{pid}`,
			HandlerFunc:	hf.Protocol,
			Protected:	false,
		},
	}
}
