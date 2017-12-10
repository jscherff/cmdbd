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

import `fmt`

// Factories contains named references to DataStore factory methods.
// The name should be the database driver name.
type Factories map[string]func(string) (DataStore, error)

// DataStores contains named references to initialized DataStore
// instances. The name is an arbitrary strings. It is up to the 
// implmenentation to deconflict the name space.
type DataStores map[string]DataStore

// factories is a centralized registry of named references to DataStore
// factory methods. It is used by the init method of the file containing
// the DataStore implementation.
var factories = make(Factories)

// dataStores is a centralized registry of named references to initialized
// DataStore instances.
var dataStores = make(DataStores)

// registerFactory allows DataStore implementations to register their
// factory methods using the database driver name. It is called by the
// init method of the file containing the DataStore implementation.
func registerFactory(driver string, factory func(string) (DataStore, error)) {
	factories[driver] = factory
}

// registerDataStore allows initialized DataStore instances to register
// references to themselves using arbitrary strngs such as the data store
// name (DSN), schema name, or qualified table name. It is called by the
// Register and Prepare methods of the DataStore implementation.
func registerDataStore(name string, dataStore DataStore) {
	dataStores[name] = dataStore
}

// GetFactory allows callers to obtain references to DataStore factory
// methods using the database driver name.
func GetFactory(driver string) (func(string) (DataStore, error), error) {
	if factory, ok := factories[driver]; !ok {
		return nil, fmt.Errorf(`factory for %q not found`, driver)
	} else {
		return factory, nil
	}
}

// GetDataStore allows callers to obtain references to DataStore instances
// using arbitrary strings such as data store name (DSN), schema name, or
// qualified table name.
func GetDataStore(name string) (DataStore, error) {
	if dataStore, ok := dataStores[name]; !ok {
		return nil, fmt.Errorf(`datastore for %q not found`, name)
	} else {
		return dataStore, nil
	}
}
