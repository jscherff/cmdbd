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

package main

import (
	`encoding/json`
	`io`
	`io/ioutil`
	`net/http`
	`github.com/gorilla/mux`
)

// usbciCheckin records a device checkin.
func v1usbciCheckin(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var host, vid, pid = vars[`host`], vars[`pid`], vars[`vid`]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, ws.HttpBodySizeLimit))

	if err != nil {
		elog.Panic(err)
	}

	if err = r.Body.Close(); err != nil {
		elog.Panic(err)
	}

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	dev := make(map[string]interface{})

	if err = json.Unmarshal(body, &dev); err != nil {

		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		elog.Print(err)

		if err = json.NewEncoder(w).Encode(err); err != nil {
			elog.Panic(err)
		}

		return
	}

	dev[`object_json`] = body
	dev[`remote_addr`] = r.RemoteAddr

	if err = SaveDeviceCheckin(dev); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		slog.Printf(`saved checkin for %q device VID %q PID %q`, host, vid, pid)
		w.WriteHeader(http.StatusCreated)
	}
}

// usbciNewSN generates a new serial number for an unserialized device.
func v1usbciNewSN(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var host, vid, pid = vars[`host`], vars[`pid`], vars[`vid`]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, ws.HttpBodySizeLimit))

	if err != nil {
		elog.Panic(err)
	}

	if err = r.Body.Close(); err != nil {
		elog.Panic(err)
	}

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	dev := make(map[string]interface{})

	if err = json.Unmarshal(body, &dev); err != nil {

		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		elog.Print(err)

		if err = json.NewEncoder(w).Encode(err); err != nil {
			elog.Panic(err)
		}

		return
	}

	dev[`object_json`] = body
	dev[`remote_addr`] = r.RemoteAddr

	var sn string = dev[`serial_number`].(string)

	if len(sn) > 0 {
		slog.Printf(`serial number was already set to SN %q`, sn)
	}

	if sn, err = GetNewSerialNumber(dev); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		slog.Printf(`generated SN %q for %q device VID %q PID %q`, sn, host, vid, pid)
		w.WriteHeader(http.StatusCreated)
	}

	if err = json.NewEncoder(w).Encode(sn); err != nil {
		elog.Panic(err)
	}
}

// usbciAudit accepts the results of a device self-audit and stores the results.
func v1usbciAudit(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var host, vid, pid, sn = vars[`host`], vars[`vid`], vars[`pid`], vars[`sn`]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, ws.HttpBodySizeLimit))

	if err != nil {
		elog.Panic(err)
	}

	if err = r.Body.Close(); err != nil {
		elog.Panic(err)
	}

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if err = SaveDeviceChanges(host, vid, pid, sn, body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		slog.Printf(`recorded audit for %q device VID %q PID %q SN %q`, host, vid, pid, sn)
		w.WriteHeader(http.StatusCreated)
	}
}

// usbciCheckout retrieves a device from the serialized device database as a
// JSON object and returns it to the caller.
func v1usbciCheckout(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var host, vid, pid, sn = vars[`host`], vars[`vid`], vars[`pid`], vars[`sn`]

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	if j, err := GetDeviceJSONObject(vid, pid, sn); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		slog.Printf(`found SN %q for %q device VID %q PID %q`, sn, host, vid, pid)
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(j); err != nil {
			elog.Panic(err)
		}
	}
}
