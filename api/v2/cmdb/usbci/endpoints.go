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
	`io`
	`io/ioutil`
	`net/http`
	`github.com/gorilla/mux`
	`github.com/jscherff/cmdbd/api`
	`github.com/jscherff/cmdbd/service`
	`github.com/jscherff/cmdbd/model/cmdb/usbci`
)

const HttpBodySizeLimit = 1048576

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

// Endpoints is a collection of URL path to handler function mappings.
var Endpoints = []api.Endpoint {

	api.Endpoint {
		Name:		`USB CI Checkin Handler`,
		Path:		`/v2/cmdb/ci/usb/checkin/{host}/{vid}/{pid}`,
		Method:		`POST`,
		HandlerFunc:	Checkin,
		Protected:	false,
	},

	api.Endpoint {
		Name:		`USB CI Checkout Handler`,
		Path:		`/v2/cmdb/ci/usb/checkout/{host}/{vid}/{pid}/{sn}`,
		Method:		`GET`,
		HandlerFunc:	CheckOut,
		Protected:	false,
	},

	api.Endpoint {
		Name:		`USB CI NewSn Handler`,
		Path:		`/v2/cmdb/ci/usb/newsn/{host}/{vid}/{pid}`,
		Method:		`POST`,
		HandlerFunc:	NewSn,
		Protected:	false,
	},

	api.Endpoint {
		Name:		`USB CI Audit Handler`,
		Path:		`/v2/cmdb/ci/usb/audit/{host}/{vid}/{pid}/{sn}`,
		Method:		`POST`,
		HandlerFunc:	Audit,
		Protected:	false,
	},
}

// Checkin records a device checkin.
func Checkin(w http.ResponseWriter, r *http.Request) {

	dev := &usbci.Checkin{}
	dev.RemoteAddr = r.RemoteAddr

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if err := decode(dev, w, r); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else if _, err := dev.Create(); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		loggerSvc.SystemLog().Printf(`checked in USB device %s %s SN %q on host %s`,
			dev.VendorID, dev.ProductID, dev.SerialNum, dev.HostName)

		w.WriteHeader(http.StatusCreated)
	}
}

// NewSn generates a new serial number for an unserialized device.
func NewSn(w http.ResponseWriter, r *http.Request) {

	dev := &usbci.SnRequest{}
	dev.RemoteAddr = r.RemoteAddr

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if err := decode(dev, w, r); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else if id, err := dev.Create(); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if sn, err := serialSvc.CreateSerial(dev.ObjectType, id); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		loggerSvc.SystemLog().Printf(`issued SN %q for USB device %s %s on host %s`,
			dev.SerialNum, dev.VendorID, dev.ProductID, dev.HostName)

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

	if err := decode(dev, w, r); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else if _, err := dev.Create(); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		loggerSvc.SystemLog().Printf(`audited USB device %s %s SN %q on host %s`,
			dev.VendorID, dev.ProductID, dev.SerialNum, dev.HostName)

		w.WriteHeader(http.StatusCreated)
	}
}

// CheckOut retrieves a serialized device and returns it as a JSON object.
func CheckOut(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	dev := &usbci.Serialized{}
	dev.SerialNum = vars[`sn`]

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if err := dev.Read(dev); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if j, err := dev.JSON(); err != nil {

		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else {

		loggerSvc.SystemLog().Printf(`checked out USB device %s %s SN %q for host %s`,
			dev.VendorID, dev.ProductID, dev.SerialNum, vars[`host`])

		w.WriteHeader(http.StatusOK)

		if _, err = w.Write(j); err != nil {
			loggerSvc.ErrorLog().Panic(err)
		}
	}
}

// decode unmarshals the JSON object in the HTTP request body to an object.
func decode(i interface{}, w http.ResponseWriter, r *http.Request) (error) {

	if body, err := ioutil.ReadAll(io.LimitReader(r.Body, HttpBodySizeLimit)); err != nil {
		loggerSvc.ErrorLog().Panic(err)
	} else if err := r.Body.Close(); err != nil {
		loggerSvc.ErrorLog().Panic(err)
	} else if err := json.Unmarshal(body, &i); err != nil {
		return err
	}

	return nil
}
