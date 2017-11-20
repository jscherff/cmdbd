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

package usb

import (
	`encoding/json`
	`net/http`
	`github.com/gorilla/mux`
	`github.com/jscherff/gox/log`
	`github.com/jscherff/cmdb/metaci/peripheral`

)

type ApiV1 interface {
	Vendor(http.ResponseWriter, *http.Request)
	Product(http.ResponseWriter, *http.Request)
	Class(http.ResponseWriter, *http.Request)
	SubClass(http.ResponseWriter, *http.Request)
	Protocol(http.ResponseWriter, *http.Request)
}

type apiV1 struct {
	errLog log.MLogger
	sysLog log.MLogger
	meta *peripheral.Usb
}

func NewApiV1(errLog, sysLog log.MLogger, meta *peripheral.Usb) ApiV1 {
	return &apiV1{
		errLog: errLog,
		sysLog: sysLog,
	}
}

// Vendor returns the USB vendor name associated with a vendor ID.
func (this *apiV1) Vendor(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	var vid = vars[`vid`]
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if v, err := this.meta.GetVendor(vid); err != nil {

		this.errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(v.String()); err != nil {
			this.errLog.Panic(err)
		}
	}
}

// Product returns the USB vendor and product names associated with
// a vendor and product ID.
func (this *apiV1) Product(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	var vid, pid = vars[`vid`], vars[`pid`]
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)


	if v, err := this.meta.GetVendor(vid); err != nil {

		this.errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if p, err := v.GetProduct(pid); err != nil {

		this.errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(p.String()); err != nil {
			this.errLog.Panic(err)
		}
	}
}

// Class returns the USB class description associated with a class ID.
func (this *apiV1) Class(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	var cid = vars[`cid`]
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)


	if c, err := this.meta.GetClass(cid); err != nil {

		this.errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(c.String()); err != nil {
			this.errLog.Panic(err)
		}
	}
}

// SubClass returns the USB class and subclass descriptions associated
// with a class and subclass ID.
func (this *apiV1) SubClass(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	var cid, sid = vars[`cid`], vars[`sid`]
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)


	if c, err := this.meta.GetClass(cid); err != nil {

		this.errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if s, err := c.GetSubClass(sid); err != nil {

		this.errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(s.String()); err != nil {
			this.errLog.Panic(err)
		}
	}
}

// Protocol returns the USB class, subclass, and protocol descriptions
// associated with a class, subclass, and protocol ID.
func (this *apiV1) Protocol(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	var cid, sid, pid = vars[`cid`], vars[`sid`], vars[`pid`]
	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)


	if c, err := this.meta.GetClass(cid); err != nil {

		this.errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if s, err := c.GetSubClass(sid); err != nil {

		this.errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if p, err := s.GetProtocol(pid); err != nil {

		this.errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(p.String()); err != nil {
			this.errLog.Panic(err)
		}
	}
}
