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

import "net/http"

// Route contains information about a REST API enpoint.
type Route struct {
	Name		string
	Method		string
	Pattern		string
	HandlerFunc	http.HandlerFunc
}

// Routes contains a collection of Route instances.
type Routes []Route

// usbCiRoutes is a collection of REST API enpoints supporting USB CIs.
var usbCiRoutes = Routes {

	Route {
		Name:		"USBCI Checkin Handler",
		Method:		"POST",
		Pattern:	"/v1/usbCi/checkin/{host}/{vid}/{pid}",
		HandlerFunc:	usbCiCheckinV1,
	},

	Route {
		Name:		"USBCI Checkout Handler",
		Method:		"GET",
		Pattern:	"/v1/usbCi/checkout/{host}/{vid}/{pid}/{sn}",
		HandlerFunc:	usbCiCheckoutV1,
	},

	Route {
		Name:		"USBCI NewSN Handler",
		Method:		"POST",
		Pattern:	"/v1/usbCi/newsn/{host}/{vid}/{pid}",
		HandlerFunc:	usbCiNewSNV1,
	},

	Route {
		Name:		"USBCI Audit Handler",
		Method:		"POST",
		Pattern:	"/v1/usbCi/audit/{host}/{vid}/{pid}/{sn}",
		HandlerFunc:	usbCiAuditV1,
	},
}

// usbMetaRoutes is a collection of REST API enpoints providing USB metadata.
var usbMetaRoutes = Routes {

	Route {
		Name:		"Metadata USB Vendor Handler",
		Method:		"GET",
		Pattern:	"/v1/meta/usb/{vid}",
		HandlerFunc:	usbMetaVendorV1
	},

	Route {
		Name:		"Metadata USB Product Handler",
		Method:		"GET",
		Pattern:	"/v1/meta/usb/{vid}/{pid}",
		HandlerFunc:	usbMetaProductV1
	},

	Route {
		Name:		"Metadata USB Class Handler",
		Method:		"GET",
		Pattern:	"/v1/meta/usb/{cid}",
		HandlerFunc:	usbMetaClassV1
	},

	Route {
		Name:		"Metadata USB Subclass Handler",
		Method:		"GET",
		Pattern:	"/v1/meta/usb/{cid}/{sid}",
		HandlerFunc:	usbMetaSubclassV1
	},

	Route {
		Name:		"Metadata USB Protocol Handler",
		Method:		"GET",
		Pattern:	"/v1/meta/usb/{cid}/{sid}/{pid}",
		HandlerFunc:	usbMetaProtocolV1
	},
}
