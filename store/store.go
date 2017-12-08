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

	`github.com/jmoiron/sqlx`
	`github.com/jscherff/cmdbd/common`
)

// driverName is the database driver name for the dataStore.
var driverName = `undefined`

// DataStore provides an enhanced CRUD interface for the dataStore.
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

// New creates a new DataStore instance using the registered factory
// method associated with the provided driver name.
func New(driver, dsn string) (DataStore, error) {
	if factory, ok := factories[driver]; !ok {
		return nil, fmt.Errorf(`driver %q not found`, driver)
	} else {
		return factory(dsn)
	}
}

// NewDataStore creates a new instance of DataStore.
func NewDataStore(driver, dsn string) (DataStore, error) {

	var this *dataStore
	driverName = driver

	if db, err := sqlx.Open(driver, dsn); err != nil {
		return nil, err
	} else if err := db.Ping(); err != nil {
		return nil, err
	} else {
		stmts := make(map[string]map[string]*sqlx.NamedStmt)
		this = &dataStore{db, stmts}
	}

	this.Register(dsn)

	return this, nil
}

// dataStore is an implementation of the DataStore interface.
type dataStore struct {
	*sqlx.DB
	stmts map[string]map[string]*sqlx.NamedStmt
}

// Register registers the DataStore in the registry using the schema name.
func (this *dataStore) Register(schemaName string) {
	registerDataStore(schemaName, this)
}

// String returns database version, schema, and other information.
func (this *dataStore) String() (string) {
	return driverName
}

// Prepare converts a collection of JSON-encoded Query objects into 
// a collection of sqlx Named Statements.
func (this *dataStore) Prepare(queryFile string) (error) {

	var queries map[string]map[string]*query

	if err := common.LoadConfig(&queries, queryFile); err != nil {
		return err
	}

	for modelName := range queries {

		if this.stmts[modelName] == nil {
			this.stmts[modelName] = make(map[string]*sqlx.NamedStmt)
		}

		for queryName, query := range queries[modelName] {

			if stmt, err := this.PrepareNamed(query.String()); err != nil {
				return err
			} else {
				this.stmts[modelName][queryName] = stmt
			}
		}
	}

	return nil
}

// Select executes a Named SELECT Statement and returns the multi-row result
// in a slice of interfaces.
func (this *dataStore) Select(queryName string, dest, arg interface{}) (error) {

	modelName := ModelName(dest)

	if destSlice, ok := dest.([]interface{}); !ok {
		return fmt.Errorf(`destination must be a slice`)
	} else if stmt, ok := this.stmts[modelName][queryName]; !ok {
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

	modelName := ModelName(dest)

	if stmt, ok := this.stmts[modelName][queryName]; !ok {
		return fmt.Errorf(`statement %q not found`, queryName)
	} else if err := stmt.Get(dest, arg); err != nil {
		return err
	}

	return nil
}

// do executes a non-SELECT Named Statement and returns a sql.Result object.
// Called by Insert(), Update(), and Delete().
func (this *dataStore) do(queryName string, arg interface{}) (sql.Result, error) {

	modelName := ModelName(arg)

	if stmt, ok := this.stmts[modelName][queryName]; !ok {
		return nil, fmt.Errorf(`statement %q not found`, queryName)
	} else if res, err := stmt.Exec(arg); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

// Close closes the database handle.
func (this *dataStore) Close() {

	for modelName := range this.stmts {
		for queryName := range this.stmts[modelName] {
			this.stmts[modelName][queryName].Close()
		}
	}

	this.Close()
}

// modelName derives the model name from the object's type.
func ModelName(t interface{}) (string) {
	return strings.TrimPrefix(fmt.Sprintf(`%T`, t), `*`)
}
