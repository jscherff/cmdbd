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

package cmdb

import (
	`fmt`
	`time`
	`github.com/jscherff/cmdbd/store`
	`golang.org/x/crypto/bcrypt`
)

var dataStore store.DataStore

func Init(ds store.DataStore) {
	dataStore = ds
}

type Event struct {
	Id		int64		`db:"id,omitempty"             json:"id"`
	Code		int		`db:"code,omitempty"           json:"code"`
	Source		string		`db:"source,omitempty"         json:"source"`
	Description	string		`db:"description,omitempty"    json:"description"`
	HostName	string		`db:"host_name,omitempty"      json:"host_name"`
	RemoteAddr	string		`db:"remote_addr,omitempty"    json:"remote_addr"`
	EventDate	time.Time	`db:"event_date,omitempty"     json:"event_date"`
}

type Sequence struct {
	Ord		int64		`db:"ord,omitempty"            json:"ord"`
	IssueDate	time.Time	`db:"issue_date,omitempty"     json:"issue_date"`
}

type User struct {
	Id		int64		`db:"id,omitempty"             json:"id"`
	Username	string		`db:"username,omitempty"       json:"username"`
	Password	string		`db:"password,omitempty"       json:"password"`
	Created		time.Time	`db:"created,omitempty"        json:"created"`
	Locked		bool		`db:"locked,omitempty"         json:"locked"`
	Role		string		`db:"role,omitempty"           json:"role"`
}

type Info struct {
	Client		string		`json:"client,omitempty"`
	Server		string		`json:"host,omitempty"`
	Proto		string		`json:"proto,omitempty"`
	Method		string		`json:"method,omitempty"`
	Scheme		string		`json:"scheme,omitempty"`
	Path		string		`json:"path,omitempty"`
	Store		string		`json:"store,omitempty"`
}

// ----------------------
// Standard CRUD Methods.
// ----------------------

func (this *Event) Create() (id int64, err error) {
	this.Id, err = dataStore.Exec(`Insert`, this)
	return this.Id, err
}

func (this *Sequence) Create() (id int64, err error) {
	this.Ord, err = dataStore.Exec(`Insert`, this)
	return this.Ord, err
}

func (this *User) Create() (id int64, err error) {
	this.Id, err = dataStore.Exec(`Insert`, this)
	return this.Id, err
}

func (this *User) Read() (error) {
	return dataStore.Read(`SelectByUniqueId`, this, this)
}

func (this *Info) Read() (error) {
	this.Store = dataStore.String()
	return nil
}

// --------------------
// Speicalized Methods.
// --------------------

func (this *User) VerifyPassword(passwd string) (error) {
	return bcrypt.CompareHashAndPassword([]byte(this.Password), []byte(passwd))
}

func (this *User) VerifyAccess(/*TODO: provide and validate role*/) (error) {
	if this.Locked {
		return fmt.Errorf(`user %q account locked`, this.Username)
	} else {
		return nil
	}
}
