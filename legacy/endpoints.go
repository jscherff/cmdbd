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

package legacy

import `github.com/jscherff/cmdbd/server`

// usbCiRoutes is a collection of REST API enpoints supporting USB CIs.
var usbCiRoutes = server.Routes {

	server.Route {
		Name:		`USBCI Checkin Handler`,
		Method:		`POST`,
		Pattern:	`/v1/usbci/checkin/{host}/{vid}/{pid}`,
		HandlerFunc:	usbCiCheckinV1,
		Protected:	true,
	},

	server.Route {
		Name:		`USBCI Checkout Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbci/checkout/{host}/{vid}/{pid}/{sn}`,
		HandlerFunc:	usbCiCheckoutV1,
		Protected:	true,
	},

	server.Route {
		Name:		`USBCI NewSn Handler`,
		Method:		`POST`,
		Pattern:	`/v1/usbci/newsn/{host}/{vid}/{pid}`,
		HandlerFunc:	usbCiNewSnV1,
		Protected:	true,
	},

	server.Route {
		Name:		`USBCI Audit Handler`,
		Method:		`POST`,
		Pattern:	`/v1/usbci/audit/{host}/{vid}/{pid}/{sn}`,
		HandlerFunc:	usbCiAuditV1,
		Protected:	true,
	},
}

// usbMetaRoutes is a collection of REST API enpoints providing USB metadata.
var usbMetaRoutes = server.Routes {

	server.Route {
		Name:		`Metadata USB Vendor Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbmeta/vendor/{vid}`,
		HandlerFunc:	usbMetaVendorV1,
		Protected:	false,
	},

	server.Route {
		Name:		`Metadata USB Product Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbmeta/vendor/{vid}/{pid}`,
		HandlerFunc:	usbMetaProductV1,
		Protected:	false,
	},

	server.Route {
		Name:		`Metadata USB Class Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbmeta/class/{cid}`,
		HandlerFunc:	usbMetaClassV1,
		Protected:	false,
	},

	server.Route {
		Name:		`Metadata USB SubClass Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbmeta/subclass/{cid}/{sid}`,
		HandlerFunc:	usbMetaSubClassV1,
		Protected:	false,
	},

	server.Route {
		Name:		`Metadata USB Protocol Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbmeta/protocol/{cid}/{sid}/{pid}`,
		HandlerFunc:	usbMetaProtocolV1,
		Protected:	false,
	},
}

var cmdbAuthRoutes = server.Routes {

	server.Route {
		Name:		`CMDB Authenticator`,
		Method:		`GET`,
		Pattern:	`/v1/cmdbauth`,
		HandlerFunc:	cmdbAuthSetTokenV1,
		Protected:	false,
	},
}
