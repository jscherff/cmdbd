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

// mysqlDataStore is a MySQL database that implements the DataStore interface.
type mysqlDataStore struct {
	*sqlx.DB
	stmts map[string]*sqlx.NamedStmt
}

// NewmysqlDataStore creates a new instance of mysqlDataStore.
func NewMysqlDataStore(configFile string) (DataStore, error) {

	config := &mysql.Config{}

	if err := common.LoadConfig(config, configFile); err != nil {
		return nil, err
	}

	if location, err := time.LoadLocation(`Local`); err != nil {
		return nil, err
	} else {
		config.Loc = location
	}

	var this *mysqlDataStore

	if db, err := sqlx.Open(`mysql`, config.FormatDSN()); err != nil {
		return nil, err
	} else {
		this = &mysqlDataStore{db, make(map[string]*sqlx.NamedStmt)}
	}

	if err := this.Ping(); err != nil {
		return nil, err
	}

	return this, nil
}

// Version returns database ver, user, and schema information.
func (this *mysqlDataStore) Version() (string, error) {

	sql := `SELECT VERSION() AS 'version',
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
		return fmt.Sprintf(`version %s (%s/%s)`, v.Version, v.User, v.Schema), nil
	}
}

// Prepare converts a collection of JSON-encoded Query objects into 
// a collection of sqlx Named Statements.
func (this *mysqlDataStore) Prepare(queryFile string) (error) {

	queries := make(Queries)

	if err := common.LoadConfig(&queries, queryFile); err != nil {
		return err
	}

	for name, query := range queries {

		if stmt, err := this.PrepareNamed(query.String()); err != nil {
			return err
		} else {
			this.stmts[name] = stmt
		}
	}

	return nil
}

// exec executes a non-SELECT Named Statement and returns a sql.Result object.
// Called by Insert(), Update(), and Delete().
func (this *mysqlDataStore) exec(queryName string, args interface{}) (sql.Result, error) {

	if stmt, ok := this.stmts[queryName]; !ok {
		return nil, fmt.Errorf(`statement %q does not exist`, queryName)
	} else if res, err := stmt.Exec(args); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

// Select executes a Named SELECT Statement and returns the multi-row result
// in a slice of interfaces.
func (this *mysqlDataStore) Select(queryName string, dest interface{}, args interface{}) (error) {

	if destSlice, ok := dest.([]interface{}); !ok {
		return fmt.Errorf(`destination must be a slice`)
	} else if stmt, ok := this.stmts[queryName]; !ok {
		return fmt.Errorf(`statement %q does not exist`, queryName)
	} else if err := stmt.Select(destSlice, args); err != nil {
		return err
	}

	return nil
}

// Insert executes a Named INSERT Statement and returns the last insert ID.
func (this *mysqlDataStore) Insert(queryName string, args interface{}) (int64, error) {

	if res, err := this.exec(queryName, args); err != nil {
		return 0, err
	} else {
		return res.LastInsertId()
	}
}

// Insert executes a Named UPDATE or DELETE Statement and returns the number 
// of rows affected.
func (this *mysqlDataStore) Exec(queryName string, args interface{}) (int64, error) {

	if res, err := this.exec(queryName, args); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

// Get executes a Named SELECT Statement and returns the single-row result
// in an interface.
func (this *mysqlDataStore) Get(queryName string, dest interface{}, args interface{}) (error) {

	if stmt, ok := this.stmts[queryName]; !ok {
		return fmt.Errorf(`statement %q does not exist`, queryName)
	} else if err := stmt.Get(dest, args); err != nil {
		return err
	}

	return nil
}

// Close closes the database handle.
func (this *mysqlDataStore) Close() {

	for _, stmt := range this.stmts {
		stmt.Close()
	}

	this.Close()
}
