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
	"path/filepath"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

const (
	auditInsertSQL string = `

		INSERT INTO audits (
			serial_number,
			field_name,
			old_value,
			new_value
		)

		VALUES (?, ?, NULLIF(?, ''), NULLIF(?, ''))`

	checkinInsertSQL string = `

		INSERT INTO checkins (
			host_name,
			vendor_id,
			product_id,
			serial_number,
			vendor_name,
			product_name,
			product_ver,
			software_id,
			object_type
		)

		VALUES (?, ?, ?, ?, NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), ?)`

	serialInsertSQL string = `
	
		INSERT INTO serials (
			host_name,
			vendor_id,
			product_id,
			serial_number,
			vendor_name,
			product_name,
			product_ver,
			software_id,
			object_type
		)

		VALUES (?, ?, ?, ?, NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), ?)`

	serialUpdateSQL string = `

		UPDATE serials
		SET serial_number = ?
		WHERE id = ?`
)

type Database struct {

	handle *sql.DB
	config *mysql.Config

	info string

	stmt struct {
		auditInsert	*sql.Stmt
		checkinInsert	*sql.Stmt
		serialInsert	*sql.Stmt
		serialUpdate	*sql.Stmt
	}
}

func NewDatabase(cf string) (this *Database, err error) {

	this = new(Database)
	appDir := filepath.Dir(os.Args[0])

	fh, err := os.Open(filepath.Join(appDir, cf))
	defer fh.Close()

	if err != nil {
		return this, err
	}

	jd := json.NewDecoder(fh)

	if err = jd.Decode(&this.config); err != nil {
		return this, err
	}

	if this.handle, err = sql.Open("mysql", this.config.FormatDSN()); err != nil {
		return this, err
	}

	if err = this.handle.Ping(); err != nil {
		return this, err
	}

	this.handle.QueryRow("SELECT VERSION()").Scan(&this.info)
	this.info = fmt.Sprintf("%s (%s@%s)", this.info, this.config.Addr, this.config.User)

	if this.stmt.auditInsert, err = this.handle.Prepare(auditInsertSQL); err != nil {
		return this, err
	}

	if this.stmt.checkinInsert, err = this.handle.Prepare(checkinInsertSQL); err != nil {
		return this, err
	}

	if this.stmt.serialInsert, err = this.handle.Prepare(serialInsertSQL); err != nil {
		return this, err
	}

	if this.stmt.serialUpdate, err = this.handle.Prepare(serialUpdateSQL); err != nil {
		return this, err
	}

	return this, err
}

func (this *Database) Close() {

	this.stmt.auditInsert.Close()
	this.stmt.checkinInsert.Close()
	this.stmt.serialInsert.Close()
	this.stmt.serialUpdate.Close()

	this.handle.Close()
}
