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

package main

import `github.com/jscherff/cmdb/metaci/peripheral`

// MetaUsb contains metatata for USB devices.
type MetaUsb struct {
	*peripheral.Usb
}

// NewMetaUsb creates and initializes a new MetaUsb instance.
func NewMetaUsb (cf string) (*MetaUsb, error) {

	if usb, err := peripheral.NewUsb(cf); err != nil {
		return nil, err
	} else {
		return &MetaUsb{usb}, nil
	}
}
