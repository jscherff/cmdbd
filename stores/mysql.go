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

package stores

import (
	`database/sql`
	`fmt`
	`time`
	`github.com/jmoiron/sqlx`
	`github.com/go-sql-driver/mysql`
	`github.com/jscherff/cmdbd/common`
)

// MySqlDataStore is a MySQL database that implements the DataStore interface.
type MySqlDataStore struct {
	*sqlx.DB
	tables []string
	columns map[string][]string
	statements map[string]*sqlx.NamedStmt
}

// NewMySqlDataStore creates a new instance of MySqlDataStore.
func NewMySqlDataStore(configFile string) (DataStore, error) {

	config := &mysql.Config{}

	if err := common.LoadConfig(config, configFile); err != nil {
		return nil, err
	}

	if location, err := time.LoadLocation(`Local`); err != nil {
		return nil, err
	} else {
		config.Loc = location
	}

	this := &MySqlDataStore{
		statements: make(map[string]*sqlx.NamedStmt),
	}

	if db, err := sqlx.Open(`mysql`, config.FormatDSN()); err != nil {
		return nil, err
	} else {
		this.DB = db
	}

	if err := this.Ping(); err != nil {
		return nil, err
	}

	return this, nil
}

// Version returns database version, user, and schema information.
func (this *MySqlDataStore) Version() (string, error) {

	sql := `SELECT VERSION(), USER(), DATABASE()`

	if params, err := this.QueryRowx(sql).SliceScan(); err != nil {
		return ``, err
	} else {
		return fmt.Sprintf(`version %v (%v/%v)`, params...), nil
	}
}

// Tables returns a slice of tables in the schema.
func (this *MySqlDataStore) Tables() ([]string, error) {

	if this.tables != nil {
		return this.tables, nil
	}

	sql := `SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = DATABASE()`

	rows, err := this.DB.Queryx(sql)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var table string

	for rows.Next() {

		if err := rows.Scan(&table); err != nil {
			return nil, err
		}

		this.tables = append(this.tables, table)
	}

	return this.tables, nil
}

// Columns returns a slice of columns in the named table.
func (this *MySqlDataStore) Columns(table string) ([]string, error) {

	if this.columns[table] != nil {
		return this.columns[table], nil
	}

	this.columns = make(map[string][]string)

	sql := `SELECT column_name, column_default, extra
		FROM information_schema.columns
		WHERE table_name = ?
		AND table_schema = DATABASE()`

	rows, err := this.Queryx(sql, table)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var v struct {
			column_name string
			column_default []byte
			extra []byte
		}

		if err := rows.StructScan(&v); err != nil {
			return nil, err
		} else if v.column_default != nil && string(v.column_default) == `CURRENT_TIMESTAMP` {
			continue
		} else if v.extra != nil && string(v.extra) == `auto_increment` {
			continue
		} else {
			this.columns[table] = append(this.columns[table], v.column_name)
		}
/*
		if err != nil {
			return nil, err
		} else if values[1] != nil && string(values[1]) == `CURRENT_TIMESTAMP` {
			continue
		} else if values[2] != nil && string(values[2]) == `auto_increment` {
			continue
		} else {
			this.columns[table] = append(this.columns[table], string(values[0]))
		}
*/
	}

	return this.columns[table], nil
}

func (this *MySqlDataStore) Query(queryName string, args interface{}) ([]interface{}, error) {
	return nil, nil
}

func (this *MySqlDataStore) Exec(queryName string, args interface{}) (sql.Result, error) {
	return nil, nil
}

// Close closes the database handle.
func (this *MySqlDataStore) Close() {

	for _, statement := range this.statements {
		statement.Close()
	}

	this.DB.Close()
}

func (this *MySqlDataStore) Prepare(queryFile string) (error) {

	queries := make(map[string]Query)

	if err := common.LoadConfig(queries, queryFile); err != nil {
		return err
	}

	for name, query := range queries {
		fmt.Println(name, query)
	}

	return nil
}
