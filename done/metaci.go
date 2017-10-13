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

// MetaCi contains metatata for various classes of CIs.
type MetaCi struct {
	Usb *peripheral.Usb
}

// NewMetaCi creates and initializes a new instance of MetaCI.
func NewMetaCi(uf string) (*MetaCi, error) {

	if usb, err := NewUsb(uf); err != nil {
		return nil, err
	} else {
		return &MetaCi{usb}, nil
	}
}
