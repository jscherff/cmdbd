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
	`fmt`
	`github.com/jscherff/cmdbd/utils`
)

// SerialSvc is an interface that creates serial numbers from seed values.
type SerialSvc interface {
	CreateSerial(key string, seed int64) (serialNum string, err error)
}

// serialSvc is a service that implements the SerialSvc interface.
type serialSvc struct {
	SerialFormat map[string]string
}

// NewSerialSvc returns an object that implements the SerialSvc interface.
func NewSerialSvc(cf string) (SerialSvc, error) {

	this := &serialSvc{}

	// Load configuration settings.

	if err := utils.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	return this, nil
}

// Format returns the format string of the provided format key.
func (this *serialSvc) serialFormat(key string) (string, error) {

	if format, ok := this.SerialFormat[key]; ok {
		return format, nil
	} else if format, ok := this.SerialFormat[`Default`]; ok {
		return format, nil
	} else {
		return ``, fmt.Errorf(`format key %q not found`, key)
	}
}

// Create generates a new serial number using the provided format key and seed.
func (this *serialSvc) CreateSerial(key string, seed int64) (string, error) {

	if format, err := this.serialFormat(key); err != nil {
		return ``, err
	} else {
		return fmt.Sprintf(format, seed), nil
	}
}
