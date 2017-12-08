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
)

// Print formats for generating named parameters.
const (
	namedParamFmt = `:%v`
	namedEqualFmt = `%[1]v = :%[1]v`
)

// query contains SQL Xquery components needed for building prepared statements.
type query struct {
	Table string
	Command string
	Filters []string
	Columns []string
	qString string
}

// String implements the Stringer interface for the Query object and returns
// the complete SQL statement string assembled from the statement components.
func (this *query) String() (string) {

	// Return stored query string if it has already been generated.

	if this.qString != `` {
		return this.qString
	}

	// Verify we have at least the minimum query elements for a valid query.

	if this.Table == `` || this.Command == `` {
		return ``
	}

	// Convenience method for generating named parameters in query strings.

	toFormat := func(format string, cols []string) (list []string) {
		for _, col := range cols {
			if col == `*` { continue }
			col = strings.ToLower(col)
			list = append(list, fmt.Sprintf(format, col))
		}
		return list
	}

	// Convert the table name to lowercase.

	table := strings.ToLower(this.Table)

	// Convert the SQL command to uppercase.

	command := strings.ToUpper(this.Command)

	// Generate comma-separated list of lowercase column names.

	columns := strings.ToLower(strings.Join(this.Columns, `, `))

	// Generate comma-separated list of named parameters.

	params := strings.Join(toFormat(namedParamFmt, this.Columns), `, `)

	// Generate comma-separated list of column assignments.

	setters := strings.Join(toFormat(namedEqualFmt, this.Columns), `, `)

	// Generate AND-separated list of condition filters.

	filters := strings.Join(toFormat(namedEqualFmt, this.Filters), ` AND `)

	switch command {

	case `INSERT`, `REPLACE`:
		this.qString = fmt.Sprintf(`%s INTO %s (%s) VALUES (%s)`,
			command, table, columns, params,
		)

	case `SELECT`:
		this.qString = fmt.Sprintf(`%s %s FROM %s`,
			command, columns, table,
		)

	case `UPDATE`:
		this.qString = fmt.Sprintf(`%s %s SET %s`,
			command, table, setters,
		)

	case `DELETE`:
		this.qString = fmt.Sprintf(`DELETE FROM %s`,
			table,
		)

	default:
		return ``
	}

	if filters != `` {
		this.qString = fmt.Sprintf(`%s WHERE %s`,
			this.qString, filters,
		)
	}

	return this.qString
}

type queries struct {
	Driver string
	Schema string
	Query map[string]*query
}

func NewQueries(cf string)
