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

import `fmt`

// SaveDeviceCheckin saves a device checkin to the database 'checkins' table.
func SaveDeviceCheckin(dev map[string]interface{}) (err error) {

	var vals []interface{}

	for _, col := range db.Columns[`usbciInsertCheckin`] {
		vals = append(vals, dev[col])
	}

	if _, err := db.Statements[`usbciInsertCheckin`].Exec(vals...); err != nil {
		elog.Print(err)
	}

	return err
}

// GetNewSerialNumber generates a new device serial number using the value
// from the auto-incremented ID column of the 'snrequest' table with the
// format string provided by the caller.
func GetNewSerialNumber(sfmt string, dev map[string]interface{}) (sn string, err error) {

	var vals []interface{}

	for _, col := range db.Columns[`usbciInsertSnRequest`] {
		vals = append(vals, dev[col])
	}

	res, err := db.Statements[`usbciInsertSnRequest`].Exec(vals...)

	if err != nil {
		elog.Print(err)
		return sn, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		elog.Print(err)
		return sn, err
	}

	sn = fmt.Sprintf(`24F%04X`, id)

	if _, err = db.Statements[`usbciUpdateSnRequest`].Exec(sn, id); err != nil {
		elog.Print(err)
	}

	return sn, err
}

// SaveDeviceChanges records changes reported in a device audit in the 'changes'
// table in JSON format.
func SaveDeviceChanges(host, vid, pid, sn string, chgs []byte) (err error) {

	if _, err := db.Statements[`usbciInsertChanges`].Exec(host, vid, pid, sn, chgs); err != nil {
		elog.Print(err)
	}

	return err
}

// GetDeviceJSONObject retreives device properties from the 'serialized' device
// table and returns them to the caller in JSON format.
func GetDeviceJSONObject(vid, pid, sn string) (j []byte, err error) {

	if err = db.Statements[`usbciSelectJSONObject`].QueryRow(vid, pid, sn).Scan(&j); err != nil {
		elog.Print(err)
	}
	return j, err
}
