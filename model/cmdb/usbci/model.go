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
	`github.com/jscherff/cmdbd/store`
)

var Stmts store.Statements

func Init(stmts store.Statements) {
	Stmts = stmts
}

type Ident struct {
	Id		interface{}	`db:"id,omitempty"`
	VendorID	string		`db:"vendor_id,omitempty"`
	ProductID	string		`db:"product_id,omitempty"`
	SerialNum	string		`db:"serial_number,omitempty"`
	HostName	string		`db:"host_name,omitempty"`
	RemoteAddr	string		`db:"remote_addr,omitempty"`
}

type Common struct {
	VendorName	string		`db:"vendor_name,omitempty"`
	ProductName	string		`db:"product_name,omitempty"`
	ProductVer	string		`db:"product_ver,omitempty"`
	FirmwareVer	string		`db:"firmware_ver,omitempty"`
	SoftwareID	string		`db:"software_id,omitempty"`
	PortNumber	int		`db:"port_number,omitempty"`
	BusNumber	int		`db:"bus_number,omitempty"`
	BusAddress	int		`db:"bus_address,omitempty"`
	BufferSize	int		`db:"buffer_size,omitempty"`
	MaxPktSize	int		`db:"max_pkt_size,omitempty"`
	USBSpec		string		`db:"usb_spec,omitempty"`
	USBClass	string		`db:"usb_class,omitempty"`
	USBSubClass	string		`db:"usb_subclass,omitempty"`
	USBProtocol	string		`db:"usb_protocol,omitempty"`
	DeviceSpeed	string		`db:"device_speed,omitempty"`
	DeviceVer	string		`db:"device_ver,omitempty"`
	DeviceSN	string		`db:"device_sn,omitempty"`
	FactorySN	string		`db:"factory_sn,omitempty"`
	DescriptorSN	string		`db:"descriptor_sn,omitempty"`
	ObjectType	string		`db:"object_type,omitempty"`
	ObjectJSON	[]byte		`db:"object_json,omitempty"`
}

type Custom struct {
	Custom01	string		`db:"custom_01,omitempty"`
	Custom02	string		`db:"custom_02,omitempty"`
	Custom03	string		`db:"custom_03,omitempty"`
	Custom04	string		`db:"custom_04,omitempty"`
	Custom05	string		`db:"custom_05,omitempty"`
	Custom06	string		`db:"custom_06,omitempty"`
	Custom07	string		`db:"custom_07,omitempty"`
	Custom08	string		`db:"custom_08,omitempty"`
	Custom09	string		`db:"custom_09,omitempty"`
	Custom10	string		`db:"custom_10,omitempty"`
}

type Checkin struct {
	Ident
	Common
	CheckinDate	time.Time	`db:"checkin_date,omitempty"`
}

type SnRequest struct {
	Ident
	Common
	RequestDate	time.Time	`db:"request_date,omitempty"`
}

type Serialized struct {
	Ident
	Common
	FirstSeen	time.Time	`db:"first_seen,omitempty"`
	LastSeen	time.Time	`db:"last_seen,omitempty"`
	Checkins	int		`db:"checkins,omitempty"`
}

type Unserialized struct {
	Ident
	Common
	FirstSeen	time.Time	`db:"first_seen,omitempty"`
	LastSeen	time.Time	`db:"last_seen,omitempty"`
	Checkins	int		`db:"checkins,omitempty"`
}

type Audit struct {
	Ident
	Changes		[]byte		`db:"changes,omitempty"`
	AuditDate	time.Time	`db:"audit_date,omitempty"`
}

type Change struct {
	Ident
	PropertyName	string		`db:"property_name,omitempty"`
	PreviousValue	string		`db:"previous_value,omitempty"`
	CurrentValue	string		`db:"current_value,omitempty"`
	AuditDate	time.Time	`db:"audit_date,omitempty"`
}

type Changes []Change

func (this *Audit) Create() (int64, error) {
	return Stmts.Insert(`Create`, this)
}

func (this *Change) Create() (int64, error) {
	return Stmts.Insert(`Create`, this)
}

func (this *Checkin) Create() (int64, error) {
	return Stmts.Insert(`Create`, this)
}

func (this *SnRequest) Create() (int64, error) {
	return Stmts.Insert(`Create`, this)
}

func (this *SnRequest) Update() (int64, error) {
	return Stmts.Update(`Update`, this)
}

func (this *Serialized) Read(arg interface{}) (error) {
	return Stmts.Get(`Read`, this, arg)
}

func (this *Serialized) JSON() ([]byte, error) {
	return json.Marshal(this)
}
