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
	`time`
	`github.com/go-sql-driver/mysql`
	`github.com/jscherff/cmdbd/common`
)

const (
	mysqlDriver = `mysql`
	mysqlLocation = `Local`
)

// mysqlDataStore is a MySQL database that implements the DataStore interface.
type mysqlDataStore struct {
	*dataStore
}

// init registers the driver name and factory method with the DataStore registry.
func init() {
	registerFactory(mysqlDriver, NewMysqlDataStore)
}

// NewmysqlDataStore creates a new instance of mysqlDataStore.
func NewMysqlDataStore(configFile string) (DataStore, error) {

	conf := &mysql.Config{}

	if err := common.LoadConfig(conf, configFile); err != nil {
		return nil, err
	}

	if loc, err := time.LoadLocation(mysqlLocation); err != nil {
		return nil, err
	} else {
		conf.Loc = loc
	}

	if this, err := NewDataStore(mysqlDriver, conf.FormatDSN()); err != nil {
		return nil, err
	} else {
		return this, nil
	}
}

// String returns database version, schema, and other information.
func (this *mysqlDataStore) String() (string) {

	sql := `SELECT VERSION() AS 'version',
		DATABASE() AS 'schema',
		USER() AS 'user'`

	var v struct {
		Version	string	`db:"version"`
		Schema	string	`db:"schema"`
		User	string	`db:"user"`
	}

	if row := this.QueryRowx(sql); row.Err() != nil {
		return mysqlDriver
	} else if err := row.StructScan(&v); err != nil {
		return mysqlDriver
	} else {
		return fmt.Sprintf(`%s version %s (%s/%s)`,
			mysqlDriver, v.Version, v.User, v.Schema,
		)
	}
}
