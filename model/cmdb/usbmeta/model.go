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
	`github.com/jmoiron/sqlx`
	`github.com/jscherff/cmdb/meta/peripheral`
	`github.com/jscherff/cmdbd/store`
)

var dataStore store.DataStore

func Init(ds store.DataStore) {
	dataStore = ds
}

type Class struct {
	ClassId		string		`db:"class_id,omitempty" json:"class_id"`
	ClassDesc	string		`db:"class_desc,omitempty" json:"class_desc"`
	LastUpdate	time.Time	`db:"last_update,omitempty" json:"last_update"`
}

type SubClass struct {
	ClassId		string		`db:"class_id,omitempty" json:"class_id"`
	SubClassId	string		`db:"subclass_id,omitempty" json:"subclass_id"`
	SubClassDesc	string		`db:"subclass_desc,omitempty" json:"subclass_desc"`
	LastUpdate	time.Time	`db:"last_update,omitempty" json:"last_update"`
}

type Protocol struct {
	ClassId		string		`db:"class_id,omitempty" json:"class_id"`
	SubClassId	string		`db:"subclass_id,omitempty" json:"subclass_id"`
	ProtocolId	string		`db:"protocol_id,omitempty" json:"protocol_id"`
	ProtocolDesc	string		`db:"protocol_desc,omitempty" json:"protocol_desc"`
	LastUpdate	time.Time	`db:"last_update,omitempty" json:"last_update"`
}

type Vendor struct {
	VendorId	string		`db:"vendor_id,omitempty" json:"vendor_id"`
	VendorName	string		`db:"vendor_name,omitempty" json:"vendor_name"`
	LastUpdate	time.Time	`db:"last_update,omitempty" json:"last_update"`
}

type Product struct {
	VendorId	string		`db:"vendor_id,omitempty" json:"vendor_id"`
	ProductId	string		`db:"product_id,omitempty" json:"product_id"`
	ProductName	string		`db:"product_name,omitempty" json:"product_name"`
	LastUpdate	time.Time	`db:"last_update,omitempty" json:"last_update"`
}

func (this *Vendor) Create() (int64, error) {
	return dataStore.Create(`Create`, this)
}

func (this *Product) Create() (int64, error) {
	return dataStore.Create(`Create`, this)
}

func (this *Class) Create() (int64, error) {
	return dataStore.Create(`Create`, this)
}

func (this *SubClass) Create() (int64, error) {
	return dataStore.Create(`Create`, this)
}

func (this *Protocol) Create() (int64, error) {
	return dataStore.Create(`Create`, this)
}

func (this *Vendor) Read(arg interface{}) (error) {
	return dataStore.Read(`Read`, this, arg)
}

func (this *Product) Read(arg interface{}) (error) {
	return dataStore.Read(`Read`, this, arg)
}

func (this *Class) Read(arg interface{}) (error) {
	return dataStore.Read(`Read`, this, arg)
}

func (this *SubClass) Read(arg interface{}) (error) {
	return dataStore.Read(`Read`, this, arg)
}

func (this *Protocol) Read(arg interface{}) (error) {
	return dataStore.Read(`Read`, this, arg)
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

	tx, err := dataStore.Begin()

	if err != nil {
		return err
	}

	var (
		vendor *Vendor
		product *Product
		class *Class
		subClass *SubClass
		protocol *Protocol
		createVendor, createProduct, createClass, createSubClass, createProtocol *sqlx.NamedStmt
	)

	if stmt, err := dataStore.Statement(`Create`, vendor); err != nil {
		return err
	} else {
		createVendor = tx.NamedStmt(stmt)
	}

	if stmt, err := dataStore.Statement(`Create`, product); err != nil {
		return err
	} else {
		createProduct = tx.NamedStmt(stmt)
	}

	if stmt, err := dataStore.Statement(`Create`, class); err != nil {
		return err
	} else {
		createClass = tx.NamedStmt(stmt)
	}

	if stmt, err := dataStore.Statement(`Create`, subClass); err != nil {
		return err
	} else {
		createSubClass = tx.NamedStmt(stmt)
	}

	if stmt, err := dataStore.Statement(`Create`, protocol); err != nil {
		return err
	} else {
		createProtocol = tx.NamedStmt(stmt)
	}

	for vid, v := range usb.Vendors {

		vendor = &Vendor{
			VendorId: vid,
			VendorName: v.String(),
		}

		if _, err := createVendor.Exec(vendor); err != nil {
			return err
		}

		for pid, p := range v.Product {

			product = &Product{
				VendorId: vid,
				ProductId: pid,
				ProductName: p.String(),
			}

			if _, err := createProduct.Exec(product); err != nil {
				return err
			}
		}
	}

	for cid, c := range usb.Classes {

		class = &Class{
			ClassId: cid,
			ClassDesc: c.String(),
		}

		if _, err := createClass.Exec(class); err != nil {
			return err
		}

		for sid, s := range c.SubClass {

			subClass = &SubClass{
				ClassId: cid,
				SubClassId: sid,
				SubClassDesc: s.String(),
			}

			if _, err := createSubClass.Exec(subClass); err != nil {
				return err
			}

			for pid, p := range s.Protocol {

				protocol = &Protocol{
					ClassId: cid,
					SubClassId: sid,
					ProtocolId: pid,
					ProtocolDesc: p.String(),
				}

				if _, err := createProtocol.Exec(protocol); err != nil {
					return err
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
