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

	"github.com/go-sql-driver/mysql"
)

// Database contains the database configuration, handle, and prepared statements.
// It is part of the systemwide configuration under Config.Database.
type Database struct {

	*sql.DB

	Driver string
	Version string
	Config *mysql.Config

	SQL map[string]string
	Stmt map[string]*sql.Stmt
}

// Init connects to the database and prepares the prepared statements.
func (this *Database) Init() (err error) {

	if this.DB, err = sql.Open(this.Driver, this.Config.FormatDSN()); err != nil {
		return err
	}

	if err = this.Ping(); err != nil {
		return err
	}

	for k, v := range this.SQL {
		if this.Stmt[k], err = this.Prepare(v); err != nil {
			return err
		}
	}

	this.QueryRow("SELECT VERSION()").Scan(&this.Version)

	return err
}

// Info provides identifying information about the database and user.
func (this *Database) Info() (string) {
	return fmt.Sprintf("Connected to %q (%s@%s/%s) using %q driver",
		this.Version,
		this.Config.User,
		this.Config.Addr,
		this.Config.DBName,
		this.Driver,
	)
}

// Close closes the prepared statements and database handle.
func (this *Database) Close() {
	for _, stmt := range this.Stmt { stmt.Close() }
	this.DB.Close()
}
