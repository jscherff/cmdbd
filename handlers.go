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
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jscherff/gocmdb/cmapi"
	"github.com/jscherff/goutils"
)

var HandlerFuncs = map[string]http.HandlerFunc {
	"usbciAction": usbciAction,
	"usbciAudit": usbciAudit,
}

// usbciAction handles various 'actions' for device gocmdb agents.
func usbciAction(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")

	vars := mux.Vars(r)
	var action = vars["action"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, ws.HttpBodySizeLimit))

	if err != nil {
		panic(goutils.ErrorDecorator(err))
	}

	if err = r.Body.Close(); err != nil {
		panic(goutils.ErrorDecorator(err))
	}

	dev := cmapi.NewUsbCi()

	if err = json.Unmarshal(body, &dev); err != nil {

		elog.WriteError(goutils.ErrorDecorator(err))
		w.WriteHeader(http.StatusUnprocessableEntity)

		if err = json.NewEncoder(w).Encode(err); err != nil {
			panic(goutils.ErrorDecorator(err))
		}

		return
	}

	switch action {

	case "fetchsn":

		var sn = dev.ID()

		if len(sn) != 0 {
			w.WriteHeader(http.StatusNoContent)
			break
		}

		var id int64
		var res sql.Result

		if res, err = usbciDeviceInsert("usbciSnRequestInsert", dev); err != nil {
			break
		}

		if id, err = res.LastInsertId(); err != nil {
			elog.WriteError(goutils.ErrorDecorator(err))
			break
		}

		sn = fmt.Sprintf("24F%04X", id)

		if _, err = usbciSnRequestUpdate("usbciSnRequestUpdate", sn, id); err != nil {
			break
		}

		w.WriteHeader(http.StatusCreated)

		if err = json.NewEncoder(w).Encode(sn); err != nil {
			panic(goutils.ErrorDecorator(err))
		}

	case "checkin":

		if _, err = usbciDeviceInsert("usbciCheckinInsert", dev); err == nil {
			w.WriteHeader(http.StatusAccepted)
		}

	case "changes":

		usbciDeviceInsert("usbciCheckinInsert", dev)

		if len(dev.Changes) == 0 {
			w.WriteHeader(http.StatusNoContent)
			break
		}

		if err = usbciChangeInserts("usbciChangeInsert", dev); err == nil {
			w.WriteHeader(http.StatusAccepted)
		}

	default:

		w.WriteHeader(http.StatusBadRequest)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// usbciAudit performs a device audit against the previous state in
// the database.
func usbciAudit(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")

	vars := mux.Vars(r)
	var vid, pid, id = vars["vid"], vars["pid"], vars["id"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, ws.HttpBodySizeLimit))

	if err != nil {
		panic(goutils.ErrorDecorator(err))
	}

	if err = r.Body.Close(); err != nil {
		panic(goutils.ErrorDecorator(err))
	}

	dev := cmapi.NewUsbCi()

	if err = json.Unmarshal(body, &dev); err != nil {

		elog.WriteError(goutils.ErrorDecorator(err))
		w.WriteHeader(http.StatusUnprocessableEntity)

		if err = json.NewEncoder(w).Encode(err); err != nil {
			panic(goutils.ErrorDecorator(err))
		}

		return
	}

	// Retrieve map of device properties from previous checkin, if any.

	map1, err := RowToMap("usbciAuditSelect", vid, pid, id)

	// Perform a new device checkin to save current device properties. Abort
	// with error if unable to save properties.

	if _, err = usbciDeviceInsert("usbciCheckinInsert", dev); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If attempt to retrieve device properties from previous checkin failed
	// or produced empty results, abort without error. (Could be a new device
	// with no previous checkin.)

	if map1 == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Retrieve map of device properties from current checkin. Abort with error
	// if unable to retrieve properties or unable to convert to map.

	map2, err := RowToMap("usbciAuditSelect", vid, pid, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Compare maps of previous and current properties and record differences in
	// current devices 'Changes' field.

	for cn, _ := range map2 {
		if map1[cn] != map2[cn] {
			dev.Changes = append(dev.Changes, []string{cn, map1[cn], map2[cn]})
		}
	}

	// If there are no differences, return status of 'No Content' to caller.

	if len(dev.Changes) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Record changes, if any, to changes table. Return status of 'Accepted' if
	// successful, or record error and return 'Internal Server Error' on failure.

	if err = usbciChangeInserts("usbciChangeInsert", dev); err == nil {
		w.WriteHeader(http.StatusAccepted)
	} else {
		elog.WriteError(goutils.ErrorDecorator(err))
	}
}
