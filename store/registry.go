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
// instances. The name should be the database or schema name, not the
// database driver name.
type DataStores map[string]DataStore

// factories is a centralized registry of named references to DataStore
// factory methods. It is used by the init method of the file containing
// the DataStore implementation.
var factories = make(Factories)

// dataStores is a centralized registry of named references to initialized
// DataStore instances. It is used by initialized DataStore instances.
var dataStores = make(DataStores)

// registerFactory allows DataStore implementations to register their
// factory methods using the database driver name. It is called by the
// init method of the file containing the DataStore implementation.
func registerFactory(name string, factory func(string) (DataStore, error)) {
	factories[name] = factory
}

// registerDataStore allows initialized DataStore instances to register
// references to themselves using the database or schema name. It is called
// by the Register method of the DataStore implementation.
func registerDataStore(name string, dataStore DataStore) {
	dataStores[name] = dataStore
}

// getFactory allows callers to obtain references to DataStore factory
// methods using only the database driver name.
func GetFactory(name string) (func(string) (DataStore, error), error) {
	if factory, ok := factories[name]; !ok {
		return nil, fmt.Errorf(`factory for %q not found`, name)
	} else {
		return factory, nil
	}
}

// getDataStore allows callers to obtain references to DataStore instances
// using only the database or schema name.
func GetDataStore(name string) (DataStore, error) {
	if dataStore, ok := dataStores[name]; !ok {
		return nil, fmt.Errorf(`datastore for %q not found`, name)
	} else {
		return dataStore, nil
	}
}
