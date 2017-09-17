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
	`fmt`
	`github.com/jscherff/gocmdb/cmapi`
	`github.com/jscherff/goutil`
)

// usbciChangeInserts stores the results of a device self-audit in the audit table.
func usbciChangeInserts(stmt string, dev *cmapi.UsbCi) (err error) {

	var tx *sql.Tx

	if tx, err = db.Begin(); err != nil {
		elog.WriteError(goutil.ErrorDecorator(err))
		return err
	}

	for _, ch := range dev.GetChanges() {

		_, err = tx.Stmt(db.Statements[stmt]).Exec(
			dev.GetHostName(),
			dev.GetVendorID(),
			dev.GetProductID(),
			dev.GetSerialNum(),
			dev.GetBusNumber(),
			dev.GetBusAddress(),
			dev.GetPortNumber(),
			ch[cmapi.FieldNameIx],
			ch[cmapi.OldValueIx],
			ch[cmapi.NewValueIx],
		)
		if err != nil {
			elog.WriteError(goutil.ErrorDecorator(err))
			break
		}
	}

	if err == nil {
		err = tx.Commit()
	} else {
		err = tx.Rollback()
	}

	if err != nil {
		elog.WriteError(goutil.ErrorDecorator(err))
	}

	return err
}

// StoreDevice stores the the device in the table referred to by the statement.
func usbciDeviceInsert(stmt string, dev *cmapi.UsbCi) (res sql.Result, err error) {

	vals, err := goutil.ObjectDbValsByCol(dev, `db`, db.Columns[stmt])

	if err == nil {
		res, err = db.Statements[stmt].Exec(vals...)
	}

	if err != nil {
		elog.WriteError(goutil.ErrorDecorator(err))
	}

	return res, err
}

// usbciDeviceSelect retrieves the device from the table referred to by the statement.
func usbciDeviceSelect(stmt string, args ...interface{}) (rows *sql.Rows, err error) {

	if rows, err = db.Statements[stmt].Query(args...); err != nil {
		elog.WriteError(goutil.ErrorDecorator(err))
	}

	return rows, err
}

// usbciSnRequestUpdate updates the serial number request with the serial number issued.
func usbciSnRequestUpdate(stmt string, sn string, id int64) (res sql.Result, err error) {

	if res, err = db.Statements[stmt].Exec(sn, id); err != nil {
		elog.WriteError(goutil.ErrorDecorator(err))
	}

	return res, err
}

// RowToMap converts a database row into a map of string values indexed by column name.
func RowToMap(stmt, vid, pid, id string) (mss map[string]string, err error) {

	rows, err := usbciDeviceSelect(stmt, vid, pid, id)
	defer rows.Close()

	if err != nil {
		elog.WriteError(goutil.ErrorDecorator(err))
		return nil, err
	}

	var cols []string

	if cols, err = rows.Columns(); err != nil {
		elog.WriteError(goutil.ErrorDecorator(err))
		return nil, err
	}

	for rows.Next() {

		vals := make([]interface{}, len(cols))
		pvals := make([]interface{}, len(cols))

		for i, _ := range vals {
			pvals[i] = &vals[i]
		}

		if err = rows.Scan(pvals...); err != nil {
			elog.WriteError(goutil.ErrorDecorator(err))
			return nil, err
		}

		mss = make(map[string]string)

		for i, cn := range cols {
			if b, ok := vals[i].([]byte); ok {
				mss[cn] = string(b)
			} else {
				mss[cn] = fmt.Sprintf(`%v`, vals[i])
			}
		}
	}

	if rows.Err() != nil {
		err = rows.Err()
		elog.WriteError(goutil.ErrorDecorator(err))
	}

	return mss, err
}
