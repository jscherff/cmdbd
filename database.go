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
	`encoding/json`
	`database/sql`
	`fmt`
	`github.com/go-sql-driver/mysql`
)

// Database contains the database configuration and handle.
type Database struct {
	*sql.DB
	Driver string
	Config *mysql.Config
}

// NewDatabase creates and initializes a new Database instance.
func NewDatabase(cf string) (this *Database, err error) {

	this = &Database{}

	if err = loadConfig(this, cf); err != nil {
		return nil, err
	}
if b, err := json.MarshalIndent(this.Config, "", "\t"); err != nil {
	fmt.Println("Error: ", err)
} else {
	fmt.Println(this.Config.FormatDSN())
	fmt.Println(string(b))
}
	if this.DB, err = sql.Open(this.Driver, this.Config.FormatDSN()); err != nil {
		return nil, err
	}

	if err = this.Ping(); err != nil {
		return nil, err
	}

	return this, nil
}

// Info provides identifying information about the database and user.
func (this *Database) Info() (string) {

	var v string

	this.QueryRow(`SELECT VERSION()`).Scan(&v)

	return fmt.Sprintf(`Database version %s (%s@%s/%s)`, v,
		this.Config.User,
		this.Config.Addr,
		this.Config.DBName,
	)
}

// Close closes the database handle.
func (this *Database) Close() {
	this.DB.Close()
}
