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
	`database/sql`
	`encoding/json`
	`fmt`
	`io`
	`io/ioutil`
	`net/http`
	`time`
	`github.com/gorilla/mux`
	`github.com/jscherff/gocmdb/cmapi`
)

var HandlerFuncs = map[string]http.HandlerFunc {
	`usbciAction`: usbciAction,
	`usbciAudit`: usbciAudit,
}

// usbciAction handles various 'actions' for device gocmdb agents.
func usbciAction(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	vars := mux.Vars(r)
	var action = vars[`action`]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, ws.HttpBodySizeLimit))

	if err != nil {
		elog.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = r.Body.Close(); err != nil {
		elog.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dev := cmapi.NewUsbCi()

	if err = json.Unmarshal(body, &dev); err != nil {

		elog.Println(err.Error())
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		if err = json.NewEncoder(w).Encode(err); err != nil {
			elog.Println(err.Error())
		}

		return
	}

	switch action {

	case `fetchsn`:

		var sn = dev.ID()

		if len(sn) != 0 {
			w.WriteHeader(http.StatusNoContent)
			err = fmt.Errorf(`serial number already set to %q`, sn)
			elog.Println(err.Error())
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			break
		}

		var id int64
		var res sql.Result

		if res, err = usbciDeviceInsert(`usbciSnRequestInsert`, dev); err != nil {
			// Error already decorated and logged.
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		if id, err = res.LastInsertId(); err != nil {
			elog.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		sn = fmt.Sprintf(`24F%04X`, id)

		if _, err = usbciSnRequestUpdate(`usbciSnRequestUpdate`, sn, id); err != nil {
			// Error already decorated and logged.
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		if err = json.NewEncoder(w).Encode(sn); err != nil {
			elog.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		w.WriteHeader(http.StatusCreated)

	case `checkin`:

		if _, err = usbciDeviceInsert(`usbciCheckinInsert`, dev); err != nil {
			// Error already decorated and logged.
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		w.WriteHeader(http.StatusCreated)

	default:
		err = fmt.Errorf(`unsupported action %q`, action)
		elog.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// usbciAudit performs a device audit against the previous state in
// the database.
func usbciAudit(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(`Content-Type`, `applicaiton/json; charset=UTF8`)

	vars := mux.Vars(r)
	var vid, pid, sn = vars[`vid`], vars[`pid`], vars[`sn`]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, ws.HttpBodySizeLimit))

	if err != nil {
		elog.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = r.Body.Close(); err != nil {
		elog.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dev := cmapi.NewUsbCi()

	if err = json.Unmarshal(body, &dev); err != nil {

		elog.Println(err.Error())
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		if err = json.NewEncoder(w).Encode(err); err != nil {
			elog.Println(err.Error())
		}

		return
	}

	slog.Printf(`auditing VID %q PID %q SN %q`, vid, pid, sn)

	// Retrieve map of device properties from previous checkin, if any.

	slog.Printf(`retrieving previous VID %q PID %q SN %q`, vid, pid, sn)

	map1, err := RowToMap(`usbciAuditSelect`, vid, pid, sn)

	if map1 != nil {
		slog.Printf(`previous VID %q PID %q SN %q found`, vid, pid, sn)
	} else {
		slog.Printf(`previous VID %q PID %q SN %q not found`, vid, pid, sn)
	}

	// Perform a new device checkin to save current device properties. Abort
	// with error if unable to save properties.

	slog.Printf(`storing current VID %q PID %q SN %q`, vid, pid, sn)

	if _, err = usbciDeviceInsert(`usbciCheckinInsert`, dev); err != nil {
		// Error already decorated and logged.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If attempt to retrieve device properties from previous checkin failed
	// or produced empty results, abort without error. (Could be a new device
	// with no previous checkin.)

	if map1 == nil {
		// Not found status already logged
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Retrieve map of device properties from current checkin. Abort with error
	// if unable to retrieve properties or unable to convert to map.

	map2, err := RowToMap(`usbciAuditSelect`, vid, pid, sn)

	if err != nil {
		// Error already decorated and logged.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Compare maps of previous and current properties and record differences in
	// current devices 'Changes' field.

	for cn, _ := range map2 {
		if map1[cn] != map2[cn] {
			dev.Changes = append(dev.Changes, []string{cn, map1[cn], map2[cn]})
		}
	}

	// If there are no differences, return status of 'Not Modified' to caller.

	if len(dev.Changes) == 0 {
		slog.Printf(`device VID %q PID %q SN %q not modified`, vid, pid, sn)
		w.WriteHeader(http.StatusNotModified)
		return
	}

	// Record changes, if any, to changes table. Return status of 'Accepted' if
	// successful, or record error and return 'Internal Server Error' on failure.

	if err = usbciChangeInserts(`usbciChangeInsert`, dev); err != nil {
		// Error already decorated and logged.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare client response by generating a change log suitable for writing
	// to a flat file or CSV file.

	changes := make([][]string, len(dev.Changes))

	for i, change := range(dev.Changes) {
		changes[i] = append(changes[i], time.Now().Local().String())
		changes[i] = append(changes[i], vid, pid, sn)
		changes[i] = append(changes[i], change...)
	}

	// Send the results in a JSON object.

	if err = json.NewEncoder(w).Encode(changes); err != nil {
		elog.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
