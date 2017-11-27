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
var Registry = map[string]func(configFile string) (DataStore, error)

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

// Columns is a slice of column names.
type Columns []string

// Tables is a map of table names to column name slices.
type Tables map[string]Columns

// Query is an interface whose methods return components of a SQL statement.
type Query interface {
	Table() string
	Command() string
	Columns() []string
	Filters() []string
	String() string
	Build(Columns) (error)
}

// Queries is an interface that represents a collection of query objects
// and whose methods load those objects from a JSON configuration file and
// return a map of SQL statement strings.
type Queries interface {
	Init(configFile string) (error)
	Get() map[string]Query
	String() map[string]string
	Build(Tables) (error)
}

// query contains SQL query components needed for building prepared statements.
type query struct {
	Table string
	Command string
	Columns []string
	Filters []string
	sqlStmt string
}

// Table returns the lowercase name of the table.
func (this *query) Table() (string) {
	return strings.ToLower(this.table)
}

// Command returns the uppercase SQL command.
func (this *query) Command() (string) {
	return strings.ToUpper(this.Commmand)
}

// Columns is a slice of target column names for INSERT, SELECT, and UPDATE
// SQL statements. Nil values indicate 'all columns.' A single empty-string
// column name indicates an empty column list for the INSERT statement (all
// defaults). A single '*' is legal only for the SELECT statement.
func (this *query) Columns() ([]string) {
	return this.Columns
}

// Filters is a slice of columns used in the conditions clause of the SQL
// statement. The interface currently only supportes ANDed conditions.
func (this *query) Filters() ([]string) {
	return this.Filters
}

// String implements the Stringer interface for the Query object and returns
// the complete SQL statement string assembled from the statement components.
func (this *query) String() (string) {
	return this.sqlStmt
}

// Build constructs the SQL query statement from the query components.
func (this *query) Build(cols Columns) (error) {

	if this.Table() == `` || this.Command() == `` {
		return fmt.Errorf(`table and command cannot be nil`)
	}

	var args, sets, flts []string

	if this.Columns() != nil {
		cols = this.Columns()
	}

	for _, col := range cols {

		if col == `*` || col == `` {
			continue
		}

		args = append(args, fmt.Sprintf(`:%s`, col))
		sets = append(sets, fmt.Sprintf(`%s = :%s`, col, col))
	}

	for _, col := range this.Filters() {
		flts = append(flts, fmt.Sprintf(`%s = :%s`, col, col))
	}

	switch this.Command() {

	case `INSERT`, `REPLACE`:
		this.sqlStmt = fmt.Sprintf(`%s INTO %s (%s) VALUES (%s)`,
			this.Command(),
			this.Table(),
			strings.Join(cols, `, `),
			strings.Join(args, `, `),
		)

	case `SELECT`:
		this.sqlStmt = fmt.Sprintf(`%s %s FROM %s`,
			this.Command(),
			strings.Join(cols, `, `),
			this.Table(),
			strings.Join(this.OrderBy(), `, `),
		)

	case `UPDATE`:
		this.sqlStmt = fmt.Sprintf(`%s %s SET %s`,
			this.Command(),
			this.Table(),
			strings.Join(sets, `, `),
		)

	case `DELETE`:
		this.sqlStmt = fmt.Sprintf(`DELETE FROM %s`,
			this.Table(),
		)

	default:
		return fmt.Errorf(`invalid command %q`, this.Command())
	}

	if flts != nil {
		sql += fmt.Sprintf(` WHERE %s`, strings.Join(flts, ` AND `))
	}

	return nil
}
