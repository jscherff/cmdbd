
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

package usbci

import (
	`encoding/json`
	`fmt`
	`net/http`
	`github.com/gorilla/mux`
	`github.com/jscherff/cmdbd/api`
	`github.com/jscherff/cmdbd/model/cmdb`
	`github.com/jscherff/cmdbd/model/cmdb/usbci`
	`github.com/jscherff/cmdbd/service`
	`github.com/jscherff/gox/log`
)

// Templates for system and error messages.
const (
	fmtDevInfo = `USB device '%s-%s'`
	fmtHostInfo = `host '%s' at '%s'`
	fmtCheckInSuccess = `checked in ` + fmtDevInfo + ` SN '%s' for ` + fmtHostInfo
	fmtNewSnDuplicate = `duplicate SN '%s' for ` + fmtDevInfo + ` on `+ fmtHostInfo
	fmtNewSnSuccess = `issued SN '%s' to ` + fmtDevInfo + ` on ` + fmtHostInfo
	fmtAuditSuccess = `audited ` + fmtDevInfo + ` SN '%s' for ` + fmtHostInfo
	fmtCheckOutSuccess = `checked out ` + fmtDevInfo + ` SN '%s' for ` + fmtHostInfo
)

// Package variables required for operation.
var (
	authSvc service.AuthSvc
	serialSvc service.SerialSvc
	systemLog log.MLogger
	errorLog log.MLogger
)

// Init initializes the package variables required for operation.
func Init(as service.AuthSvc, ss service.SerialSvc, sl, el log.MLogger) {
	authSvc, serialSvc, systemLog, errorLog = as, ss, sl, el
}

// CheckIn records a device checkin.
func CheckIn(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	dev := &usbci.Checkin{}
	dev.HostName, dev.RemoteAddr = vars[`host`], r.RemoteAddr

	w.Header().Set(`Content-Type`, `application/json; charset=UTF8`)

	if err := api.DecodeBody(dev, r); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else if _, err := dev.Create(); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		systemLog.Printf(fmtCheckInSuccess, dev.VendorId,
			dev.ProductId, dev.SerialNum, dev.HostName, dev.RemoteAddr)

		w.WriteHeader(http.StatusCreated)
	}
}

// CheckOut retrieves a serialized device and returns it as a JSON object.
func CheckOut(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	dev := &usbci.Serialized{}
	dev.VendorId, dev.ProductId, dev.SerialNum = vars[`vid`], vars[`pid`], vars[`sn`]

	w.Header().Set(`Content-Type`, `application/json; charset=UTF8`)

	if err := dev.Read(); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if j, err := dev.JSON(); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else {

		systemLog.Printf(fmtCheckOutSuccess, dev.VendorId,
			dev.ProductId, dev.SerialNum, vars[`host`], r.RemoteAddr)

		w.WriteHeader(http.StatusOK)

		if _, err = w.Write(j); err != nil {
			errorLog.Panic(err)
		}
	}
}

// NewSn generates a new serial number for an unserialized device.
func NewSn(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	seq := &cmdb.Sequence{}
	dev := &usbci.SnRequest{}
	dev.HostName, dev.RemoteAddr = vars[`host`], r.RemoteAddr

	w.Header().Set(`Content-Type`, `application/json; charset=UTF8`)

	if err := api.DecodeBody(dev, r); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else if seed, err := seq.Create(); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if sn, err := serialSvc.CreateSerial(dev.ObjectType, seed); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if _, err := dev.CreateWithSn(sn); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if exists := dev.DeviceExists(); exists {

		err := fmt.Errorf(fmtNewSnDuplicate, dev.SerialNum,
			dev.VendorId, dev.ProductId, dev.HostName, dev.RemoteAddr)

		errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusConflict)

	} else {

		systemLog.Printf(fmtNewSnSuccess, dev.SerialNum,
			dev.VendorId, dev.ProductId, dev.HostName, dev.RemoteAddr)

		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(sn); err != nil {
			errorLog.Panic(err)
		}
	}
}

// Audit accepts the results of a device self-audit and stores the results.
func Audit(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	aud := &usbci.Audit{}
        aud.VendorId, aud.ProductId, aud.SerialNum, aud.HostName, aud.RemoteAddr =
		vars[`vid`], vars[`pid`], vars[`sn`], vars[`host`], r.RemoteAddr

	w.Header().Set(`Content-Type`, `application/json; charset=UTF8`)

	if body, err := api.ReadBody(r); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return

	} else {

		aud.Changes = body
	}


	if _, err := aud.Create(); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if changes, err := aud.Expand(); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else if _, err := changes.Create(); err != nil {

		errorLog.Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		systemLog.Printf(fmtAuditSuccess, aud.VendorId,
			aud.ProductId, aud.SerialNum, aud.HostName, aud.RemoteAddr)

		w.WriteHeader(http.StatusCreated)
	}
}
