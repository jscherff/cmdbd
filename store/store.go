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
	`strings`
	`time`

	`github.com/jmoiron/sqlx`
	`github.com/go-sql-driver/mysql`
	`github.com/jscherff/cmdbd/common`
)

// DataStore is an interface that represents a data store.
type DataStore interface {
	Register(schemaName string)
	Prepare(queryFile string) (err error)
	Select(queryName string, dest, arg interface{}) (err error)
	Insert(queryName string, arg interface{}) (id int64, err error)
	Update(queryName string, arg interface{}) (rows int64, err error)
	Delete(queryName string, arg interface{}) (rows int64, err error)
	Get(queryName string, dest, arg interface{}) (err error)
	String() (info string)
	Close()
}

// New creates a new DataStore instance using only the the provided driver
// name and configuration file/string using the registered factory method
// associated with the driver name from the factories registry.
func New(driver, config string) (DataStore, error) {
	if factory, ok := factories[driver]; !ok {
		return nil, fmt.Errorf(`driver %q not found`, driver)
	} else {
		return factory(config)
	}
}

// dataStore is a MySQL database that implements the DataStore interface.
type dataStore struct {
	*sqlx.DB
	Stmts NamedStmts
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
func (this *dataStore) Prepare(queryFile string) (error) {

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

// Select executes a Named SELECT Statement and returns the multi-row result
// in a slice of interfaces.
func (this *dataStore) Select(queryName string, dest, arg interface{}) (error) {

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
func (this *dataStore) Insert(queryName string, arg interface{}) (int64, error) {

	if res, err := this.do(queryName, arg); err != nil {
		return 0, err
	} else {
		return res.LastInsertId()
	}
}

// Update executes a Named UPDATE Statement and returns the number of
// rows affected.
func (this *dataStore) Update(queryName string, arg interface{}) (int64, error) {

	if res, err := this.do(queryName, arg); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

// Delete executes a Named DELETE Statement and returns the number of
// rows affected.
func (this *dataStore) Delete(queryName string, arg interface{}) (int64, error) {

	if res, err := this.do(queryName, arg); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

// Get executes a Named SELECT Statement and returns the single-row result
// in an interface.
func (this *dataStore) Get(queryName string, dest, arg interface{}) (error) {

	if stmt, ok := this.stmts[queryName]; !ok {
		return fmt.Errorf(`statement %q not found`, queryName)
	} else if err := stmt.Get(dest, arg); err != nil {
		return err
	}

	return nil
}

// Close closes the database handle.
func (this *dataStore) Close() {

	for _, stmt := range this.Stmts {
		stmt.Close()
	}

	this.Close()
}

// do executes a non-SELECT Named Statement and returns a sql.Result object.
// Called by Insert(), Update(), and Delete().
func (this *dataStore) do(queryName string, arg interface{}) (sql.Result, error) {

	if stmt, ok := this.stmts[queryName]; !ok {
		return nil, fmt.Errorf(`statement %q not found`, queryName)
	} else if res, err := stmt.Exec(arg); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

// NamedStmt extends sqlx.NamedStmt.
type NamedStmt struct {
	*sqlx.NamedStmt
}

// NamedStmts is a map of NamedStmt instances.
type NamedStmts map[string]*NamedStmt

