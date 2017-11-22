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

package service

import (
	`crypto/rsa`
	`net/http`
	`time`
	`database/sql`
	`github.com/jmoiron/sqlx`
)

const (
	listTables = `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = DATABASE();
	`
	listColumns = `
		SELECT column_name
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
	db *sqlx.DB
	Stmts map[string]*sqlx.NamedStmt
	Tables map[string]string
	Columns map[string]map[string]string
	Queries map[string]*Query
}

func (this *store) Query(queryName string, args interface{}) ([]interface{}, error) {
	return nil, nil
}

func (this *store) Exec(queryName string, args interface{}) (sql.Result, error) {
	return nil, nil
}

func NewStore(driver, dsn
// StoreService is an interface that creates new Stores.
type StoreService interface {
	Create(config map[string]string) (Store, error)
}

// storeService is a service that implements the StoreService interface.
type storeService struct {}

// NewAuthTokenService returns an object that implements the AuthTokenService interface.
func NewStoreService(db *sqlx.DB, table string, cols []string) *storeService {
	return &storeService{}
}




{
	"Query": {
		"cmdbInsertSequence": [
			"INSERT_EMPTY",
			"cmdb_sequence"
		],
		"cmdbSelectUserPassword": [
			"SELECT_LIST",
			"password",
			"cmdb_users",
			"username = ?"
		],
		"usbCiInsertChanges": [
			"INSERT_ALL",
			"usbci_changes"
		],
		"usbCiInsertCheckin": [
		       	"INSERT_ALL",
			"usbci_checkins"
		],
		"usbCiInsertSnRequest": [
			"INSERT_ALL",
			"usbci_snrequests"
		],
		"usbCiUpdateSnRequest": [
       			"UPDATE_LIST",
			"usbci_snrequests",
			"serial_number = ?",
			"id = ?"
		],
		"usbCiSelectSerialized": [
			"SELECT_ALL",
		       	"usbci_serialized",
			"vendor_id = ? AND product_id = ? AND serial_number = ?"
		],
		"usbCiSelectJSONObject": [
			"SELECT_LIST",
			"object_json",
			"usbci_serialized",
			"vendor_id = ? AND product_id = ? AND serial_number = ?"
		],
		"usbMetaReplaceVendor": [
			"REPLACE_LIST",
			"usbmeta_vendor",
			"vendor_id, vendor_name",
			"?, ?"
		],
		"usbMetaReplaceProduct": [
			"REPLACE_LIST",
			"usbmeta_product",
			"vendor_id, product_id, product_name",
			"?, ?, ?"
		],
		"usbMetaReplaceClass": [
			"REPLACE_LIST",
			"usbmeta_class",
			"class_id, class_desc",
			"?, ?"
		],
		"usbMetaReplaceSubClass": [
			"REPLACE_LIST",
			"usbmeta_subclass",
			"class_id, subclass_id, subclass_desc",
			"?, ?, ?"
		],
		"usbMetaReplaceProtocol": [
			"REPLACE_LIST",
			"usbmeta_protocol",
			"class_id, subclass_id, protocol_id, protocol_desc",
			"?, ?, ?, ?"
		]
	}
}
