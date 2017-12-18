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

package usbmeta

import (
	`time`
	`github.com/jscherff/cmdb/meta/peripheral`
	`github.com/jscherff/cmdbd/store`
)

var Stmts store.Statements

func Init(stmts store.Statements) {
	Stmts = stmts
}

type Class struct {
	ClassID		string		`db:"class_id,omitempty"`
	ClassDesc	string		`db:"class_desc,omitempty"`
	LastUpdate	time.Time	`db:"last_update,omitempty"`
}

type SubClass struct {
	ClassID		string		`db:"class_id,omitempty"`
	SubClassID	string		`db:"subclass_id,omitempty"`
	SubClassDesc	string		`db:"subclass_desc,omitempty"`
	LastUpdate	time.Time	`db:"last_update,omitempty"`
}

type Protocol struct {
	ClassID		string		`db:"class_id,omitempty"`
	SubClassID	string		`db:"subclass_id,omitempty"`
	ProtocolID	string		`db:"protocol_id,omitempty"`
	ProtocolDesc	string		`db:"protocol_desc,omitempty"`
	LastUpdate	time.Time	`db:"last_update,omitempty"`
}

type Vendor struct {
	VendorID	string		`db:"vendor_id,omitempty"`
	VendorName	string		`db:"vendor_name,omitempty"`
	LastUpdate	time.Time	`db:"last_update,omitempty"`
}

type Product struct {
	VendorID	string		`db:"vendor_id,omitempty"`
	ProductID	string		`db:"product_id,omitempty"`
	ProductName	string		`db:"product_name,omitempty"`
	LastUpdate	time.Time	`db:"last_update,omitempty"`
}

func (this *Vendor) Create() (int64, error) {
	return Stmts.Insert(`Create`, this)
}

func (this *Product) Create() (int64, error) {
	return Stmts.Insert(`Create`, this)
}

func (this *Class) Create() (int64, error) {
	return Stmts.Insert(`Create`, this)
}

func (this *SubClass) Create() (int64, error) {
	return Stmts.Insert(`Create`, this)
}

func (this *Protocol) Create() (int64, error) {
	return Stmts.Insert(`Create`, this)
}

func (this *Vendor) Read(arg interface{}) (error) {
	return Stmts.Get(`Read`, this, arg)
}

func (this *Product) Read(arg interface{}) (error) {
	return Stmts.Get(`Read`, this, arg)
}

func (this *Class) Read(arg interface{}) (error) {
	return Stmts.Get(`Read`, this, arg)
}

func (this *SubClass) Read(arg interface{}) (error) {
	return Stmts.Get(`Read`, this, arg)
}

func (this *Protocol) Read(arg interface{}) (error) {
	return Stmts.Get(`Read`, this, arg)
}

func (this *Vendor) String() (string) {
	return this.VendorName
}

func (this *Product) String() (string) {
	return this.ProductName
}

func (this *Class) String() (string) {
	return this.ClassDesc
}

func (this *SubClass) String() (string) {
	return this.SubClassDesc
}

func (this *Protocol) String() (string) {
	return this.ProtocolDesc
}

// Load updates the USB metadata tables in the database.
func Load(usb *peripheral.Usb) (error) {

	for vid, v := range usb.Vendors {

		vendor := &Vendor{
			VendorID: vid,
			VendorName: v.String(),
		}

		if _, err := vendor.Create(); err != nil {
			return err
		}

		for pid, p := range v.Product {

			product := &Product{
				VendorID: vid,
				ProductID: pid,
				ProductName: p.String(),
			}

			if _, err := product.Create(); err != nil {
				return err
			}
		}
	}

	for cid, c := range usb.Classes {

		class := &Class{
			ClassID: cid,
			ClassDesc: c.String(),
		}

		if _, err := class.Create(); err != nil {
			return err
		}

		for sid, s := range c.SubClass {

			subClass := &SubClass{
				ClassID: cid,
				SubClassID: sid,
				SubClassDesc: s.String(),
			}

			if _, err := subClass.Create(); err != nil {
				return err
			}

			for pid, p := range s.Protocol {

				protocol := &Protocol{
					ClassID: cid,
					SubClassID: sid,
					ProtocolID: pid,
					ProtocolDesc: p.String(),
				}

				if _, err := protocol.Create(); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
