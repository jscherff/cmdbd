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
	`github.com/jscherff/cmdbd/service`
	`github.com/jscherff/cmdbd/model/cmdb/usbci`
	`github.com/jscherff/gox/log`
)

const HttpBodySizeLimit = 1048576

// HandlersV2 contains http.HandleFunc signatures of CMDBd APIv2.
type HandlersV2 interface {
	Checkin(http.ResponseWriter, *http.Request)
	NewSn(http.ResponseWriter, *http.Request)
	Audit(http.ResponseWriter, *http.Request)
	CheckOut(http.ResponseWriter, *http.Request)
}

// handlersV2 implements the HandlersV2 interface.
type handlersV2 struct {
	errorLog log.MLogger
	systemLog log.MLogger
	serialNumSvc service.SerialNumService
}

// NewHandlersV2 returns a new handlersV2 instance.
func NewHandlersV2(errLog, sysLog log.MLogger, snSvc service.SerialNumService) HandlersV2 {
	return &handlersV2{
		errorLog: errLog,
		systemLog: sysLog,
		serialNumSvc: snSvc,
	}
}

// Checkin records a device checkin.
func (this *handlersV2) Checkin(w http.ResponseWriter, r *http.Request) {

	dev := &usbci.Checkin{}
	dev.RemoteAddr = r.RemoteAddr

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if err := this.decode(dev, w, r); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else if _, err := dev.Create(); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		this.systemLog.Printf(`checked in USB device %s %s SN %q on host %s`,
			dev.VendorID, dev.ProductID, dev.SerialNum, dev.HostName)

		w.WriteHeader(http.StatusCreated)
	}
}

// NewSn generates a new serial number for an unserialized device.
func (this *handlersV2) NewSn(w http.ResponseWriter, r *http.Request) {

	dev := &usbci.SnRequest{}
	dev.RemoteAddr = r.RemoteAddr

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if err := this.decode(dev, w, r); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else if id, err := dev.Create(); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if sn, err := this.serialNumSvc.Create(dev.ObjectType, id); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		this.systemLog.Printf(`issued SN %q for USB device %s %s on host %s`,
			dev.SerialNum, dev.VendorID, dev.ProductID, dev.HostName)

		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(sn); err != nil {
			this.errorLog.Panic(err)
		}
	}
}

// Audit accepts the results of a device self-audit and stores the results.
func (this *handlersV2) Audit(w http.ResponseWriter, r *http.Request) {

	dev := &usbci.Audit{}
	dev.RemoteAddr = r.RemoteAddr

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if err := this.decode(dev, w, r); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else if _, err := dev.Create(); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		this.systemLog.Printf(`audited USB device %s %s SN %q on host %s`,
			dev.VendorID, dev.ProductID, dev.SerialNum, dev.HostName)

		w.WriteHeader(http.StatusCreated)
	}
}

// CheckOut retrieves a serialized device and returns it as a JSON object.
func (this *handlersV2) CheckOut(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	dev := &usbci.Serialized{}
	dev.SerialNum = vars[`sn`]

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if err := dev.Read(dev); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if j, err := dev.JSON(); err != nil {

		this.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

	} else {

		this.systemLog.Printf(`checked out USB device %s %s SN %q for host %s`,
			dev.VendorID, dev.ProductID, dev.SerialNum, vars[`host`])

		w.WriteHeader(http.StatusOK)

		if _, err = w.Write(j); err != nil {
			this.errorLog.Panic(err)
		}
	}
}

// decode unmarshals the JSON object in the HTTP request body to an object.
func (this *handlersV2) decode(i interface{}, w http.ResponseWriter, r *http.Request) (error) {

	if body, err := ioutil.ReadAll(io.LimitReader(r.Body, HttpBodySizeLimit)); err != nil {
		this.errorLog.Panic(err)
	} else if err := r.Body.Close(); err != nil {
		this.errorLog.Panic(err)
	} else if err := json.Unmarshal(body, &i); err != nil {
		return err
	}

	return nil
}
