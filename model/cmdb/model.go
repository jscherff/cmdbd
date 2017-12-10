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
	`time`
	`github.com/jscherff/cmdbd/model`
	`golang.org/x/crypto/bcrypt`
)

var Model = model.New()

type Error struct {
	Id		interface{}	`db:"id,omitempty"`
	Code		int		`db:"code,omitempty"`
	Source		string		`db:"source,omitempty"`
	Description	string		`db:"description,omitempty"`
	EventDate	time.Time	`db:"event_date,omitempty"`
}

type Sequence struct {
	Ord		interface{}	`db:"ord,omitempty"`
	IssueDate	time.Time	`db:"issue_date,omitempty"`
}

type User struct {
	Id		interface{}	`db:"id,omitempty"`
	Username	string		`db:"username,omitempty"`
	Password	string		`db:"password,omitempty"`
	Created		time.Time	`db:"created,omitempty"`
	Locked		bool		`db:"locked,omitempty"`
	Role		string		`db:"role,omitempty"`
}

func (this *Error) Create() (int64, error) {
	return Model.Stmts().Insert(`Create`, this)
}

func (this *Sequence) Create() (int64, error) {
	return Model.Stmts().Insert(`Create`, this)
}

func (this *User) Create() (int64, error) {
	return Model.Stmts().Insert(`Create`, this)
}

func (this *User) Read(arg interface{}) (error) {
	return Model.Stmts().Get(`Read`, this, arg)
}

func (this *User) Verify(passwd string) (error) {
	return bcrypt.CompareHashAndPassword([]byte(this.Password), []byte(passwd))
}
