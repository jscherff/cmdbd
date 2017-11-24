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

	sql := `SELECT table_name, table_type
		FROM information_schema.tables
		WHERE table_schema = DATABASE()`

	rows, err := this.DB.Queryx(sql)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var v struct {
		TabName	string	`db:"table_name"`
		TabType	string	`db:"table_type"`
	}

	for rows.Next() {

		if err := rows.StructScan(&v); err != nil {
			return nil, err
		} else if v.TabType != `BASE TABLE` {
			continue
		}

		this.tables = append(this.tables, v.TabName)
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

	var v struct {
		ColName	string	`db:"column_name"`
		ColDflt	[]byte	`db:"column_default"`
		Extra	[]byte	`db:"extra"`
	}

	for rows.Next() {

		if err := rows.StructScan(&v); err != nil {
			return nil, err
		} else if v.ColDflt != nil && string(v.ColDflt) == `CURRENT_TIMESTAMP` {
			continue
		} else if v.Extra != nil && string(v.Extra) == `auto_increment` {
			continue
		}

		this.columns[table] = append(this.columns[table], v.ColName)
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
