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
	`github.com/jscherff/gox/log`
)

const HttpBodySizeLimit = 1048576

var errLog, sysLog log.MLogger

func SetLoggers(eLog, sLog log.MLogger) {
	errLog, sysLog = eLog, sLog
}

// Vendor returns the USB vendor name associated with a vendor ID.
func VendorV1(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	usb := conf.MetaUsb

	var vid = vars[`vid`]
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if v, err := usb.GetVendor(vid); err != nil {

		errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(v.String()); err != nil {
			errLog.Panic(err)
		}
	}
}

// Product returns the USB vendor and product names associated with
// a vendor and product ID.
func ProductV1(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	usb := conf.MetaUsb

	var vid, pid = vars[`vid`], vars[`pid`]
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)


	if v, err := usb.GetVendor(vid); err != nil {

		errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if p, err := v.GetProduct(pid); err != nil {

		errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(p.String()); err != nil {
			errLog.Panic(err)
		}
	}
}

// Class returns the USB class description associated with a class ID.
func ClassV1(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	usb := conf.MetaUsb

	var cid = vars[`cid`]
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)


	if c, err := usb.GetClass(cid); err != nil {

		errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(c.String()); err != nil {
			errLog.Panic(err)
		}
	}
}

// SubClass returns the USB class and subclass descriptions associated
// with a class and subclass ID.
func SubClassV1(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	usb := conf.MetaUsb

	var cid, sid = vars[`cid`], vars[`sid`]
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)


	if c, err := usb.GetClass(cid); err != nil {

		errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if s, err := c.GetSubClass(sid); err != nil {

		errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(s.String()); err != nil {
			errLog.Panic(err)
		}
	}
}

// Protocol returns the USB class, subclass, and protocol descriptions
// associated with a class, subclass, and protocol ID.
func ProtocolV1(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	usb := conf.MetaUsb

	var cid, sid, pid = vars[`cid`], vars[`sid`], vars[`pid`]
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)


	if c, err := usb.GetClass(cid); err != nil {

		errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if s, err := c.GetSubClass(sid); err != nil {

		errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if p, err := s.GetProtocol(pid); err != nil {

		errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(p.String()); err != nil {
			errLog.Panic(err)
		}
	}
}
