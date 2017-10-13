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

	for _, col := range qy.Cols[`usbCiInsertCheckin`] {
		vals = append(vals, dev[col])
	}

	if _, err = qy.Stmt[`usbCiInsertCheckin`].Exec(vals...); err != nil {
		el.Print(err)
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

	for _, col := range qy.Cols[`usbCiInsertSnRequest`] {
		vals = append(vals, dev[col])
	}

	if res, err := qy.Stmt[`usbCiInsertSnRequest`].Exec(vals...); err != nil {
		el.Print(err)
		return sn, err
	} else if id, err = res.LastInsertId(); err != nil {
		el.Print(err)
		return sn, err
	}

	if res, err := qy.Stmt[`cmdbInsertSequence`].Exec(); err != nil {
		el.Print(err)
		return sn, err
	} else if sq, err := res.LastInsertId(); err != nil {
		el.Print(err)
		return sn, err
	} else {
		sn = fmt.Sprintf(conf.Options.SerialFormat, sq)
	}

	if _, err = qy.Stmt[`usbCiUpdateSnRequest`].Exec(sn, id); err != nil {
		el.Print(err)
	}

	if err != nil {
		el.Print(err)
	} else if err := tx.Commit(); err != nil {
		el.Print(err)
	}

	return sn, err
}

// SaveDeviceChanges records changes reported in a device audit in the 'changes'
// table in JSON format.
func SaveDeviceChanges(host, vid, pid, sn string, chgs []byte) (err error) {

	if _, err = qy.Stmt[`usbCiInsertChanges`].Exec(host, vid, pid, sn, chgs); err != nil {
		el.Print(err)
	}

	return err
}

// GetDeviceJSONObject retreives device properties from the 'serialized' device
// table and returns them to the caller in JSON format.
func GetDeviceJSONObject(vid, pid, sn string) (j []byte, err error) {

	if err = qy.Stmt[`usbCiSelectJSONObject`].QueryRow(vid, pid, sn).Scan(&j); err != nil {
		el.Print(err)
	}

	return j, err
}

// SaveUsbMeta updates the USB meta tables in the database.
func SaveUsbMeta() error {

	tx, err := db.Begin()

	if err != nil {
		el.Print(err)
		return err
	}

	vendorStmt := tx.Stmt(qy.Stmt[`usbMetaReplaceVendor`])
	productStmt := tx.Stmt(qy.Stmt[`usbMetaReplaceProduct`])
	classStmt := tx.Stmt(qy.Stmt[`usbMetaReplaceClass`])
	subclassStmt := tx.Stmt(qy.Stmt[`usbMetaReplaceSubclass`])
	protocolStmt := tx.Stmt(qy.Stmt[`usbMetaReplaceProtocol`])

	VendorLoop:
	for vid, v := range conf.MetaCi.Usb.Vendors {

		if _, err = vendorStmt.Exec(vid, v.String()); err != nil {
			break VendorLoop
		}

		for pid, p := range v.Product {

			if _, err = productStmt.Exec(vid, pid, p.String()); err != nil {
				break VendorLoop
			}
		}
	}

	ClassLoop:
	for cid, c := range conf.MetaCi.Usb.Classes {

		if _, err = classStmt.Exec(cid, c.String()); err != nil {
			break ClassLoop
		}

		for sid, s := range c.Subclass {

			if _, err = subclassStmt.Exec(cid, sid, s.String()); err != nil {
				break ClassLoop
			}

			for pid, p := range s.Protocol {

				if _, err = protocolStmt.Exec(cid, sid, pid, p.String()); err != nil {
					break ClassLoop
				}
			}
		}
	}

	if err != nil {
		el.Print(err)
	} else if err := tx.Commit(); err != nil {
		el.Print(err)
	}

	return err
}
