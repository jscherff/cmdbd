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

func init() {
	Register(`mysql`, NewMySqlDataStore)
}

// MySqlDataStore is a MySQL database that implements the DataStore interface.
type MySqlDataStore struct {
	*sqlx.DB
	ver	string
	tbls	[]string
	cols	map[string][]string
	query	map[string]*sqlx.NamedStmt
	exec	map[string]*sqlx.NamedStmt
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
		query:	make(map[string]*sqlx.NamedStmt),
		exec:	make(map[string]*sqlx.NamedStmt),
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

// Version returns database ver, user, and schema information.
func (this *MySqlDataStore) Version() (string, error) {

	if this.ver != `` {
		return this.ver, nil
	}

	sql := `SELECT VERSION() AS 'ver',
		DATABASE() AS 'schema',
		USER() AS 'user'`

	var v struct {
		Version	string	`db:"version"`
		Schema	string	`db:"schema"`
		User	string	`db:"user"`
	}

	if row := this.QueryRowx(sql); row.Err() != nil {
		return ``, row.Err()
	} else if err := row.StructScan(&v); err != nil {
		return ``, err
	} else {
		this.ver = fmt.Sprintf(`ver %s (%s/%s)`, v.Version, v.User, v.Schema)
	}

	return this.ver, nil
}

// Tables returns a slice of tbls in the schema.
func (this *MySqlDataStore) Tables() ([]string, error) {

	if this.tbls != nil {
		return this.tbls, nil
	}

	sql := `SELECT table_name, table_type
		FROM information_schema.tbls
		WHERE table_schema = DATABASE()`

	var v struct {
		TabName	string	`db:"table_name"`
		TabType	string	`db:"table_type"`
	}

	rows, err := this.Queryx(sql)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		if err := rows.StructScan(&v); err != nil {
			return nil, err
		} else if v.TabType != `BASE TABLE` {
			continue
		}

		this.tbls = append(this.tbls, v.TabName)
	}

	return this.tbls, nil
}

// Columns returns a slice of cols in the named table.
func (this *MySqlDataStore) Columns(table string) ([]string, error) {

	if this.cols[table] != nil {
		return this.cols[table], nil
	}

	this.cols = make(map[string][]string)

	sql := `SELECT column_name, column_default, extra
		FROM information_schema.cols
		WHERE table_name = ?
		AND table_schema = DATABASE()`

	var v struct {
		ColName	string	`db:"column_name"`
		ColDflt	[]byte	`db:"column_default"`
		Extra	[]byte	`db:"extra"`
	}

	rows, err := this.Queryx(sql, table)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		if err := rows.StructScan(&v); err != nil {
			return nil, err
		} else if v.ColDflt != nil && string(v.ColDflt) == `CURRENT_TIMESTAMP` {
			continue
		} else if v.Extra != nil && string(v.Extra) == `auto_increment` {
			continue
		}

		this.cols[table] = append(this.cols[table], v.ColName)
	}

	return this.cols[table], nil
}

// Prepare converts a collection of JSON-encoded Query objects into 
// a collection of sqlx Named Statements.
func (this *MySqlDataStore) Prepare(queryFile string) (error) {

	var queries = make(map[string]*Query)

	if err := common.LoadConfig(&queries, queryFile); err != nil {
		return err
	}

	for name, query := range queries {

		cols, err := this.Columns(query.Table)

		if err != nil {
			return err
		}

		if sql, err := query.SQL(cols); err != nil {
			return err
		} else if stmt, err := this.PrepareNamed(sql); err != nil {
			return err
		} else if strings.ToUpper
			this.query[name] = stmt
		}
	}

	return nil
}

// Query executes a select Named Statement and returns the multiple-row result
// in a slice of interfaces.
func (this *MySqlDataStore) Query(queryName string, dest []interface{}, args interface{}) error) {

	if stmt, ok := this.query[queryName]; !ok {
		return fmt.Errorf(`statement %q does not exist`, queryName)
	} else if err := stmt.Select(destSlice, args); err != nil {
		return err
	}

	return nil
}

// Query executes a non-select Named Statement and returns the last insert ID
// and number of rows affected in a sql.Result object.
func (this *MySqlDataStore) Exec(queryName string, args interface{}) (sql.Result, error) {

	if stmt, ok := this.exec[queryName]; !ok {
		return nil, fmt.Errorf(`statement %q does not exist`, queryName)
	} else if res, err := stmt.Get(dest, args); err != nil {
		return nil, err
	}

	return res, nil

}

// Get executes a select Named Statement and returns the single-row
// result in an interface.
func (this *MySqlDataStore) Get(queryName string, dest interface{}, args interface{}) (error) {

	if stmt, ok := this.query[queryName]; !ok {
		return fmt.Errorf(`statement %q does not exist`, queryName)
	} else if err := stmt.Get(dest, args); err != nil {
		return err
	}

	return nil
}

// Close closes the database handle.
func (this *MySqlDataStore) Close() {

	for _, stmt := range this.query {
		stmt.Close()
	}

	for _, stmt := range this.exec {
		exec.Close()
	}

	this.Close()
}
