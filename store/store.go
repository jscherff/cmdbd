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
	`github.com/jmoiron/sqlx`
	`github.com/jscherff/cmdbd/common`
)

// DataStore provides an enhanced CRUD interface for the dataStore.
type DataStore interface {
	Register(schemaName string)
	Prepare(queryFile string) (stmts Statements, err error)
	Beginx() (trans *sqlx.Tx, err error)
	DriverName() (driver string)
	String() (version string)
	Close() (error)
}

// dataStore is an implementation of the DataStore interface.
type dataStore struct {
	*sqlx.DB
	queries map[string]map[string]*query
	statements map[string]map[string]*sqlx.NamedStmt
}

// NewDataStore creates a new instance of DataStore.
func NewDataStore(driver, dsn string) (DataStore, error) {

	if ds, err := newDataStore(driver, dsn); err != nil {
		return nil, err
	} else {
		return ds, nil
	}
}

// newDataStore performs common tasks for creating a DataStore instance.
func newDataStore(driver, dsn string) (*dataStore, error) {

	var this *dataStore

	if db, err := sqlx.Open(driver, dsn); err != nil {
		return nil, err
	} else if err := db.Ping(); err != nil {
		return nil, err
	} else {
		this = &dataStore{
			DB: db,
			queries: make(map[string]map[string]*query),
			statements: make(map[string]map[string]*sqlx.NamedStmt),
		}
	}

	return this, nil
}

// Register registers the DataStore in the registry using arbitrary names.
func (this *dataStore) Register(name string) {
	registerDataStore(name, this)
}

// String returns database version, schema, and other information.
func (this *dataStore) String() (string) {
	return this.DriverName()
}

// Prepare converts a collection of JSON-encoded Query objects into 
// a collection of Named Statements.
func (this *dataStore) Prepare(queryFile string) (Statements, error) {

	qries := make(queries)
	stmts := make(statements)

	if err := common.LoadConfig(&qries, queryFile); err != nil {
		return nil, err
	}

	for modelName := range qries {

		if stmts[modelName] == nil {
			stmts[modelName] = make(map[string]*statement)
		}

		for queryName, query := range qries[modelName] {

			if stmt, err := this.PrepareNamed(query.String()); err != nil {
				return nil, err
			} else {
				stmts[modelName][queryName] = &statement{stmt}
			}
		}
	}

	return stmts, nil
}
