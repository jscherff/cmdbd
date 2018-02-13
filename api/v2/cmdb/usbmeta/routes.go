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

import `github.com/jscherff/cmdbd/api`

// Routes is a collection of HTTP verb/path-to-handler-function mappings.
var Routes = api.Routes {

	api.Route {
		Name:		`USB Metadata Vendor Handler`,
		Path:		`/v2/cmdb/meta/usb/vendor/{vid}`,
		Method:		`GET`,
		HandlerFunc:	Vendor,
		Protected:	false,
	},

	api.Route {
		Name:		`USB Metadata Product Handler`,
		Path:		`/v2/cmdb/meta/usb/vendor/{vid}/{pid}`,
		Method:		`GET`,
		HandlerFunc:	Product,
		Protected:	false,
	},

	api.Route {
		Name:		`USB Metadata Class Handler`,
		Path:		`/v2/cmdb/meta/usb/class/{cid}`,
		Method:		`GET`,
		HandlerFunc:	Class,
		Protected:	false,
	},

	api.Route {
		Name:		`USB Metadata SubClass Handler`,
		Path:		`/v2/cmdb/meta/usb/subclass/{cid}/{sid}`,
		Method:		`GET`,
		HandlerFunc:	SubClass,
		Protected:	false,
	},

	api.Route {
		Name:		`USB Metadata Protocol Handler`,
		Path:		`/v2/cmdb/meta/usb/protocol/{cid}/{sid}/{pid}`,
		Method:		`GET`,
		HandlerFunc:	Protocol,
		Protected:	false,
	},
}
