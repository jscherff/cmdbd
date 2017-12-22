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

	if ds, err := newDataStore(mysqlDriver, conf.FormatDSN()); err != nil {
		return nil, err
	} else {
		return &mysqlDataStore{ds}, nil
	}
}

func (this *mysqlDataStore) String() (string) {

	info := this.DriverName()

	var ver struct {
		Version	string	`db:"version"`
		Schema	string	`db:"schema"`
		User	string	`db:"user"`
	}

	sql := `SELECT VERSION() AS 'version',
		DATABASE() AS 'schema',
		USER() AS 'user'`

	if err := this.Get(&ver, sql); err != nil {
		return info
	} else {
		return fmt.Sprintf(`%s %s %s/%s`,
			info, ver.Version, ver.User, ver.Schema,
		)
	}
}
