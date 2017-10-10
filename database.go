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
	`strings`
	`github.com/go-sql-driver/mysql`
)

// Database contains the database configuration, handle, and prepared statements.
// It is part of the systemwide configuration under Config.Database.
type Database struct {

	*sql.DB

	Driver string
	Version string
	Config *mysql.Config
	Queries map[string][]string
	Columns map[string][]string
	Statements map[string]*sql.Stmt
}

// Init connects to the database and prepares the prepared statements.
func (this *Database) Init() (err error) {

	if this.DB, err = sql.Open(this.Driver, this.Config.FormatDSN()); err != nil {
		elog.Print(err)
		return err
	}

	if err = this.Ping(); err != nil {
		elog.Print(err)
		return err
	}

	if err = this.BuildSQL(); err != nil {
		// Error already decorated and logged.
		return err
	}

	this.QueryRow(`SELECT VERSION()`).Scan(&this.Version)

	return err
}

// Info provides identifying information about the database and user.
func (this *Database) Info() (string) {
	return fmt.Sprintf(`Database %q (%s@%s/%s) using %q driver`,
		this.Version,
		this.Config.User,
		this.Config.Addr,
		this.Config.DBName,
		this.Driver,
	)
}

// Close closes the prepared statements and database handle.
func (this *Database) Close() {
	for _, stmt := range this.Statements { stmt.Close() }
	this.DB.Close()
}

func (this *Database) BuildSQL() (err error) {

	for key, query := range this.Queries {

		rows, err := this.Query(`CALL proc_usbci_columns_for_table(?)`, query[1])

		if err != nil {
			elog.Print(err)
			return err
		}

		defer rows.Close()

		for rows.Next() {

			var col string

			if err = rows.Scan(&col); err != nil {
				elog.Print(err)
				return err
			}

			this.Columns[key] = append(this.Columns[key], col)
		}

		if err = rows.Err(); err != nil {
			elog.Print(err)
			return err
		}

		var sql string

		switch query[0] {

		case `INSERT`:

			sql = fmt.Sprintf(`INSERT INTO %s VALUES (%s?)`,
				query[1],
				strings.Repeat(`?, `, len(this.Columns[key]) - 1),
			)

		case `UPDATE`:

			sql = fmt.Sprintf(`UPDATE %s SET %s WHERE %s`,
				query[1],
				query[2],
				query[3],
			)

		case `SELECT`:

			sql = fmt.Sprintf(`SELECT %s FROM %s WHERE %s`,
				query[1],
				query[2],
				query[3],
			)
		}

		if this.Statements[key], err = this.Prepare(sql); err != nil {
			elog.Print(err)
			return err
		}
	}

	return err
}
