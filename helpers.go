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
	"fmt"
	"path/filepath"
	"runtime"
	"github.com/jscherff/goutils/dbutils"
	"github.com/jscherff/gocmdb/cmapi"
)

// ErrorDecorator prepends function filename, line number, and function name
// to error messages.
func ErrorDecorator(ue error) (de error) {

	var msg string

	pc, file, line, success := runtime.Caller(1)
	function := runtime.FuncForPC(pc)

	if success {
		msg = fmt.Sprintf("%s:%d: %s()", filepath.Base(file), line, function.Name())
	} else {
		msg = "unknown goroutine"
	}

	return fmt.Errorf("%s: %v", msg, ue)
}

// usbChangeInserts stores the results of a device self-audit in the audit table.
func usbChangeInserts(stmt string, dev *cmapi.UsbCi) (err error) {

	var tx *sql.Tx

	if tx, err = db.Begin(); err != nil {
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
			break
		}
	}

	if err == nil {
		err = tx.Commit()
	} else {
		elog.WriteError(ErrorDecorator(err))
		err = tx.Rollback()
	}

	return err
}

// StoreDevice stores the the device in the table referred to by the statement
// and returns the LAST_INSERT_ID().
func usbciInsertDevice(stmt string, dev *cmapi.UsbCi) (res sql.Result, err error) {

	vals, err := dbutils.ObjectDbVals(dev, "db") //TODO: change to ...ValsByCol

	if err != nil {
		fmt.Println(err) //TODO
	}

	res, err = db.Statements[stmt].Exec(vals...)

	return res, err
}

func usbciSelectDevice(stmt string, dev *cmapi.UsbCi) (rows *sql.Rows, err error) {
	return db.Statements[stmt].Query(dev.VID(), dev.PID(), dev.ID())
}

// usbSnRequestUpdate updates the serial number request record with the serial number
// issued.
func usbSnRequestUpdate(stmt string, sn string, id int64) (res sql.Result, err error) {
	res, err = db.Statements[stmt].Exec(sn, id)
	return res, err
}

func RowToMap(rows *sql.Rows) (mss map[string]string, err error) {

	var cols []string

	if cols, err = rows.Columns(); err != nil {
		return nil, err
	}

	for rows.Next() {

		vals := make([]interface{}, len(cols))
		pvals := make([]interface{}, len(cols))

		for i, _ := range vals {
			pvals[i] = &vals[i]
		}

		if err = rows.Scan(pvals...); err != nil {
			return nil, err
		}

		mss = make(map[string]string)

		for i, cn := range cols {
			if b, ok := vals[i].([]byte); ok {
				mss[cn] = string(b)
			} else {
				mss[cn] = fmt.Sprintf("%v", vals[i])
			}
		}
	}

	if rows.Err() != nil {
		err = rows.Err()
	}

	return mss, err
}
