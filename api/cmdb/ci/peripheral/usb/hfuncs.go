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
	`io`
	`io/ioutil`
	`net/http`
	`github.com/gorilla/mux`
	`github.com/jscherff/gox/log`
)

const (
	HttpBodySizeLimit = 1048576
)

type HandlerFuncV1 interface {
	CheckIn(http.ResponseWriter, *http.Request)
	NewSn(http.ResponseWriter, *http.Request)
	Audit(http.ResponseWriter, *http.Request)
	CheckOut(http.ResponseWriter, *http.Request)
}

type handlerFuncV1 struct {
	errLog log.MLogger
	sysLog log.MLogger
}

func NewHandlerFuncV1(errLog, sysLog log.MLogger) HandlerFuncV1 {
	return &handlerFuncV1{
		errLog: errLog,
		sysLog: sysLog,
	}
}

// CheckIn records a device checkin.
func (this *handlerFuncV1) CheckIn(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var host, vid, pid = vars[`host`], vars[`pid`], vars[`vid`]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, HttpBodySizeLimit))

	if err != nil {
		this.errLog.Panic(err)
	}

	if err = r.Body.Close(); err != nil {
		this.errLog.Panic(err)
	}

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	dev := make(map[string]interface{})

	if err = json.Unmarshal(body, &dev); err != nil {

		this.errLog.Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		if err = json.NewEncoder(w).Encode(err); err != nil {
			this.errLog.Panic(err)
		}

		return
	}

	dev[`object_json`] = body
	dev[`remote_addr`] = r.RemoteAddr

	if err = SaveDeviceCheckin(dev); err != nil {

		this.errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		this.sysLog.Printf(`saved checkin for host %q device VID %q PID %q`, host, vid, pid)
		w.WriteHeader(http.StatusCreated)
	}
}

// newSn generates a new serial number for an unserialized device.
func (this *handlerFuncV1) NewSn(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var host, vid, pid = vars[`host`], vars[`pid`], vars[`vid`]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, HttpBodySizeLimit))

	if err != nil {
		this.errLog.Panic(err)
	}

	if err = r.Body.Close(); err != nil {
		this.errLog.Panic(err)
	}

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	dev := make(map[string]interface{})

	if err = json.Unmarshal(body, &dev); err != nil {

		this.errLog.Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		if err = json.NewEncoder(w).Encode(err); err != nil {
			this.errLog.Panic(err)
		}

		return
	}

	dev[`object_json`] = body
	dev[`remote_addr`] = r.RemoteAddr

	var sn string = dev[`serial_number`].(string)

	if len(sn) > 0 {
		this.sysLog.Printf(`host %q device VID %q PID %q SN was already set to %q`, host, vid, pid, sn)
	}

	if sn, err = GetNewSerialNumber(dev); err != nil {

		this.errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		this.sysLog.Printf(`generated SN %q for host %q device VID %q PID %q`, sn, host, vid, pid)
		w.WriteHeader(http.StatusCreated)
	}

	if err = json.NewEncoder(w).Encode(sn); err != nil {
		this.errLog.Panic(err)
	}
}

// Audit accepts the results of a device self-audit and stores the results.
func (this *handlerFuncV1) Audit(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var host, vid, pid, sn = vars[`host`], vars[`vid`], vars[`pid`], vars[`sn`]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, HttpBodySizeLimit))

	if err != nil {
		this.errLog.Panic(err)
	}

	if err = r.Body.Close(); err != nil {
		this.errLog.Panic(err)
	}

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if err = SaveDeviceChanges(host, vid, pid, sn, body); err != nil {

		this.errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		this.sysLog.Printf(`recorded audit for host %q device VID %q PID %q SN %q`, host, vid, pid, sn)
		w.WriteHeader(http.StatusCreated)
	}
}

// CheckOut retrieves a device from the serialized device database as a
// JSON object and returns it to the caller.
func (this *handlerFuncV1) CheckOut(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var host, vid, pid, sn = vars[`host`], vars[`vid`], vars[`pid`], vars[`sn`]

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if j, err := GetDeviceJSONObject(vid, pid, sn); err != nil {

		this.errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		this.sysLog.Printf(`found SN %q for host %q device VID %q PID %q`, sn, host, vid, pid)
		w.WriteHeader(http.StatusOK)

		if _, err = w.Write(j); err != nil {
			this.errLog.Panic(err)
		}
	}
}
