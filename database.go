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

// SEE:
// http://www.alexedwards.net/blog/practical-persistence-sql
// https://stackoverflow.com/questions/17112852/get-the-new-record-primary-key-id-from-mysql-insert-query
// https://mariadb.com/kb/en/the-mariadb-library/create-sequence/
// https://github.com/go-sql-driver/mysql/wiki/Examples

package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var checkinSQL string = `

	INSERT INTO device_checkins (
		host_name,
		vendor_id,
		product_id,
		serial_number,
		vendor_name,
		product_name,
		product_ver,
		software_id,
		object_type)

	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

var deviceSQL string =

	`INSERT INTO devices (
		host_name,
		vendor_id,
		product_id,
		serial_number,
		vendor_name,
		product_name,
		product_ver,
		software_id,
		object_type)

	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)

	ON DUPLICATE KEY UPDATE 
		host_name = ?,
		vendor_id = ?,
		product_id = ?,
		vendor_name = ?,
		product_name = ?,
		product_ver = ?,
		software_id = ?,
		object_type = ?;`

var auditSQL string =

	`INSERT INTO device_audits (
		serial_number,
		field_name,
		old_value,
		new_value)

	VALUES (?, ?, ?, ?)`

var serialSQL string =

	`INSERT INTO issued_serials (
		host_name,
		vendor_id,
		product_id,
		serial_number,
		vendor_name,
		product_name,
		product_ver,
		software_id,
		object_type)

	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

var db *sql.DB

var checkinStmt *sql.Stmt
var deviceStmt *sql.Stmt
var auditStmt *sql.Stmt
var serialStmt *sql.Stmt

func init() {

	config, err := getConfig()

	if err != nil {
		log.Fatalf("%v", err)
	}

	db, err := sql.Open("mysql", config.FormatDSN())

	if err != nil {
		log.Fatalf("%v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("%v", err)
	}

	checkinStmt, err = db.Prepare(checkinSQL)

	if err != nil {
		log.Fatalf("%v", err)
	}

	deviceStmt, err = db.Prepare(deviceSQL)

	if err != nil {
		log.Fatalf("%v", err)
	}

	auditStmt, err = db.Prepare(auditSQL)

	if err != nil {
		log.Fatalf("%v", err)
	}

	serialStmt, err = db.Prepare(serialSQL)

	if err != nil {
		log.Fatalf("%v", err)
	}

	var version string

	db.QueryRow("SELECT VERSION()").Scan(&version)
	log.Printf("Connected to %s", version)
}
