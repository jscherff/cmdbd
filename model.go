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

	if _, err = db.Statements[`usbciInsertCheckin`].Exec(vals...); err != nil {
		elog.Print(err)
	}

	return err
}

// GetNewSerialNumber generates a new device serial number using the value
// from the auto-incremented ID column of the 'snrequest' table with the
// format string provided by the caller.
func GetNewSerialNumber(dev map[string]interface{}) (sn string, err error) {

	var (
		id int64
		vals []interface{}
	)

	tx, err := db.Begin()

	if err != nil {
		return sn, err
	}

	for _, col := range db.Columns[`usbciInsertSnRequest`] {
		vals = append(vals, dev[col])
	}

	if res, err := db.Statements[`usbciInsertSnRequest`].Exec(vals...); err != nil {
		elog.Print(err)
		return sn, err
	} else if id, err = res.LastInsertId(); err != nil {
		elog.Print(err)
		return sn, err
	}

	if res, err := db.Statements[`cmdbInsertSequence`].Exec(); err != nil {
		elog.Print(err)
		return sn, err
	} else if sq, err := res.LastInsertId(); err != nil {
		elog.Print(err)
		return sn, err
	} else {
		sn = fmt.Sprintf(conf.Options.SerialFormat, sq)
	}

	if _, err = db.Statements[`usbciUpdateSnRequest`].Exec(sn, id); err != nil {
		elog.Print(err)
	}

	if err != nil {
		elog.Print(err)
	} else if err := tx.Commit(); err != nil {
		elog.Print(err)
	}

	return sn, err
}

// SaveDeviceChanges records changes reported in a device audit in the 'changes'
// table in JSON format.
func SaveDeviceChanges(host, vid, pid, sn string, chgs []byte) (err error) {

	if _, err = db.Statements[`usbciInsertChanges`].Exec(host, vid, pid, sn, chgs); err != nil {
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

// SaveUsbMeta updates the USB meta tables in the database.
func SaveUsbMeta() error {

	tx, err := db.Begin()

	if err != nil {
		elog.Print(err)
		return err
	}

	vendorStmt := tx.Stmt(db.Statements[`metaReplaceUsbVendor`])
	productStmt := tx.Stmt(db.Statements[`metaReplaceUsbProduct`])
	classStmt := tx.Stmt(db.Statements[`metaReplaceUsbClass`])
	subclassStmt := tx.Stmt(db.Statements[`metaReplaceUsbSubclass`])
	protocolStmt := tx.Stmt(db.Statements[`metaReplaceUsbProtocol`])

	VendorLoop:
	for vid, v := range conf.Data.UsbMeta.Vendors {

		svid, sv := vid.String(), v.String()

		if _, err = vendorStmt.Exec(svid, sv); err != nil {
			break VendorLoop
		}

		for pid, p := range v.Product {

			spid, sp := pid.String(), p.String()

			if _, err = productStmt.Exec(svid, spid, sp); err != nil {
				break VendorLoop
			}
		}
	}

	ClassLoop:
	for cid, c := range conf.Data.UsbMeta.Classes {

		scid, sc := fmt.Sprintf(`%02x`, uint8(cid)), c.String()

		if _, err = classStmt.Exec(scid, sc); err != nil {
			break ClassLoop
		}

		for sid, s := range c.Subclass {

			ssid, ss := fmt.Sprintf(`%02x`, uint8(sid)), s.String()

			if _, err = subclassStmt.Exec(scid, ssid, ss); err != nil {
				break ClassLoop
			}

			for pid, p := range s.Protocol {

				spid, sp := fmt.Sprintf(`%02x`, uint8(pid)), p.String()

				if _, err = protocolStmt.Exec(scid, ssid, spid, sp); err != nil {
					break ClassLoop
				}
			}
		}
	}

	if err != nil {
		elog.Print(err)
	} else if err := tx.Commit(); err != nil {
		elog.Print(err)
	}

	return err
}
