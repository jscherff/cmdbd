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

package services

import (
	`crypto/rsa`
	`net/http`
	`time`
	`database/sql`
	`github.com/jmoiron/sqlx`
)

type dataStore struct {
	Queries struct {
		ListTables string
		ListColumns string

	ListTablesSQL string
	ListColm
var dataStores
const (
	listTables = `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = DATABASE();
	`
	listColumns = `
		SELECT column_name, column_default, extra
		FROM information_schema.columns
		WHERE table_name = ?
		AND table_schema = DATABASE();
	`
)

// Query contains SQL query components needed for building prepared statements.
type Query struct {
	Table string
	Command string
	Columns []string
	Filters []string
}

// queries is a collection of named query objects.
type Queries map[string]*Query

// Store is an interface that represents a data store.
type Store interface {
	//Config(fileName string) (error)
	Query(queryName string, args interface{}) ([]interface{}, error)
	Exec(queryName string, args interface{}) (sql.Result, error)
}

// store is an object that implements the Store interface.
type store struct {
	Db *sqlx.DB
	Stmts map[string]*sqlx.NamedStmt
	Tables map[string]string
	Columns map[string]map[string]string
}

func (this *store) Query(queryName string, args interface{}) ([]interface{}, error) {
	return nil, nil
}

func (this *store) Exec(queryName string, args interface{}) (sql.Result, error) {
	return nil, nil
}

func NewStore(driver, dsn string) (Store, error) {

	this := &store{}

	if this.db, err := sqlx.Open(driver, dsn); err != nil {
		return nil, err
	}




// storeService is a service that implements the StoreService interface.
type storeService struct {}

// NewAuthTokenService returns an object that implements the AuthTokenService interface.
func NewStoreService(db *sqlx.DB, table string, cols []string) *storeService {
	return &storeService{}
}
	if this.DB, err = sql.Open(this.Driver, this.Config.FormatDSN()); err != nil {
		return nil, err
	}

	if err = this.Ping(); err != nil {
		return nil, err
	}

	return this, nil
}

// Info provides identifying information about the database and user.
func (this *Database) Info() (string) {

	var v string

	this.QueryRow(`SELECT VERSION()`).Scan(&v)

	return fmt.Sprintf(`Database version %s (%s@%s/%s)`, v,
		this.Config.User,
		this.Config.Addr,
		this.Config.DBName,
	)
}

// Close closes the database handle.
func (this *Database) Close() {
	this.DB.Close()
}
