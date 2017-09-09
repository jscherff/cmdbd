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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/jscherff/gocmdb/usbci"
)

// SerialHandler creates a new record in the 'serials' table when a device
// requests a serial number. It generates a new device serial number based
// on the INT primary key of the table, offers it to the device, then updates
// the 'serial_number' column of the table with the new serial number.
func SerialHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, HttpBodySizeLimit))

	if err != nil {
		errorLog.WriteError(ErrorDecorator(err))
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		errorLog.WriteError(ErrorDecorator(err))
		panic(err)
	}

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")

	dev := usbci.NewWSAPI()

	if err := json.Unmarshal(body, &dev); err != nil {

		errorLog.WriteError(ErrorDecorator(err))
		w.WriteHeader(http.StatusUnprocessableEntity)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			errorLog.WriteError(ErrorDecorator(err))
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")

	if len(dev.GetSerialNum()) != 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var id int64
	var sn string

	if id, err = storeDevice(db.Stmt.SerialInsert, dev); err != nil {
		errorLog.WriteError(ErrorDecorator(err))
	} else {
		sn = fmt.Sprintf("24F%04x", id)
		_, err = updateSerial(db.Stmt.SerialUpdate, sn, id)
	}

	if err != nil {
		errorLog.WriteError(ErrorDecorator(err))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(sn); err != nil {
			errorLog.WriteError(ErrorDecorator(err))
			panic(err)
		}
	}
}

// CheckinHandler creates a new record in the 'checkin' table when a device
// checks in. A DB trigger then creates a new record in the 'devices' table
// if one does not exist or updates the existing record with data from every
// column except the serial number. The trigger also updates the 'last_seen'
// column of the 'devices' table with the checkin date.
func CheckinHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, HttpBodySizeLimit))

	if err != nil {
		errorLog.WriteError(ErrorDecorator(err))
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		errorLog.WriteError(ErrorDecorator(err))
		panic(err)
	}

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")

	dev := usbci.NewWSAPI()

	if err = json.Unmarshal(body, &dev); err != nil {

		errorLog.WriteError(ErrorDecorator(err))
		w.WriteHeader(http.StatusUnprocessableEntity)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			errorLog.WriteError(ErrorDecorator(err))
			panic(err)
		}

		return
	}

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")

	if _, err = storeDevice(db.Stmt.CheckinInsert, dev); err != nil {
		errorLog.WriteError(ErrorDecorator(err))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

// AuditHandler records property changes reported by the device in the 'audits'
// table. Each report is associated with a single serial number but may contain
// multiple changes.
func AuditHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, HttpBodySizeLimit))

	if err != nil {
		errorLog.WriteError(ErrorDecorator(err))
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		errorLog.WriteError(ErrorDecorator(err))
		panic(err)
	}

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")

	dev := usbci.NewWSAPI()

	if err := json.Unmarshal(body, &dev); err != nil {

		errorLog.WriteError(ErrorDecorator(err))
		w.WriteHeader(http.StatusUnprocessableEntity)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			errorLog.WriteError(ErrorDecorator(err))
			panic(err)
		}

		return
	}

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")

	if _, err = storeDevice(db.Stmt.CheckinInsert, dev); err != nil {
		errorLog.WriteError(ErrorDecorator(err))
	}

	if err = storeAudit(db.Stmt.AuditInsert, dev); err != nil {
		errorLog.WriteError(ErrorDecorator(err))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}
