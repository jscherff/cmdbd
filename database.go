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

var auditInsertSQL string = `

	INSERT INTO audits (
		serial_number,
		field_name,
		old_value,
		new_value
	)

	VALUES (?, ?, ?, ?)`

var checkinInsertSQL string = `

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

	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

var serialInsertSQL string = `

	INSERT INTO serials (
		host_name,
		vendor_id,
		product_id,
		vendor_name,
		product_name,
		product_ver,
		software_id,
		object_type
	)

	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

var serialUpdateSQL string = `

	UPDATE serials
	SET serial_number = ?
	WHERE id = ?`

var db *sql.DB

var auditInsertStmt *sql.Stmt
var checkinInsertStmt *sql.Stmt
var serialInsertStmt *sql.Stmt
var serialUpdateStmt *sql.Stmt

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

	auditInsertStmt, err = db.Prepare(auditInsertSQL)

	if err != nil {
		log.Fatalf("%v", err)
	}

	checkinInsertStmt, err = db.Prepare(checkinInsertSQL)

	if err != nil {
		log.Fatalf("%v", err)
	}

	serialInsertStmt, err = db.Prepare(serialInsertSQL)

	if err != nil {
		log.Fatalf("%v", err)
	}

	serialUpdateStmt, err = db.Prepare(serialUpdateSQL)

	if err != nil {
		log.Fatalf("%v", err)
	}

	var version string

	db.QueryRow("SELECT VERSION()").Scan(&version)
	log.Printf("Connected to %s", version)
}
