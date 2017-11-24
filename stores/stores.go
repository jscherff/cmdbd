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

package stores

import (
	`database/sql`
)

// Query contains SQL query components needed for building prepared statements.
type Query struct {
	Table string
	Command string
	Columns []string
	Filters []string
}

// DataStore is an interface that represents a data store.
type DataStore interface {
	Version() (string, error)
	Tables() ([]string, error)
	Columns(table string) ([]string, error)
	Prepare(queryFile string) (error)
	Query(queryName string, args interface{}) ([]interface{}, error)
	Exec(queryName string, args interface{}) (sql.Result, error)
	Close()
}
/*
// buildQuery converts a Query object into a sqlx Named Query.
func buildQuery(query Query, store DataStore) (error) {

	// SELECT <column list or *> FROM <table> WHERE <filters>
	// INSERT INTO <table> [(<columns>)] VALUES (<values>)
	// UPDATE <table> SET <column> = <value>, <column> = <value> WHERE <filter>
	// DELETE FROM <table> WHERE <filter>

	// SELECT INSERT UPDATE DELETE

	// SELECT (columns == nil or columns == empty)
	// 	SELECT *
	// INSERT (columns == nil)
	//	INSERT INTO <table> () VALUES ()
	// INSERT (columns == [])
	//	INSERT INTO <table> (all columns) VALUES (all values)
	//

	if query.Table == nil || query.Command == nil {
		return fmt.Errorf(`query table and command must not be nil`)
	}

	var command, columns, values, filters string

	if query.Columns != nil {
		columns = strings.Join(` ,`, query.Columns)
	}

	if len(columns) == 0 {
		columns = strings.Join(` ,`, store.Columns(query.Table))
	}

	if query.Values

	if query.Filters != nil {

		for idx, col := range query.Filters {
			query.Filters[idx] = fmt.Sprintf(`%[1]s = :%[1]s`, col)
		}

		filters = strings.Join(` AND `, query.Filters)
	}



	switch Query.Command {

	case `insert`:
		command = fmt.Sprintf(`INSERT INTO %v VALUES`, Query.Table)

	case `select`:
		sql += `SELECT


{
	"InsertSequence": {
		"Table": "cmdb_sequence",
		"Command": "insert",
		"Columns": null,
		"Filters": null
	},

	"InsertError": {
		"Table": "cmdb_errors",
		"Command": "insert",
		"Columns": ["error_source", "error_code", "error_desc"],
		"Filters": null
	},

	"SelectPassword": {
		"Table": "cmdb_users",
		"Command": "select",
		"Columns": ["password"],
		"Filters": ["username"]
	},

	"InsertChanges": {
		"Table": "usbci_changes",
		"Command": "insert",
		"Columns": [],
		"Filters": null
	},

	"InsertCheckIn": {
		"Table": "usbci_checkins",
		"Command": "insert",
		"Columns": [],
		"Filters": null
	},

	"InsertSnRequest": {
		"Table": "usbci_snrequests",
		"Command": "insert",
		"Columns": [],
		"Filters": null
	},

	"UpdateSnRequest": {
		"Table": "usbci_snrequests",
		"Command": "update",
		"Columns": ["serial_number"],
		"Filters": ["id"]
	},

	"SelectSerialized": {
	       	"Table": "usbci_serialized",
		"Command": "select",
		"Columns": [],
		"Filters": ["vendor_id", "product_id", "serial_number"]
	},

	"SelectJSONObject": {
		"Table": "usbci_serialized",
		"Command": "select",
		"Columns": ["object_json"],
		"Filters": ["vendor_id", "product_id", "serial_number"]
	},

	"ReplaceVendor": {
		"Table": "usbmeta_vendor",
		"Command": "replace",
		"Columns": ["vendor_id", "vendor_name"],
		"Filters": null
	},

	"ReplaceProduct": {
		"Table": "usbmeta_product",
		"Command": "replace",
		"Columns": ["vendor_id", "product_id", "product_name"],
		"Filters": null
	},

	"ReplaceClass": {
		"Table": "usbmeta_class",
		"Command": "replace",
		"Columns": ["class_id", "class_desc"],
		"Filters": null
	},

	"ReplaceSubClass": {
		"Table": "usbmeta_subclass",
		"Command": "replace",
		"Columns": ["class_id", "subclass_id", "subclass_desc"],
		"Filters": null
	},

	"ReplaceProtocol": {
		"Table": "usbmeta_protocol",
		"Command": "replace",
		"Columns": ["class_id", "subclass_id", "protocol_id", "protocol_desc"],
		"Filters": null
	}
}
*/
