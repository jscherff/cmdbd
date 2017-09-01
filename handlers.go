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

	"github.com/jscherff/gocmdb/webapi"
	"github.com/gorilla/mux"
)

const sizeLimit int64 = 1048576

// Serial creates a new record in the 'serials' table when a device
// requests a serial number. It generates a new device serial number
// based on the INT primary key of the table, offers it to the device,
// then updates the 'serial_number' column of the table with the new
// serial number.
func Serial(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	objectType := vars["objectType"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, sizeLimit))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	device := new(webapi.Device)

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")

	if err := json.Unmarshal(body, &device); err != nil {

		w.WriteHeader(http.StatusUnprocessableEntity)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	if len(device.SerialNum) != 0 {

		w.WriteHeader(http.StatusNoContent)
		return
	}

	var insertId int64

	result, err := serialInsertStmt.Exec(
		device.HostName,
		device.VendorID,
		device.ProductID,
		device.VendorName,
		device.ProductName,
		device.ProductVer,
		device.SoftwareID,
		objectType,
	)

	if err == nil {
		insertId, err = result.LastInsertId()
	}

	if err == nil {
		device.SerialNum = fmt.Sprintf("24F%04x", insertId)
		result, err = serialUpdateStmt.Exec(device.SerialNum, insertId)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(device); err != nil {
			panic(err)
		}
	}
}

// Checkin creates a new record in the 'checkin' table when a device
// checks in. A DB trigger then creates a new record in the 'devices'
// table if one does not exist or updates the existing record with data
// from every column except the serial number. The trigger also updates
// the 'last_seen' column of the 'devices' table with the checkin date.
func Checkin(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	objectType := vars["objectType"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, sizeLimit))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	device := new(webapi.Device)

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")

	if err = json.Unmarshal(body, &device); err != nil {

		w.WriteHeader(http.StatusUnprocessableEntity)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}

		return
	}
fmt.Println(device)
	_, err = checkinInsertStmt.Exec(
		device.HostName,
		device.VendorID,
		device.ProductID,
		device.SerialNum,
		device.VendorName,
		device.ProductName,
		device.ProductVer,
		device.SoftwareID,
		objectType,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

// Audit records property changes reported by the device in the 'audits'
// table. Each report is associated with a single serial number but may
// contain multiple changes.
func Audit(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	serialNum := vars["serialNum"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, sizeLimit))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	changes := new(webapi.Changes)

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")

	if err := json.Unmarshal(body, &changes); err != nil {

		w.WriteHeader(http.StatusUnprocessableEntity)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}

		return
	}
/*
	for _, c := range changes {

	_, err = auditInsertStmt.Exec(
		device.HostName,
		device.VendorID,
		device.ProductID,
		device.SerialNum,
		device.VendorName,
		device.ProductName,
		device.ProductVer,
		device.SoftwareID,
		objectType,
		device.HostName,
		device.VendorID,
		device.ProductID,
		device.VendorName,
		device.ProductName,
		device.ProductVer,
		device.SoftwareID,
		objectType,
*/
	fmt.Println(changes, serialNum)	//TODO: record changes to database

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")
	w.WriteHeader(http.StatusAccepted)
}
