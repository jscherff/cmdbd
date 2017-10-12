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

type Route struct {
	Name		string
	Method		string
	Pattern		string
	HandlerFunc	http.HandlerFunc
}

type Routes []Route

var routes = Routes {

	Route {
		Name:		"USBCI Checkin Handler",
		Method:		"POST",
		Pattern:	"/v1/usbci/checkin/{host}/{vid}/{pid}",
		HandlerFunc:	v1usbciCheckin,
	},

	Route {
		Name:		"USBCI Checkout Handler",
		Method:		"GET",
		Pattern:	"/v1/usbci/checkout/{host}/{vid}/{pid}/{sn}",
		HandlerFunc:	v1usbciCheckout,
	},

	Route {
		Name:		"USBCI NewSN Handler",
		Method:		"POST",
		Pattern:	"/v1/usbci/newsn/{host}/{vid}/{pid}",
		HandlerFunc:	v1usbciNewSN,
	},

	Route {
		Name:		"USBCI Audit Handler",
		Method:		"POST",
		Pattern:	"/v1/usbci/audit/{host}/{vid}/{pid}/{sn}",
		HandlerFunc:	v1usbciAudit,
	},
/*
	Route {
		Name:		"Metadata USB Vendor Handler",
		Method:		"GET",
		Pattern:	"/v1/meta/usb/{vid}",
		HandlerFunc:	v1usbMetaVendorName
	},

	Route {
		Name:		"Metadata USB Product Handler",
		Method:		"GET",
		Pattern:	"/v1/meta/usb/{vid}/{pid}",
		HandlerFunc:	v1usbMetaProductName
	},

	Route {
		Name:		"Metadata USB Class Handler",
		Method:		"GET",
		Pattern:	"/v1/meta/usb/{cid}",
		HandlerFunc:	v1usbMetaClassDesc
	},

	Route {
		Name:		"Metadata USB Subclass Handler",
		Method:		"GET",
		Pattern:	"/v1/meta/usb/{cid}/{sid}",
		HandlerFunc:	v1usbMetaSubclassDesc
	},

	Route {
		Name:		"Metadata USB Protocol Handler",
		Method:		"GET",
		Pattern:	"/v1/meta/usb/{cid}/{sid}/{pid}",
		HandlerFunc:	v1usbMetaProtocolDesc
	},
*/
}
