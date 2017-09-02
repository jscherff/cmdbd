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
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

const (
	dbConfigFile string = "dbconfig.json"

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

var (
	db *sql.DB
	dbConfig *mysql.Config
	auditInsertStmt *sql.Stmt
	checkinInsertStmt *sql.Stmt
	serialInsertStmt *sql.Stmt
	serialUpdateStmt *sql.Stmt
)

func init() {

	var err error

	// Open and decode the database JSON configuration file.

	fn := filepath.Join(filepath.Dir(os.Args[0]), dbConfigFile)
	fh, err := os.Open(fn)

	if err != nil {
		log.Fatalf("%v", err)
	}

	defer fh.Close()
	jd := json.NewDecoder(fh)

	if err = jd.Decode(&dbConfig); err != nil {
		log.Fatalf("%v", err)
	}

	// Open the database and test connectivity.

	if db, err = sql.Open("mysql", dbConfig.FormatDSN()); err != nil {
		log.Fatalf("%v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("%v", err)
	}

	var dbVer string

	db.QueryRow("SELECT VERSION()").Scan(&dbVer)
	log.Printf("Connected to %s on %s as %s.", dbVer, dbConfig.Addr, dbConfig.User)

	// Create database prepared statements.

	if auditInsertStmt, err = db.Prepare(auditInsertSQL); err != nil {
		log.Fatalf("%v", err)
	}

	if checkinInsertStmt, err = db.Prepare(checkinInsertSQL); err != nil {
		log.Fatalf("%v", err)
	}

	if serialInsertStmt, err = db.Prepare(serialInsertSQL); err != nil {
		log.Fatalf("%v", err)
	}

	if serialUpdateStmt, err = db.Prepare(serialUpdateSQL); err != nil {
		log.Fatalf("%v", err)
	}
}
