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

import(
	"database/sql"
	"github.com/jscherff/gocmdb/usbci"
)

// storeAudit stores the results of a device self-audit in the audit table.
func storeAudit(stmt *sql.Stmt, dev *usbci.WSAPI) (err error) {

	var tx *sql.Tx

	if tx, err = conf.Database.Begin(); err == nil {

		for _, ch := range dev.GetChanges() {

			_, err = tx.Stmt(stmt).Exec(
				dev.GetHostName(),
				dev.GetVendorID(),
				dev.GetProductID(),
				dev.GetSerialNum(),
				dev.GetBusNumber(),
				dev.GetBusAddress(),
				dev.GetPortNumber(),
				ch[usbci.FieldNameIx],
				ch[usbci.OldValueIx],
				ch[usbci.NewValueIx],
			)
			if err != nil {
				break
			}
		}
	}

	if err != nil {
		// Decorate error here because rollback result overwrites it.
		conf.Log.Writer[Error].WriteError(ErrorDecorator(err))
		err = tx.Rollback()
	} else {
		err = tx.Commit()
	}

	return err
}

// StoreDevice stores the the device in the table referred to by the statement
// and returns the LAST_INSERT_ID().
func storeDevice(stmt *sql.Stmt, dev *usbci.WSAPI) (id int64, err error) {

	err = stmt.QueryRow(
		dev.GetHostName(),
		dev.GetVendorID(),
		dev.GetProductID(),
		dev.GetSerialNum(),
		dev.GetVendorName(),
		dev.GetProductName(),
		dev.GetProductVer(),
		dev.GetSoftwareID(),
		dev.GetBufferSize(),
		dev.GetBusNumber(),
		dev.GetBusAddress(),
		dev.GetPortNumber(),
		dev.GetUSBSpec(),
		dev.GetUSBClass(),
		dev.GetUSBSubclass(),
		dev.GetUSBProtocol(),
		dev.GetDeviceSpeed(),
		dev.GetDeviceVer(),
		dev.GetMaxPktSize(),
		dev.GetDeviceSN(),
		dev.GetFactorySN(),
		dev.GetDescriptorSN(),
		dev.GetObjectType(),
	).
		Scan(&id)

	return id, err
}

// updateSerial updates the serial number request record with the serial number
// issued.
func updateSerial(stmt *sql.Stmt, sn string, id int64) (res sql.Result, err error) {
	res, err = stmt.Exec(sn, id)
	return res, err
}
