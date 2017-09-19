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

const (
	selectColumnsSQL = `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_name = ?
		AND table_schema = 'gocmdb'
	`
)

// Database contains the database configuration, handle, and prepared statements.
// It is part of the systemwide configuration under Config.Database.
type Database struct {

	*sql.DB

	Driver string
	Version string
	Config *mysql.Config
	Tables map[string]string
	Columns map[string][]string
	Queries map[string]string
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

	var wheres = func(fields string) (string) {

		var conds []string

		for _, field := range strings.Split(fields, `,`) {
			conds = append(conds, field + ` = ?`)
		}

		return ` WHERE ` + strings.Join(conds, ` AND `)
	}

	var sets = func(fields string) (string) {

		var sets []string

		for _, field := range strings.Split(fields, `,`) {
			sets = append(sets, field + ` = ?`)
		}

		return ` SET ` + strings.Join(sets, `, `)
	}

	for k, v := range this.Tables {

		parts := strings.Split(v, `|`)

		rows, err := this.Query(selectColumnsSQL, parts[0])

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

			this.Columns[k] = append(this.Columns[k], col)
		}

		if err = rows.Err(); err != nil {
			elog.Print(err)
			return err
		}

		clist := strings.Join(this.Columns[k], `, `)

		switch parts[1] {

		case `INSERT`:

			this.Queries[k] = fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s?)`,
				parts[0], clist, strings.Repeat(`?, `, len(this.Columns[k]) - 1))

			if len(parts) > 2 {
				this.Queries[k] += wheres(parts[2])
			}

		case `UPDATE`:

			this.Queries[k] = fmt.Sprintf(`UPDATE %s`, parts[0])
			this.Queries[k] += sets(parts[2])

			if len(parts) > 3 {
				this.Queries[k] += wheres(parts[3])
			}

		case `SELECT`:

			if parts[2] == `` {
				parts[2] = clist
			}

			this.Queries[k] = fmt.Sprintf(`SELECT %s FROM %s`,
				parts[2], parts[0])

			if len(parts) > 2 {
				this.Queries[k] += wheres(parts[3])
			}
		}

		if this.Statements[k], err = this.Prepare(this.Queries[k]); err != nil {
			elog.Print(err)
			return err
		}
	}

	return err
}
