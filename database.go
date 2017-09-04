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

// Database contains the database configuration, handle, and prepared statements.
type Database struct {

	Handle *sql.DB
	Config *mysql.Config

	Info string

	Stmt struct {
		AuditInsert	*sql.Stmt
		CheckinInsert	*sql.Stmt
		SerialInsert	*sql.Stmt
		SerialUpdate	*sql.Stmt
	}
}

// NewDatbase initializes the database object, initializes the database handle,
// and prepares the prepared statements.
func NewDatabase(driver, cf string) (this *Database, err error) {

	this = new(Database)
	appDir := filepath.Dir(os.Args[0])

	fh, err := os.Open(filepath.Join(appDir, cf))
	defer fh.Close()

	if err != nil {
		return this, err
	}

	jd := json.NewDecoder(fh)

	if err = jd.Decode(&this.Config); err != nil {
		return this, err
	}

	if this.Handle, err = sql.Open(driver, this.Config.FormatDSN()); err != nil {
		return this, err
	}

	if err = this.Handle.Ping(); err != nil {
		return this, err
	}

	this.Handle.QueryRow("SELECT VERSION()").Scan(&this.Info)
	this.Info = fmt.Sprintf("Connected to %s (%s@%s)", this.Info, this.Config.User, this.Config.Addr)

	if this.Stmt.AuditInsert, err = this.Handle.Prepare(AuditInsertSQL); err != nil {
		return this, err
	}

	if this.Stmt.CheckinInsert, err = this.Handle.Prepare(CheckinInsertSQL); err != nil {
		return this, err
	}

	if this.Stmt.SerialInsert, err = this.Handle.Prepare(SerialInsertSQL); err != nil {
		return this, err
	}

	if this.Stmt.SerialUpdate, err = this.Handle.Prepare(SerialUpdateSQL); err != nil {
		return this, err
	}

	return this, err
}

// Close closes the prepared statements and database handle.
func (this *Database) Close() {

	this.Stmt.AuditInsert.Close()
	this.Stmt.CheckinInsert.Close()
	this.Stmt.SerialInsert.Close()
	this.Stmt.SerialUpdate.Close()

	this.Handle.Close()
}
