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

// Store is an interface that represents data storage and a collection of
// methods that perform CRUD DML operations on the store.
type Store interface {
	Tables() ([]string, error)
	Columns(table string) ([]string, error)
	Select(queryName string, args interface{}) ([]interface, error)
	Update(queryName string, args interface{}) (sql.Result, error)
	Insert(queryName string, items []interface{}) (sql.Result, error)
	Delete(queryName string, args interface{}) (sql.Result, error)
	NewStmt(name, sql string) (error)
}

// store is an object that implements the Store interface.
type store struct {
	db *sqlx.DB
	tables map[string]string
	columns map[string]map[string]string
	query map[string]*sqlx.NamedStmt
}

func (this *store) Select(queryName string, args interface{}) ([]interface{}, error) {
	return new([]interface{}), nil
}

func (this *store) Update(queryName string, args interface{}) (sql.Result, error) {
	return nil
}

func (this *store) Insert(items []interface{}) (sql.Result, error) {
	return nil
}

func (this *store) Delete(args interface{}) (sql.Result, error) {
	return nil
}

func NewStore(
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
