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
	`net/http`
	`github.com/gorilla/mux`
	`github.com/jscherff/cmdbd/api`
	`github.com/jscherff/cmdbd/api/v2/cmdb/usbci`
	`github.com/jscherff/cmdbd/service`
)

// Package variables required for operation.
var (
	loggerSvc service.LoggerSvc
)

// Init initializes the package variables required for operation.
func Init(ls service.LoggerSvc) {
	loggerSvc = ls
}

// Audit accepts the results of a device self-audit and stores the results.
func Audit(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	dev := struct {
		VendorId	string	`json:"vendor_id"`
		ProductId	string	`json:"product_id"`
		SerialNum	string	`json:"serial_number"`
		HostName	string	`json:"host_name"`
		Changes		[]byte	`json:"changes"`
	} {
		VendorId:	vars[`vid`],
		ProductId:	vars[`pid`],
		SerialNum:	vars[`sn`],
		HostName:	vars[`host`],
	}

	if body, err := api.ReadBody(r); err != nil {
		loggerSvc.ErrorLog().Panic(err)
	} else {
		dev.Changes = body
	}

	if body, err := json.Marshal(dev); err != nil {
		loggerSvc.ErrorLog().Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		api.WriteBody(r, body)
		usbci.Audit(w, r)
	}
}
