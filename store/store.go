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
	`fmt`
	`strings`
	`github.com/jmoiron/sqlx`
	`github.com/jscherff/cmdbd/common`
)

type DataStore interface {
	String() (string)
	Prepare(queryFile string) (error)
	Statement(queryName string, obj interface{}) (*sqlx.NamedStmt, error)
	Read(queryName string, dest, arg interface{}) (error)
	Create(queryName string, arg interface{}) (int64, error)
	Update(queryName string, arg interface{}) (int64, error)
	Delete(queryName string, arg interface{}) (int64, error)
	Begin() (*sqlx.Tx, error)
	Close() (error)
}

// dataStore is an implementation of the dataStore interface.
type dataStore struct {
	*sqlx.DB
	queries map[string]map[string]*query
	namedStmts map[string]map[string]*sqlx.NamedStmt
}

// NewDataStore returns a DataStore interface.
func NewDataStore(driver, dsn string) (DataStore, error) {
	return newDataStore(driver, dsn)
}

// newDataStore performs common tasks for creating a dataStore instance.
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
			namedStmts: make(map[string]map[string]*sqlx.NamedStmt),
		}
	}

	return this, nil
}

// String returns database version, schema, and other information.
func (this *dataStore) String() (string) {
	return this.DriverName()
}

// Prepare converts a collection of JSON-encoded Query objects into 
// a collection of Named NamedStmts.
func (this *dataStore) Prepare(queryFile string) (error) {

	if err := common.LoadConfig(&this.queries, queryFile); err != nil {
		return err
	}

	for modelName := range this.queries {

		if this.namedStmts[modelName] == nil {
			this.namedStmts[modelName] = make(map[string]*sqlx.NamedStmt)
		}

		for queryName, query := range this.queries[modelName] {

			if namedStmt, err := this.PrepareNamed(query.String()); err != nil {
				return err
			} else {
				this.namedStmts[modelName][queryName] = namedStmt
			}
		}
	}

	return nil
}

// NamedStmt looks up a NamedStmt by query name and model name and returns it.
func (this *dataStore) Statement(queryName string, obj interface{}) (*sqlx.NamedStmt, error) {

	var modelName string

	if mn, ok := obj.(string); !ok {
		modelName = strings.TrimPrefix(fmt.Sprintf(`%T`, obj), `*`)
	} else {
		modelName = mn
	}

	if stmt, ok := this.namedStmts[modelName][queryName]; !ok {
		return nil, fmt.Errorf(`statement %q for %q not found`, queryName, modelName)
	} else {
		return stmt, nil
	}
}

// Read executes a Named SELECT NamedStmt and returns the results in the 
// destination object.
func (this *dataStore) Read(queryName string, dest, arg interface{}) (error) {

	if stmt, err := this.Statement(queryName, dest); err != nil {
		return err
	} else {
		return stmt.Select(dest, arg)
	}
}

// Insert executes a Named INSERT NamedStmt and returns the last insert ID.
func (this *dataStore) Create(queryName string, arg interface{}) (int64, error) {

	if stmt, err := this.Statement(queryName, arg); err != nil {
		return 0, err
	} else if res, err := stmt.Exec(arg); err != nil {
		return 0, err
	} else {
		return res.LastInsertId()
	}
}

// Update executes a Named UPDATE NamedStmt and returns number of rows affected.
func (this *dataStore) Update(queryName string, arg interface{}) (int64, error) {

	if stmt, err := this.Statement(queryName, arg); err != nil {
		return 0, err
	} else if res, err := stmt.Exec(arg); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

// Delete executes a Named DELETE NamedStmt and returns number of rows affected.
func (this *dataStore) Delete(queryName string, arg interface{}) (int64, error) {

	if stmt, err := this.Statement(queryName, arg); err != nil {
		return 0, err
	} else if res, err := stmt.Exec(arg); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

// Begin beings a transaction and returns an *sqlx.Tx.
func (this *dataStore) Begin() (*sqlx.Tx, error) {
	return this.Beginx()
}

// Close closes all the statements.
func (this *dataStore) Close() (error) {

	for modelName := range this.namedStmts {
		for queryName := range this.namedStmts[modelName] {
			this.namedStmts[modelName][queryName].Close()
		}
	}

	return this.Close()
}
