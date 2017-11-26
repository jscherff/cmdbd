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
	`fmt`
	`strings`
)

// Registry is a function map of data store types to factory methods.
var Registry = map[string]func(string) (DataStore, error)

// Register registers the factory function of the named datastore.
func Register(name string, factory func(string) (DataStore, error)) {
	Registry[name] = factory
}

// Factory returns the data store factory method of the named datastore.
func Factory(name string) (func(string) (DataStore, error), error) {
	if factory, ok := Registry[name]; !ok {
		return nil, fmt.Errorf(`data store %q does not exist`, name)
	} else {
		return factory, nil
	}
}

// DataStore is an interface that represents a data store.
type DataStore interface {
	Version() (string, error)
	Tables() ([]string, error)
	Columns(table string) ([]string, error)
	Prepare(queryFile string) (error)
	Query(queryName string, dest []interface{}, args []interface{}) (error)
	Get(queryName string, dest interface{}, args interface{}) (error)
	Exec(queryName string, args interface{}) (sql.Result, error)
	Close()
}


type Query interface {
	Table() string
	Command() string
	Columns() string
	Filters() string
	//Order() string
	SQL() string
}

type Queries interface {
	Load(queryFile string) (error)
}

// Query contains SQL query components needed for building prepared statements.
type query struct {
	Table string
	Command string
	Columns []string
	Filters []string
}

type queries []query

func (this *query) Table() (string) {
	return strings.ToLower(this.table)
}

func (this *query) Command() (string) {
	return strings.ToUpper(this.Commmand)
}

func (this *query) Columns() ([]string) {
	return this.Columns
}

func (this *query) Filters() ([]string) {
	return this.Filters
}

// SQL converts a Query object into a SQL string.
func (this *Query) SQL(allColumns []string) (string, error) {

	if this.Table == `` || this.Command == `` {
		return ``, fmt.Errorf(`table and command must not be nil`)
	}

	table, command := this.Table, strings.ToUpper(this.Command)

	var (
		sql string
		columns, params, setters, filters []string
	)

	if this.Columns != nil {
		columns = this.Columns
	} else {
		columns = allColumns
	}

	for _, col := range columns {

		if col == `*` || col == `` {
			continue
		}

		params = append(params, fmt.Sprintf(`:%s`, col))
		setters = append(setters, fmt.Sprintf(`%s = :%s`, col, col))
	}

	for _, col := range this.Filters {
		filters = append(filters, fmt.Sprintf(`%s = :%s`, col, col))
	}

	switch command {

	case `INSERT`, `REPLACE`:
		sql = fmt.Sprintf(`%s INTO %s (%s) VALUES (%s)`,
			command,
			table,
			strings.Join(columns, `, `),
			strings.Join(params, `, `),
		)

	case `SELECT`:
		sql = fmt.Sprintf(`%s %s FROM %s`,
			command,
			strings.Join(columns, `, `),
			table,
		)

	case `UPDATE`:
		sql = fmt.Sprintf(`%s %s SET %s`,
			command,
			table,
			strings.Join(setters, `, `),
		)

	case `DELETE`:
		sql = fmt.Sprintf(`DELETE FROM %s`,
			table,
		)

	default:
		return ``, fmt.Errorf(`invalid command %q`, this.Command)
	}

	if filters != nil {
		sql += fmt.Sprintf(` WHERE %s`, strings.Join(filters, ` AND `))
	}

	return sql, nil
}
