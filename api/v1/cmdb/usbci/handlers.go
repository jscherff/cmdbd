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
	`github.com/jscherff/gox/log`
)

const HttpBodySizeLimit = 1048576

var errLog, sysLog log.MLogger

func SetLoggers(eLog, sLog log.MLogger) {
	errLog, sysLog = eLog, sLog
}

// Checkin records a device checkin.
func CheckinV1(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var host, vid, pid = vars[`host`], vars[`pid`], vars[`vid`]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, HttpBodySizeLimit))

	if err != nil {
		errLog.Panic(err)
	}

	if err = r.Body.Close(); err != nil {
		errLog.Panic(err)
	}

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	dev := make(map[string]interface{})

	if err = json.Unmarshal(body, &dev); err != nil {

		errLog.Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		if err = json.NewEncoder(w).Encode(err); err != nil {
			errLog.Panic(err)
		}

		return
	}

	dev[`object_json`] = body
	dev[`remote_addr`] = r.RemoteAddr

	if err = SaveDeviceCheckin(dev); err != nil {

		errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		sysLog.Printf(`saved checkin for host %q device VID %q PID %q`, host, vid, pid)
		w.WriteHeader(http.StatusCreated)
	}
}

// NewSn generates a new serial number for an unserialized device.
func NewSnV1(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var host, vid, pid = vars[`host`], vars[`pid`], vars[`vid`]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, HttpBodySizeLimit))

	if err != nil {
		errLog.Panic(err)
	}

	if err = r.Body.Close(); err != nil {
		errLog.Panic(err)
	}

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	dev := make(map[string]interface{})

	if err = json.Unmarshal(body, &dev); err != nil {

		errLog.Print(err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		if err = json.NewEncoder(w).Encode(err); err != nil {
			errLog.Panic(err)
		}

		return
	}

	dev[`object_json`] = body
	dev[`remote_addr`] = r.RemoteAddr

	var sn string = dev[`serial_number`].(string)

	if len(sn) > 0 {
		sysLog.Printf(`host %q device VID %q PID %q SN was already set to %q`, host, vid, pid, sn)
	}

	if sn, err = GetNewSerialNumber(dev); err != nil {

		errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		sysLog.Printf(`generated SN %q for host %q device VID %q PID %q`, sn, host, vid, pid)
		w.WriteHeader(http.StatusCreated)
	}

	if err = json.NewEncoder(w).Encode(sn); err != nil {
		errLog.Panic(err)
	}
}

// Audit accepts the results of a device self-audit and stores the results.
func AuditV1(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var host, vid, pid, sn = vars[`host`], vars[`vid`], vars[`pid`], vars[`sn`]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, HttpBodySizeLimit))

	if err != nil {
		errLog.Panic(err)
	}

	if err = r.Body.Close(); err != nil {
		errLog.Panic(err)
	}

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if err = SaveDeviceChanges(host, vid, pid, sn, body); err != nil {

		errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		sysLog.Printf(`recorded audit for host %q device VID %q PID %q SN %q`, host, vid, pid, sn)
		w.WriteHeader(http.StatusCreated)
	}
}

// Checkout retrieves a device from the serialized device database as a
// JSON object and returns it to the caller.
func CheckoutV1(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var host, vid, pid, sn = vars[`host`], vars[`vid`], vars[`pid`], vars[`sn`]

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if j, err := GetDeviceJSONObject(vid, pid, sn); err != nil {

		errLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		sysLog.Printf(`found SN %q for host %q device VID %q PID %q`, sn, host, vid, pid)
		w.WriteHeader(http.StatusOK)

		if _, err = w.Write(j); err != nil {
			errLog.Panic(err)
		}
	}
}
