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
	`time`
	`github.com/jmoiron/sqlx`
	`github.com/jscherff/cmdbd/utils`
)

type DataStore interface {
	String() (string)
	SetPool(configFile string) (error)
	Prepare(queryFile string) (error)
	NamedStmt(queryName string, obj interface{}) (*sqlx.NamedStmt, error)
	Exec(queryName string, arg interface{}) (int64, error)
	Read(queryName string, dest, arg interface{}) (error)
	Begin() (*sqlx.Tx, error)
	Close() (error)
}

// connPool contains database/sql settings for the connection pool.
type connPool struct {
	ConnMaxLifetime time.Duration
	MaxIdleConns int
	MaxOpenConns int
}

// namedStmt extends sqlx.NamedStmt with a query object.
type namedStmt struct {
	*sqlx.NamedStmt
	*query
}

// dataStore is an implementation of the dataStore interface.
type dataStore struct {
	*sqlx.DB
	namedStmts map[string]map[string]*namedStmt
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
			namedStmts: make(map[string]map[string]*namedStmt),
		}
	}

	return this, nil
}

// String returns database version, schema, and other information.
func (this *dataStore) String() (string) {
	return this.DriverName()
}

// SetPool configures the database connection pool.
func (this *dataStore) SetPool(configFile string) (error) {

	conf := &connPool{}

	if err := utils.LoadConfig(conf, configFile); err != nil {
		return err
	}

	this.SetConnMaxLifetime(conf.ConnMaxLifetime * time.Second)
	this.SetMaxIdleConns(conf.MaxIdleConns)
	this.SetMaxOpenConns(conf.MaxOpenConns)

	return nil
}

// Prepare converts a collection of JSON-encoded Query objects into 
// a collection of Named NamedStmts.
func (this *dataStore) Prepare(queryFile string) (error) {

	var queries map[string]map[string]*query

	if err := utils.LoadConfig(&queries, queryFile); err != nil {
		return err
	}

	for modelName := range queries {

		if this.namedStmts[modelName] == nil {
			this.namedStmts[modelName] = make(map[string]*namedStmt)
		}

		for queryName, query := range queries[modelName] {

			if stmt, err := this.PrepareNamed(query.String()); err != nil {
				return err
			} else {
				this.namedStmts[modelName][queryName] = &namedStmt{
					NamedStmt: stmt,
					query: query,
				}
			}
		}
	}

	return nil
}

// namedStmt looks up a namedStmt by query name and model name and returns it.
func (this *dataStore) namedStmt(queryName string, obj interface{}) (*namedStmt, error) {

	var modelName string

	if mnString, ok := obj.(string); !ok {
		modelName = strings.TrimPrefix(fmt.Sprintf(`%T`, obj), `*`)
	} else {
		modelName = mnString
	}

	if stmt, ok := this.namedStmts[modelName][queryName]; !ok {
		return nil, fmt.Errorf(`statement %q for %q not found`, queryName, modelName)
	} else {
		return stmt, nil
	}
}

// NamedStmt looks up a NamedStmt by query name and model name and returns it.
func (this *dataStore) NamedStmt(queryName string, obj interface{}) (*sqlx.NamedStmt, error) {

	if stmt, err := this.namedStmt(queryName, obj); err != nil {
		return nil, err
	} else {
		return stmt.NamedStmt, nil
	}
}

// Read executes a SELECT NamedStmt and returns the result in a struct for a
// single-row result or a slice of structs for a multi-row result.
func (this *dataStore) Read(queryName string, dest, arg interface{}) (error) {

	if stmt, err := this.namedStmt(queryName, dest); err != nil {
		return err
	} else if stmt.query.Command != `select` {
		return fmt.Errorf(`invalid SQL command for Read: %s`, stmt.query.Command)
	} else if stmt.query.MultiRow == true {
		return stmt.Select(dest, arg)
	} else {
		return stmt.Get(dest, arg)
	}
}

// Exec executes an INSERT, UPDATE, or DELETE NamedStmt and returns the last
// insert ID (for INSERT) or number of rows affected (for UPDATE or DELETE).
func (this *dataStore) Exec(queryName string, arg interface{}) (int64, error) {

	if stmt, err := this.namedStmt(queryName, arg); err != nil {
		return 0, err
	} else if res, err := stmt.NamedStmt.Exec(arg); err != nil {
		return 0, err
	} else if stmt.query.Command == `insert` {
		return res.LastInsertId()
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
