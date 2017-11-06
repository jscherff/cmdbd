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

package main

import (
	`net/http`
)

// Route contains information about a REST API enpoint.
type Route struct {
	Name		string
	Method		string
	Pattern		string
	HandlerFunc	http.HandlerFunc
}

// Regular expression for hostname input validation
var hostRgx = `^(?:(?:\w|\w[\w\-]*\w)\.)*(?:\w|\w[\w\-]*\w)$`

// Regular expression for vendor and product ID input validation.
var hex4Rgx = `[0-9A-Fa-f]{4}`

// Regular expression for class, sublcass, and protocol ID input validation.
var hex2Rgx = `[0-9A-Fa-f]{2}`

// Regular expression for serial number input validation.
var snRgx = `^[\w\-]+$`

// Routes contains a collection of Route instances.
type Routes []Route

// usbCiRoutes is a collection of REST API enpoints supporting USB CIs.
var usbCiRoutes = Routes {

	Route {
		Name:		`USBCI Checkin Handler`,
		Method:		`POST`,
		Pattern:	`/v1/usbci/checkin/{host}/{vid}/{pid}`,
		HandlerFunc:	usbCiCheckinV1,
	},

	Route {
		Name:		`USBCI Checkout Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbci/checkout/{host}/{vid}/{pid}/{sn}`,
		HandlerFunc:	usbCiCheckoutV1,
	},

	Route {
		Name:		`USBCI NewSn Handler`,
		Method:		`POST`,
		Pattern:	`/v1/usbci/newsn/{host}/{vid}/{pid}`,
		HandlerFunc:	usbCiNewSnV1,
	},

	Route {
		Name:		`USBCI Audit Handler`,
		Method:		`POST`,
		Pattern:	`/v1/usbci/audit/{host}/{vid}/{pid}/{sn}`,
		HandlerFunc:	usbCiAuditV1,
	},
}

// usbMetaRoutes is a collection of REST API enpoints providing USB metadata.
var usbMetaRoutes = Routes {

	Route {
		Name:		`Metadata USB Vendor Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbmeta/vendor/{vid}`,
		HandlerFunc:	usbMetaVendorV1,
	},

	Route {
		Name:		`Metadata USB Product Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbmeta/vendor/{vid}/{pid}`,
		HandlerFunc:	usbMetaProductV1,
	},

	Route {
		Name:		`Metadata USB Class Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbmeta/class/{cid}`,
		HandlerFunc:	usbMetaClassV1,
	},

	Route {
		Name:		`Metadata USB SubClass Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbmeta/subclass/{cid}/{sid}`,
		HandlerFunc:	usbMetaSubClassV1,
	},

	Route {
		Name:		`Metadata USB Protocol Handler`,
		Method:		`GET`,
		Pattern:	`/v1/usbmeta/protocol/{cid}/{sid}/{pid}`,
		HandlerFunc:	usbMetaProtocolV1,
	},
}

var cmdbAuthRoutes = Routes {

	Route {
		Name:		`CMDB Authenticator Set Token`,
		Method:		`GET`,
		Pattern:	`/v1/cmdbauth/settoken`,
		HandlerFunc:	cmdbAuthSetTokenV1,
	},
}
