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
	"io/ioutil"
	"net/http"
	"fmt"
	"io"

	"github.com/jscherff/gocmdb/webapi"
	"github.com/gorilla/mux"
)

const sizeLimit int64 = 1048576

func Serial(w http.ResponseWriter, r *http.Request) {

	// Need object type because different types of devices may 
	// have different formats and series of serial numbers.

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

	if err == nil {

		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(device); err != nil {
			panic(err)
		}

	} else {

		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Checkin(w http.ResponseWriter, r *http.Request) {

	// Need object type in order to instantiate the correct
	// object from the 'gocmdb' package.

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
		device.HostName,
		device.VendorID,
		device.ProductID,
		device.VendorName,
		device.ProductName,
		device.ProductVer,
		device.SoftwareID,
		objectType,
	)

	fmt.Println(err)

	if err == nil {

		w.WriteHeader(http.StatusAccepted)

	} else {

		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Audit(w http.ResponseWriter, r *http.Request) {

	// Need only serial number, not object type, because method
	// will only log changes in the form {date, name, old, new}
	// associated with a device serial number. Serial number can
	// be matched to registration record.

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

	fmt.Println(changes, serialNum)	//TODO: record changes to database

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")
	w.WriteHeader(http.StatusAccepted)
}
