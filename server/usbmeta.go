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

package server

import `github.com/jscherff/cmdb/meta/peripheral`

// NewUsbMeta creates and initializes a new UsbMeta instance.
func NewUsbMeta (cf string, refresh bool) (*peripheral.Usb, error) {

	this := &peripheral.Usb{}

	if usb, err := peripheral.NewUsb(cf); err != nil {
		return nil, err
	} else {
		this = usb
	}

	if refresh {
		if err := this.Refresh(); err != nil {
			return this, err
		} else if err := this.Save(cf); err != nil {
			return this, err
		}
	}

	return this, nil
}
