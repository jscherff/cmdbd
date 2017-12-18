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
)

// Statement is a simplified CRUD interface for the sqlx.NamedStmt.
type Statement interface {
	Select(dest interface{}, arg interface{}) (error)
	Insert(arg interface{}) (int64, error)
	Update(arg interface{}) (int64, error)
	Delete(arg interface{}) (int64, error)
	Get(dest interface{}, arg interface{}) (error)
	String() (string)
	Close() error
}

// statement implements the Statement interface.
type statement struct {
	*sqlx.NamedStmt
}

// Insert simplifies the Exec method by returning just last insert ID vs sql.Result.
func (this *statement) Insert(arg interface{}) (int64, error) {

	if res, err := this.Exec(arg); err != nil {
		return 0, err
	} else {
		return res.LastInsertId()
	}
}

// Update simplifies the Exec method by returning just rows affected vs sql.Result.
func (this *statement) Update (arg interface{}) (int64, error) {

	if res, err := this.Exec(arg); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

// Delete simplifies the Exec method by returning just rows affected vs sql.Result.
func (this *statement) Delete (arg interface{}) (int64, error) {

	if res, err := this.Exec(arg); err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

// String returns the query string.
func (this *statement) String() (string) {
	return this.QueryString
}

// Statements exposes the Statement methods of the underlying collection of 
// named statements using a query name and (embedded) model name identifier.
type Statements interface {
	Statement(queryName string, obj interface{}) (Statement, error)
	Select(queryName string, dest interface{}, arg interface{}) (error)
	Insert(queryName string, arg interface{}) (int64, error)
	Update(queryName string, arg interface{}) (int64, error)
	Delete(queryName string, arg interface{}) (int64, error)
	Get(queryName string, dest interface{}, arg interface{}) (error)
	Close()
}

// statements implements the Statements interface over is a collection of
// named statements identified by query name and model name.
type statements map[string]map[string]*statement

// Statement looks up a Statement by query name and model name and returns it.
func (this statements) Statement(queryName string, obj interface{}) (Statement, error) {

	var modelName string

	if mn, ok := obj.(string); !ok {
		modelName = strings.TrimPrefix(fmt.Sprintf(`%T`, obj), `*`)
	} else {
		modelName = mn
	}

	if stmt, ok := this[modelName][queryName]; !ok {
		return nil, fmt.Errorf(`statement %q for %q not found`, queryName, modelName)
	} else {
		return stmt, nil
	}
}

// Select executes a Named SELECT Statement using filters provided by fields in
// the 'arg' struct/map and places the results in the 'dest' struct/map, which
// can also be slice of structs for multi-row results.
func (this statements) Select(queryName string, dest, arg interface{}) (error) {

	if stmt, err := this.Statement(queryName, dest); err != nil {
		return err
	} else {
		return stmt.Select(dest, arg)
	}
}

// Insert executes a Named INSERT Statement and returns the last insert ID.
func (this statements) Insert(queryName string, arg interface{}) (int64, error) {

	if stmt, err := this.Statement(queryName, arg); err != nil {
		return 0, err
	} else {
		return stmt.Insert(arg)
	}
}

// Update executes a Named UPDATE Statement and returns number of rows affected.
func (this statements) Update(queryName string, arg interface{}) (int64, error) {

	if stmt, err := this.Statement(queryName, arg); err != nil {
		return 0, err
	} else {
		return stmt.Update(arg)
	}
}

// Delete executes a Named DELETE Statement and returns number of rows affected.
func (this statements) Delete(queryName string, arg interface{}) (int64, error) {

	if stmt, err := this.Statement(queryName, arg); err != nil {
		return 0, err
	} else {
		return stmt.Delete(arg)
	}
}

// Get executes a Named SELECT Statement using filters provided by fields in
// the 'arg' struct/map and places the results in the 'dest' struct/map.
func (this statements) Get(queryName string, dest, arg interface{}) (error) {

	if stmt, err := this.Statement(queryName, dest); err != nil {
		return err
	} else {
		return stmt.Get(dest, arg)
	}
}

// Close closes all the statements.
func (this statements) Close() {

	for modelName := range this {
		for queryName := range this[modelName] {
			this[modelName][queryName].Close()
		}
	}
}
