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

package usbci

import (
	`encoding/json`
	`time`
	`github.com/jmoiron/sqlx`
	`github.com/jscherff/cmdbd/store`
	`github.com/jscherff/cmdbd/utils`
)

var dataStore store.DataStore

func Init(ds store.DataStore) {
	dataStore = ds
}

type Ident struct {
	Id		int64		`db:"id,omitempty"             json:"id"`
	VendorId	string		`db:"vendor_id,omitempty"      json:"vendor_id"`
	ProductId	string		`db:"product_id,omitempty"     json:"product_id"`
	SerialNum	string		`db:"serial_number,omitempty"  json:"serial_number"`
	HostName	string		`db:"host_name,omitempty"      json:"host_name"`
	RemoteAddr	string		`db:"remote_addr,omitempty"    json:"remote_addr"`
}

type Common struct {
	VendorName	string		`db:"vendor_name,omitempty"    json:"vendor_name"`
	ProductName	string		`db:"product_name,omitempty"   json:"product_name"`
	ProductVer	string		`db:"product_ver,omitempty"    json:"product_ver"`
	FirmwareVer	string		`db:"firmware_ver,omitempty"   json:"firmware_ver"`
	SoftwareId	string		`db:"software_id,omitempty"    json:"software_id"`
	PortNumber	int		`db:"port_number,omitempty"    json:"port_number"`
	BusNumber	int		`db:"bus_number,omitempty"     json:"bus_number"`
	BusAddress	int		`db:"bus_address,omitempty"    json:"bus_address"`
	BufferSize	int		`db:"buffer_size,omitempty"    json:"buffer_size"`
	MaxPktSize	int		`db:"max_pkt_size,omitempty"   json:"max_pkt_size"`
	USBSpec		string		`db:"usb_spec,omitempty"       json:"usb_spec"`
	USBClass	string		`db:"usb_class,omitempty"      json:"usb_class"`
	USBSubClass	string		`db:"usb_subclass,omitempty"   json:"usb_subclass"`
	USBProtocol	string		`db:"usb_protocol,omitempty"   json:"usb_protocol"`
	DeviceSpeed	string		`db:"device_speed,omitempty"   json:"device_speed"`
	DeviceVer	string		`db:"device_ver,omitempty"     json:"device_ver"`
	DeviceSN	string		`db:"device_sn,omitempty"      json:"device_sn"`
	FactorySN	string		`db:"factory_sn,omitempty"     json:"factory_sn"`
	DescriptorSN	string		`db:"descriptor_sn,omitempty"  json:"descriptor_sn"`
	ObjectType	string		`db:"object_type,omitempty"    json:"object_type"`
	ObjectJSON	[]byte		`db:"object_json,omitempty"    json:"object_json"`
}

type Custom struct {
	Custom01	string		`db:"custom_01,omitempty"      json:"custom_01"`
	Custom02	string		`db:"custom_02,omitempty"      json:"custom_02"`
	Custom03	string		`db:"custom_03,omitempty"      json:"custom_03"`
	Custom04	string		`db:"custom_04,omitempty"      json:"custom_04"`
	Custom05	string		`db:"custom_05,omitempty"      json:"custom_05"`
	Custom06	string		`db:"custom_06,omitempty"      json:"custom_06"`
	Custom07	string		`db:"custom_07,omitempty"      json:"custom_07"`
	Custom08	string		`db:"custom_08,omitempty"      json:"custom_08"`
	Custom09	string		`db:"custom_09,omitempty"      json:"custom_09"`
	Custom10	string		`db:"custom_10,omitempty"      json:"custom_10"`
}

type Audit struct {
	Ident
	Changes		[]byte		`db:"changes,omitempty"        json:"changes"`
	AuditDate	time.Time	`db:"audit_date,omitempty"     json:"audit_date"`
}

type Change struct {
	Ident
	AuditId		int64		`db:"audit_id,omitempty"       json:"audit_id"`
	PropertyName	string		`db:"property_name,omitempty"  json:"property_name"`
	PreviousValue	string		`db:"previous_value,omitempty" json:"previous_value"`
	CurrentValue	string		`db:"current_value,omitempty"  json:"current_value"`
	ChangeDate	time.Time	`db:"change_date,omitempty"    json:"change_date"`
}

type Checkin struct {
	Ident
	Common
	CheckinDate	time.Time	`db:"checkin_date,omitempty"   json:"checkin_date"`
}

type Serialized struct {
	Ident
	Common
	FirstSeen	time.Time	`db:"first_seen,omitempty"     json:"first_seen"`
	LastSeen	time.Time	`db:"last_seen,omitempty"      json:"last_seen"`
	Checkins	int		`db:"checkins,omitempty"       json:"checkins"`
}

type SnRequest struct {
	Ident
	Common
	RequestDate	time.Time	`db:"request_date,omitempty"   json:"request_date"`
}

type Unserialized struct {
	Ident
	Common
	FirstSeen	time.Time	`db:"first_seen,omitempty"     json:"first_seen"`
	LastSeen	time.Time	`db:"last_seen,omitempty"      json:"last_seen"`
	Checkins	int		`db:"checkins,omitempty"       json:"checkins"`
}

type Changes []*Change

// ----------------------
// Standard CRUD methods.
// ----------------------

func (this *Audit) Create() (id int64, err error) {
	this.Id, err = dataStore.Exec(`Create`, this)
	return this.Id, err
}

func (this *Change) Create() (id int64, err error) {
	this.Id, err = dataStore.Exec(`Create`, this)
	return this.Id, err
}

func (this *Checkin) Create() (id int64, err error) {
	this.Id, err = dataStore.Exec(`Create`, this)
	return this.Id, err
}

func (this *Serialized) Create() (id int64, err error) {
	this.Id, err = dataStore.Exec(`Create`, this)
	return this.Id, err
}

func (this *SnRequest) Create() (id int64, err error) {
	this.Id, err = dataStore.Exec(`Create`, this)
	return this.Id, err
}

func (this *Unserialized) Create() (id int64, err error) {
	this.Id, err = dataStore.Exec(`Create`, this)
	return this.Id, err
}

func (this *Audit) Read() (error) {
	return dataStore.Read(`Read`, this, this)
}

func (this *Change) Read() (error) {
	return dataStore.Read(`Read`, this, this)
}

func (this *Checkin) Read() (error) {
	return dataStore.Read(`Read`, this, this)
}

func (this *Serialized) Read() (error) {
	return dataStore.Read(`Read`, this, this)
}

func (this *SnRequest) Read() (error) {
	return dataStore.Read(`Read`, this, this)
}

func (this *Unserialized) Read() (error) {
	return dataStore.Read(`Read`, this, this)
}

func (this *SnRequest) Update() (int64, error) {
	return dataStore.Exec(`Update`, this)
}

// --------------------
// Specialized methods.
// --------------------

func (this *SnRequest) UpdateSn(sn string) (int64, error) {
	this.SerialNum = sn
	return this.Update()
}

func (this *SnRequest) Unique() (bool) {

	dev := &Serialized{}
	utils.DeepCopy(this, dev)

	if err := dev.Read(); err != nil {
		return true
	} else {
		return false
	}
}

func (this *Serialized) JSON() ([]byte, error) {
	return json.Marshal(this)
}

func (this *Audit) Expand() (Changes, error) {

	var (
		changesIn [][]string
		changesOut Changes
	)

	if err := json.Unmarshal(this.Changes, &changesIn); err != nil {
		return nil, err
	}

	ident := Ident{
		VendorId:	this.VendorId,
		ProductId:	this.ProductId,
		SerialNum:	this.SerialNum,
		HostName:	this.HostName,
		RemoteAddr:	this.RemoteAddr,
	}

	for _, changeIn := range changesIn {

		changeOut := &Change{
			Ident:		ident,
			AuditId:	this.Id,
			PropertyName:	changeIn[0],
			PreviousValue:	changeIn[1],
			CurrentValue:	changeIn[2],
			ChangeDate:	this.AuditDate,
		}

		changesOut = append(changesOut, changeOut)
	}

	return changesOut, nil
}

func (this Changes) Create() (int64, error) {

	tx, err := dataStore.Begin()

	if err != nil {
		return 0, err
	}

	var (
		rows int64
		change *Change
		createChange *sqlx.NamedStmt
	)

	if stmt, err := dataStore.NamedStmt(`Create`, change); err != nil {
		return 0, err
	} else {
		createChange = tx.NamedStmt(stmt)
	}

	for _, change = range this {

		if _, err := createChange.Exec(change); err != nil {
			return 0, err
		}

		rows++
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return rows, nil
}
