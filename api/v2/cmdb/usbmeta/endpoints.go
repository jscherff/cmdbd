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
	`encoding/json`
	`net/http`
	`github.com/gorilla/mux`
	`github.com/jscherff/cmdb/meta/peripheral`
	`github.com/jscherff/cmdbd/api`
	`github.com/jscherff/cmdbd/service`

)

// Package variables required for operation.
var (
	usbMetadata *peripheral.Usb
	loggerSvc service.LoggerSvc
)

// Init initializes the package variables required for operation.
func Init(um *peripheral.Usb, ls service.LoggerSvc) {
	usbMetadata, loggerSvc = um, ls
}

// Endpoints is a collection of URL path to handler function mappings.
var Endpoints = []api.Endpoint {

	api.Endpoint {
		Name:		`USB Metadata Vendor Handler`,
		Path:		`/v2/cmdb/meta/usb/vendor/{vid}`,
		Method:		`GET`,
		HandlerFunc:	Vendor,
		Protected:	false,
	},

	api.Endpoint {
		Name:		`USB Metadata Product Handler`,
		Path:		`/v2/cmdb/meta/usb/vendor/{vid}/{pid}`,
		Method:		`GET`,
		HandlerFunc:	Product,
		Protected:	false,
	},

	api.Endpoint {
		Name:		`USB Metadata Class Handler`,
		Path:		`/v2/cmdb/meta/usb/class/{cid}`,
		Method:		`GET`,
		HandlerFunc:	Class,
		Protected:	false,
	},

	api.Endpoint {
		Name:		`USB Metadata SubClass Handler`,
		Path:		`/v2/cmdb/meta/usb/subclass/{cid}/{sid}`,
		Method:		`GET`,
		HandlerFunc:	SubClass,
		Protected:	false,
	},

	api.Endpoint {
		Name:		`USB Metadata Protocol Handler`,
		Path:		`/v2/cmdb/meta/usb/protocol/{cid}/{sid}/{pid}`,
		Method:		`GET`,
		HandlerFunc:	Protocol,
		Protected:	false,
	},
}

// Vendor returns the USB vendor name associated with a vendor ID.
func Vendor(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if v, err := usbMetadata.GetVendor(vars[`vid`]); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(v.String()); err != nil {
			loggerSvc.ErrorLog().Panic(err)
		}
	}
}

// Product returns the USB vendor and product names associated with
// a vendor and product ID.
func Product(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if v, err := usbMetadata.GetVendor(vars[`vid`]); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if p, err := v.GetProduct(vars[`pid`]); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(p.String()); err != nil {
			loggerSvc.ErrorLog().Panic(err)
		}
	}
}

// Class returns the USB class description associated with a class ID.
func Class(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if c, err := usbMetadata.GetClass(vars[`cid`]); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(c.String()); err != nil {
			loggerSvc.ErrorLog().Panic(err)
		}
	}
}

// SubClass returns the USB class and subclass descriptions associated
// with a class and subclass ID.
func SubClass(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if c, err := usbMetadata.GetClass(vars[`cid`]); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if s, err := c.GetSubClass(vars[`sid`]); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(s.String()); err != nil {
			loggerSvc.ErrorLog().Panic(err)
		}
	}
}

// Protocol returns the USB class, subclass, and protocol descriptions
// associated with a class, subclass, and protocol ID.
func Protocol(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if c, err := usbMetadata.GetClass(vars[`cid`]); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if s, err := c.GetSubClass(vars[`sid`]); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if p, err := s.GetProtocol(vars[`pid`]); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(p.String()); err != nil {
			loggerSvc.ErrorLog().Panic(err)
		}
	}
}
