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

package store

import (
	`database/sql`
	`fmt`
	`time`
	`github.com/jmoiron/sqlx`
	`github.com/go-sql-driver/mysql`
	`github.com/jscherff/cmdbd/common`
)

// NamedStmt extends sqlx.NamedStmt.
type NamedStmt struct {
	*sqlx.NamedStmt
}

// NamedStmts is a map of NamedStmt instances.
type NamedStmts map[string]*NamedStmt

// mysqlDataStore is a MySQL database that implements the DataStore interface.
type dataStore struct {
	*sqlx.DB
}

// Register registers the DataStore in the registry using the schema name.
func (this *dataStore) Register(schemaName string) {
	registerDataStore(schemaName, this)
}

// String returns database version, schema, and other information.
func (this *dataStore) String() (string) {
	return ``
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
func (this *mysqlDataStore) exec(queryName string, arg interface{}) (sql.Result, error) {

	if stmt, ok := this.stmts[queryName]; !ok {
		return nil, fmt.Errorf(`statement %q not found`, queryName)
	} else if res, err := stmt.Exec(arg); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

// Select executes a Named SELECT Statement and returns the multi-row result
// in a slice of interfaces.
func (this *mysqlDataStore) Select(queryName string, dest, arg interface{}) (error) {

	if destSlice, ok := dest.([]interface{}); !ok {
		return fmt.Errorf(`destination must be a slice`)
	} else if stmt, ok := this.stmts[queryName]; !ok {
		return fmt.Errorf(`statement %q not found`, queryName)
	} else if err := stmt.Select(destSlice, arg); err != nil {
		return err
	}

	return nil
}

// Insert executes a Named INSERT Statement and returns the last insert ID.
func (this *mysqlDataStore) Insert(queryName string, arg interface{}) (int64, error) {

	if res, err := this.exec(queryName, arg); err != nil {
		return 0, err
	} else {
		return res.LastInsertId()
	}
}

// Insert executes a Named UPDATE or DELETE Statement and returns the number 
// of rows affected.
func (this *mysqlDataStore) Exec(queryName string, arg interface{}) (int64, error) {

	if res, err := this.exec(queryName, arg); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

// Get executes a Named SELECT Statement and returns the single-row result
// in an interface.
func (this *mysqlDataStore) Get(queryName string, dest, arg interface{}) (error) {

	if stmt, ok := this.stmts[queryName]; !ok {
		return fmt.Errorf(`statement %q not found`, queryName)
	} else if err := stmt.Get(dest, arg); err != nil {
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
