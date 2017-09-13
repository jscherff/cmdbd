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
)

// usbciActionHandler handles various 'actions' for device gocmdb agents.
func usbciActionHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	action := vars["action"]

	body, err := ioutil.ReadAll(io.LimitReader(
		r.Body, conf.Server.HttpBodySizeLimit))

	if err != nil {
		panic(ErrorDecorator(err))
	}

	if err := r.Body.Close(); err != nil {
		panic(ErrorDecorator(err))
	}

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")

	dev := cmapi.NewUsbCi()

	if err := json.Unmarshal(body, &dev); err != nil {
		w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")

		elog.WriteError(ErrorDecorator(err))
		w.WriteHeader(http.StatusUnprocessableEntity)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(ErrorDecorator(err))
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

		if res, err = usbciInsertDevice("usbSnRequestInsert", dev); err != nil {
			break
		}

		if id, err = res.LastInsertId(); err != nil {
			break
		}

		sn = fmt.Sprintf("24F%04X", id)

		if _, err = usbSnRequestUpdate("usbSnRequestUpdate", sn, id); err != nil {
			break
		}

		w.WriteHeader(http.StatusCreated)

		if err = json.NewEncoder(w).Encode(sn); err != nil {
			panic(ErrorDecorator(err))
		}

	case "checkin":

		if _, err = usbciInsertDevice("usbCheckinInsert", dev); err == nil {
			w.WriteHeader(http.StatusAccepted)
		}

	case "changes":

		if _, errC := usbciInsertDevice("usbCheckinInsert", dev); errC != nil {
			elog.WriteError(ErrorDecorator(errC))
		}

		if len(dev.Changes) == 0 {
			w.WriteHeader(http.StatusNoContent)
			break
		}

		if err = usbChangeInserts("usbChangeInsert", dev); err == nil {
			w.WriteHeader(http.StatusAccepted)
		}
	}

	if err != nil {
		elog.WriteError(ErrorDecorator(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// AllowedMethodHandler restricts requests to methods listed in the AllowedMethods
// slice in the systemwide configuration.
func AllowedMethodHandler(h http.Handler, methods ...string) http.Handler {

	return http.HandlerFunc(

		func(w http.ResponseWriter, r *http.Request) {

			for _, m := range methods {
				if r.Method == m {
					h.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, fmt.Sprintf("Unsupported method %q", r.Method),
				http.StatusMethodNotAllowed)
		},
	)
}

// usbciAuditHandler performs a device audit against the previous state in
// the database.
func usbciAuditHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var vid, pid, sn = vars["vid"], vars["pid"], vars["sn"]

	body, err := ioutil.ReadAll(io.LimitReader(
		r.Body, conf.Server.HttpBodySizeLimit))

	if err != nil {
		panic(ErrorDecorator(err))
	}

	if err := r.Body.Close(); err != nil {
		panic(ErrorDecorator(err))
	}

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")

	dev := cmapi.NewUsbCi()

	if err := json.Unmarshal(body, &dev); err != nil {

		elog.WriteError(ErrorDecorator(err))
		w.WriteHeader(http.StatusUnprocessableEntity)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(ErrorDecorator(err))
		}

		return
	}

	if len(vid) == 0 || len(pid) == 0 || len(sn) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retrieve map of device properties from previous checkin, if any.

	var map1 map[string]string

	rows, err := usbciSelectDevice("usbAuditSelect", dev)

	if err == nil {
		defer rows.Close()
		map1, err = RowToMap(rows)
	}
	if err != nil {
		elog.WriteError(ErrorDecorator(err))
	}

	// Perform a new device checkin to save current device properties. Abort
	// with error if unable to save properties.

	if _, errC := usbciInsertDevice("usbCheckinInsert", dev); errC != nil {
		elog.WriteError(ErrorDecorator(errC))
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

	var map2 map[string]string

	rows, err = usbciSelectDevice("usbAuditSelect", dev)

	if err == nil {
		defer rows.Close()
		map2, err = RowToMap(rows)
	}
	if err != nil {
		elog.WriteError(ErrorDecorator(err))
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

	if err = usbChangeInserts("usbChangeInsert", dev); err == nil {
		w.WriteHeader(http.StatusAccepted)
	} else {
		elog.WriteError(ErrorDecorator(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
