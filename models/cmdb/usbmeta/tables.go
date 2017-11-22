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

import `time`

type Class struct {
	ClassID		string		`db:"class_id"`
	ClassDesc	string		`db:"class_desc"`
	LastUpdate	time.Time	`db:"last_update"`
}

type SubClass struct {
	ClassID		string		`db:"class_id"`
	SubClassID	string		`db:"subclass_id"`
	SubClassDesc	string		`db:"subclass_desc"`
	LastUpdate	time.Time	`db:"last_update"`
}

type Protocol struct {
	ClassID		string		`db:"class_id"`
	SubClassID	string		`db:"subclass_id"`
	ProtocolID	string		`db:"protocol_id"`
	ProtocolDesc	string		`db:"protocol_desc"`
	LastUpdate	time.Time	`db:"last_update"`
}

type Vendor struct {
	VendorID	string		`db:"vendor_id"`
	VendorName	string		`db:"vendor_name"`
	LastUpdate	time.Time	`db:"last_update"`
}

type Product struct {
	VendorID	string		`db:"vendor_id"`
	ProductID	string		`db:"product_id"`
	ProductName	string		`db:"product_name"`
	LastUpdate	time.Time	`db:"last_update"`
}
