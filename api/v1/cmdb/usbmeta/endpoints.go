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

import `github.com/jscherff/cmdbd/server`

// Routes is a collection of REST API enpoints providing USB metadata.
var Routes = server.Routes {

	server.Route {
		Name:		`Metadata USB Vendor Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbmeta/vendor/{vid}`,
		HandlerFunc:	VendorV1,
		Protected:	false,
	},

	server.Route {
		Name:		`Metadata USB Product Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbmeta/vendor/{vid}/{pid}`,
		HandlerFunc:	ProductV1,
		Protected:	false,
	},

	server.Route {
		Name:		`Metadata USB Class Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbmeta/class/{cid}`,
		HandlerFunc:	ClassV1,
		Protected:	false,
	},

	server.Route {
		Name:		`Metadata USB SubClass Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbmeta/subclass/{cid}/{sid}`,
		HandlerFunc:	SubClassV1,
		Protected:	false,
	},

	server.Route {
		Name:		`Metadata USB Protocol Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbmeta/protocol/{cid}/{sid}/{pid}`,
		HandlerFunc:	ProtocolV1,
		Protected:	false,
	},
}
