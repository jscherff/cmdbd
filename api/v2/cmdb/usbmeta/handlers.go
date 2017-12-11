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
	`github.com/jscherff/cmdb/meta/peripheral`

)

// Handlers contains http.HandleFunc signatures of CMDBd APIv2.
type Handlers interface {
	Vendor(http.ResponseWriter, *http.Request)
	Product(http.ResponseWriter, *http.Request)
	Class(http.ResponseWriter, *http.Request)
	SubClass(http.ResponseWriter, *http.Request)
	Protocol(http.ResponseWriter, *http.Request)
}

// handlers implements the Handlers interface.
type handlers struct {
	errorLog log.MLogger
	systemLog log.MLogger
	meta *peripheral.Usb
}

// NewHandlers returns a new handlers instance.
func NewHandlers(errLog, sysLog log.MLogger, meta *peripheral.Usb) Handlers {
	return &handlers{
		errorLog: errLog,
		systemLog: sysLog,
	}
}

// Vendor returns the USB vendor name associated with a vendor ID.
func (this *handlers) Vendor(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if v, err := this.meta.GetVendor(vars[`vid`]); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(v.String()); err != nil {
			this.errorLog.Panic(err)
		}
	}
}

// Product returns the USB vendor and product names associated with
// a vendor and product ID.
func (this *handlers) Product(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if v, err := this.meta.GetVendor(vars[`vid`]); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if p, err := v.GetProduct(vars[`pid`]); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(p.String()); err != nil {
			this.errorLog.Panic(err)
		}
	}
}

// Class returns the USB class description associated with a class ID.
func (this *handlers) Class(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if c, err := this.meta.GetClass(vars[`cid`]); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(c.String()); err != nil {
			this.errorLog.Panic(err)
		}
	}
}

// SubClass returns the USB class and subclass descriptions associated
// with a class and subclass ID.
func (this *handlers) SubClass(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if c, err := this.meta.GetClass(vars[`cid`]); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if s, err := c.GetSubClass(vars[`sid`]); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(s.String()); err != nil {
			this.errorLog.Panic(err)
		}
	}
}

// Protocol returns the USB class, subclass, and protocol descriptions
// associated with a class, subclass, and protocol ID.
func (this *handlers) Protocol(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if c, err := this.meta.GetClass(vars[`cid`]); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if s, err := c.GetSubClass(vars[`sid`]); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if p, err := s.GetProtocol(vars[`pid`]); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(p.String()); err != nil {
			this.errorLog.Panic(err)
		}
	}
}
