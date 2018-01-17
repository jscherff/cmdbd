
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
)

// Templates for system and error messages.
const (
	fmtCheckInSuccess = `checked in USB device %s-%s SN %q for host %s at %s`
	fmtNewSnDuplicate = `duplicate SN %q for USB device %s-%s on host %s at %s`
	fmtNewSnSuccess = `issued SN '%q' for USB device %s-%s on host %s at %s`
	fmtAuditSuccess = `audited USB device %s-%s SN %q for host %s at %s`
	fmtCheckOutSuccess = `checked out USB device %s-%s SN %q for host %s at %s`
)

// Package variables required for operation.
var (
	authSvc service.AuthSvc
	serialSvc service.SerialSvc
	loggerSvc service.LoggerSvc
)

// Init initializes the package variables required for operation.
func Init(as service.AuthSvc, ss service.SerialSvc, ls service.LoggerSvc) {
	authSvc, serialSvc, loggerSvc = as, ss, ls
}

// CheckIn records a device checkin.
func CheckIn(w http.ResponseWriter, r *http.Request) {

	dev := &usbci.Checkin{}
	dev.RemoteAddr = r.RemoteAddr

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if err := api.DecodeBody(dev, r); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else if _, err := dev.Create(); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		loggerSvc.SystemLog().Printf(fmtCheckInSuccess, dev.VendorId,
			dev.ProductId, dev.SerialNum, dev.HostName, dev.RemoteAddr)

		w.WriteHeader(http.StatusCreated)
	}
}

// NewSn generates a new serial number for an unserialized device.
func NewSn(w http.ResponseWriter, r *http.Request) {

	seq := &cmdb.Sequence{}
	dev := &usbci.SnRequest{}
	dev.RemoteAddr = r.RemoteAddr

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if err := api.DecodeBody(dev, r); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else if seed, err := seq.Create(); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if sn, err := serialSvc.CreateSerial(dev.ObjectType, seed); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if _, err := dev.CreateWithSn(sn); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if exists := dev.DeviceExists(); exists {

		err := fmt.Errorf(fmtNewSnDuplicate, dev.SerialNum,
			dev.VendorId, dev.ProductId, dev.HostName, dev.RemoteAddr)

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusConflict)

	} else {

		loggerSvc.SystemLog().Printf(fmtNewSnSuccess, dev.SerialNum,
			dev.VendorId, dev.ProductId, dev.HostName, dev.RemoteAddr)

		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(sn); err != nil {
			loggerSvc.ErrorLog().Panic(err)
		}
	}
}

// Audit accepts the results of a device self-audit and stores the results.
func Audit(w http.ResponseWriter, r *http.Request) {

	dev := &usbci.Audit{}
	dev.RemoteAddr = r.RemoteAddr

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if err := api.DecodeBody(dev, r); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else if id, err := dev.Create(); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		dev.Id = id
	}

	if err := dev.Read(); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if changes, err := dev.Expand(); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else if _, err := changes.Create(); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		loggerSvc.SystemLog().Printf(fmtAuditSuccess, dev.VendorId,
			dev.ProductId, dev.SerialNum, dev.HostName, dev.RemoteAddr)

		w.WriteHeader(http.StatusCreated)
	}
}

// CheckOut retrieves a serialized device and returns it as a JSON object.
func CheckOut(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	dev := &usbci.Serialized{}
	dev.VendorId, dev.ProductId, dev.SerialNum = vars[`vid`], vars[`pid`], vars[`sn`]

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if err := dev.Read(); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if j, err := dev.JSON(); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else {

		loggerSvc.SystemLog().Printf(fmtCheckOutSuccess, dev.VendorId,
			dev.ProductId, dev.SerialNum, vars[`host`], r.RemoteAddr)

		w.WriteHeader(http.StatusOK)

		if _, err = w.Write(j); err != nil {
			loggerSvc.ErrorLog().Panic(err)
		}
	}
}
