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
	`time`
	`github.com/jscherff/cmdb/meta/peripheral`
)

// MetaUsbSvc is an interface that creates serial numbers from seed values.
type MetaUsbSvc interface {
	VendorName(vid string) (string, error)
	ProductName(vid, pid string) (string, error)
	ClassDesc(cid string) (string, error)
	SubClassDesc(cid, sid string) (string, error)
	ProtocolDesc(cid, sid, pid string) (string, error)
	LastUpdate() (time.Time)
	Raw() (*peripheral.Usb)
}

// metaUsbSvc is a service that implements the MetaUsbSvc interface.
type metaUsbSvc struct {
	*peripheral.Usb
}

// NewMetaUsbSvc returns an object that implements the MetaUsbSvc interface.
func NewMetaUsbSvc(cf string, refresh bool) (MetaUsbSvc, error) {

	var this *metaUsbSvc

	if usb, err := peripheral.NewUsb(cf); err != nil {
		return nil, err
	} else {
		this = &metaUsbSvc{usb}
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

// VendorName returns the name of the vendor given the vendor ID.
func (this *metaUsbSvc) VendorName(vid string) (string, error) {
	if vendor, err := this.GetVendor(vid); err != nil {
		return ``, err
	} else {
		return vendor.String(), nil
	}
}

// ProductName returns the name of the product given the vendor and product ID.
func (this *metaUsbSvc) ProductName(vid, pid string) (string, error) {
	if vendor, err := this.GetVendor(vid); err != nil {
		return ``, err
	} else if product, err := vendor.GetProduct(pid); err != nil {
		return ``, err
	} else {
		return product.String(), nil
	}
}

// ClassDesc returns the class description for the given class ID.
func (this *metaUsbSvc) ClassDesc(cid string) (string, error) {
	if class, err := this.GetClass(cid); err != nil {
		return ``, err
	} else {
		return class.String(), nil
	}
}

// SubClassDesc returns the subclass description for the given class and subclass ID.
func (this *metaUsbSvc) SubClassDesc(cid, sid string) (string, error) {
	if class, err := this.GetClass(cid); err != nil {
		return ``, err
	} else if subClass, err := class.GetSubClass(sid); err != nil {
		return ``, err
	} else {
		return subClass.String(), nil
	}
}

// ProtocolDesc returns the protocol description for the given class, subclass, and protocol ID.
func (this *metaUsbSvc) ProtocolDesc(cid, sid, pid string) (string, error) {
	if class, err := this.GetClass(cid); err != nil {
		return ``, err
	} else if subClass, err := class.GetSubClass(sid); err != nil {
		return ``, err
	} else if protocol, err := subClass.GetProtocol(pid); err != nil {
		return ``, err
	} else {
		return protocol.String(), nil
	}
}

// LastUpdate returns the date the metadata was last updated from source.
func (this *metaUsbSvc) LastUpdate() (time.Time) {
	return this.LastUpdate()
}

// Raw returns the raw underlying metadata object.
func (this *metaUsbSvc) Raw() (*peripheral.Usb) {
	return this.Usb
}
