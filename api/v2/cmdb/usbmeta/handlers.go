
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
	`github.com/jscherff/cmdbd/api`
	`github.com/jscherff/cmdbd/service`
	`github.com/jscherff/gox/log`
)

// Package variables required for operation.
var (
	systemLog, errorLog log.MLogger
	metaUsbSvc service.MetaUsbSvc
)

// Init initializes the package variables required for operation.
func Init(ms service.MetaUsbSvc, sl, el log.MLogger) {
	metaUsbSvc, systemLog, errorLog = ms, sl, el
}

// Vendor returns the USB vendor name associated with a vendor Id.
func Vendor(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.Header().Set(`Content-Type`, `application/json; charset=UTF8`)

	if vName, err := metaUsbSvc.VendorName(vars[`vid`]); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusNotFound)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(vName); err != nil {
			errorLog.Panic(err)
		}
	}
}

// Product returns the USB vendor and product names associated with
// a vendor and product Id.
func Product(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.Header().Set(`Content-Type`, `application/json; charset=UTF8`)

	if pName, err := metaUsbSvc.ProductName(vars[`vid`], vars[`pid`]); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusNotFound)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(pName); err != nil {
			errorLog.Panic(err)
		}
	}
}

// Class returns the USB class description associated with a class Id.
func Class(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.Header().Set(`Content-Type`, `application/json; charset=UTF8`)

	if cDesc, err := metaUsbSvc.ClassDesc(vars[`cid`]); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusNotFound)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(cDesc); err != nil {
			errorLog.Panic(err)
		}
	}
}

// SubClass returns the USB class and subclass descriptions associated
// with a class and subclass Id.
func SubClass(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.Header().Set(`Content-Type`, `application/json; charset=UTF8`)

	if sDesc, err := metaUsbSvc.SubClassDesc(vars[`cid`], vars[`sid`]); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusNotFound)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(sDesc); err != nil {
			errorLog.Panic(err)
		}
	}
}

// Protocol returns the USB class, subclass, and protocol descriptions
// associated with a class, subclass, and protocol Id.
func Protocol(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.Header().Set(`Content-Type`, `application/json; charset=UTF8`)

	if pDesc, err := metaUsbSvc.ProtocolDesc(vars[`cid`], vars[`sid`], vars[`pid`]); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusNotFound)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(pDesc); err != nil {
			errorLog.Panic(err)
		}
	}
}
