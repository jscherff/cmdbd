// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use conf file except in compliance with the License.
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

import (
	`fmt`
	`os`
	`path/filepath`

	//`github.com/jscherff/cmdbd/legacy`
	`github.com/jscherff/cmdbd/common`
	`github.com/jscherff/cmdbd/server`
	`github.com/jscherff/cmdbd/service`
	`github.com/jscherff/cmdbd/store`

	//`github.com/jscherff/cmdbd/model`
	//`github.com/jscherff/cmdbd/model/cmdb`
	//`github.com/jscherff/cmdbd/model/cmdb/usbci`
	//`github.com/jscherff/cmdbd/model/cmdb/usbmeta`

	//v1cmdb `github.com/jscherff/cmdbd/api/v1/cmdb`
	//v1usbci `github.com/jscherff/cmdbd/api/v1/cmdb/usbci`
	//v1usbmeta `github.com/jscherff/cmdbd/api/v1/cmdb/usbmeta`

	//v2cmdb `github.com/jscherff/cmdbd/api/v2/cmdb`
	//v2usbci `github.com/jscherff/cmdbd/api/v2/cmdb/usbci`
	//v2usbmeta `github.com/jscherff/cmdbd/api/v2/cmdb/usbmeta`
)

// Global variables.
var (
	program		string = filepath.Base(os.Args[0])
	version		string = `undefined`
)

// Shared configurations and services.
type Config struct {

	AuthMaxAge	time.Duration
	AuthTokenSvc	service.AuthTokenService
	AuthCookieSvc	service.AuthCookieService
	SerialNumSvc	service.SerialNumService

	DataStore	store.DataStore
	MiddleWare	server.MiddleWare
	Loggers		server.Loggers
	SystemLog	server.Logger
	AccessLog	server.Logger
	ErrorLog	server.Logger
	Syslog		server.Syslog

	//MetaUsb		*legacy.MetaUsb
	//Database	*legacy.Database
	//Queries		*legacy.Queries
	Router		*server.Router
	Server		*server.Server

	ConfigFile	map[string]string
	CryptoFile	map[string]string
	SerialFormat	map[string]string

	Console		bool
	Refresh		bool
}

// NewConfig creates a new Config object and reads its config
// from the provided JSON configuration file.
func NewConfig(cf string, console, refresh bool) (*Config, error) {

	// Load the base config needed to load remaining configs.

	this := &Config{Console: console, Refresh: refresh}

	if err := common.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	// Prepend the master config directory to other filenames.

	for key, fn := range this.ConfigFile {
		this.ConfigFile[key] = filepath.Join(filepath.Dir(cf), fn)
	}

	// Create and initialize the DataStore object.

	if ds, err := store.New(`mysql`, this.ConfigFile[`DataStore`]); err != nil {
		return nil, err
	} else {
		this.DataStore = ds
	}

	// Create and initialize Syslog object.

	if sl, err := NewSyslog(this); err != nil {
		return nil, err
	} else {
		this.Syslog = sl
	}

	if ls, err := NewLoggers(this); err != nil {
		return nil, err
	} else {
		this.Loggers = ls
	}

	// Retrieve Loggers.

	if sl, err := this.Loggers.Logger(`System`); err != nil {
		return nil, err
	} else {
		this.SystemLog = sl
	}

	if al, err := this.Loggers.Logger(`Access`); err != nil {
		return nil, err
	} else {
		this.AccessLog = al
	}

	if el, err := this.Loggers.Logger(`Error`); err != nil {
		return nil, err
	} else {
		this.ErrorLog = el
	}

	// Create and initialize MiddleWare object.

	if mw, err := NewMiddleWare(this); err != nil {
		return nil, err
	} else {
		this.MiddleWare = mw
	}

	// Create and initialize Router object.

	if rt, err := NewRouter(this); err != nil {
		return nil, err
	} else {
		this.Router = rt
	}
/*
	this.Router.
		AddRoutes(usbCiRoutes).
		AddRoutes(usbMetaRoutes).
		AddRoutes(cmdbAuthRoutes)

	// Create and initialize MetaUsb object.

	if mu, err := NewMetaUsb(this); err != nil {
		return nil, err
	} else {
		this.MetaUsb = mu
	}

	// Create and initialize Database object.

	if db, err := NewDatabase(this); err != nil {
		return nil, err
	} else {
		this.Database = db
	}

	// Create and initialize Queries object.

	if qs, err := NewQueries(this); err != nil {
		return nil, err
	} else {
		this.Queries = qs
	}
*/
	// Create and initialize Server object.

	if ws, err := NewServer(this); err != nil {
		return nil, err
	} else {
		ws.Handler = this.Router
		this.Server = ws
	}

	return this, nil
}

// displayVersion displays the program version.
func DisplayVersion() {
	fmt.Fprintf(os.Stderr, "%s version %s\n", program, version)
}
// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use conf file except in compliance with the License.
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
	`path/filepath`
	`time`

	`github.com/jscherff/cmdbd/common`
)

// Constants.
const (
	priKeyName	string = `PriKey`
	pubKeyName	string = `PubKey`
)

// Shared configurations and services.
type Config struct {
}

// NewConfig creates a new Config object and reads its config
// from the provided JSON configuration file.
func NewConfig(cf string) (*Config, error) {

	// Load the base config needed to load remaining configs.

	this := &Config{}

	if err := common.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	// Set the maximum age of auth cookies and tokens.

	this.AuthMaxAge *= time.Minute

	// Prepend the master config directory to other filenames.

	for key, fn := range this.ConfigFile {
		this.ConfigFile[key] = filepath.Join(filepath.Dir(cf), fn)
	}

	for key, fn := range this.CryptoFile {
		this.CryptoFile[key] = filepath.Join(filepath.Dir(cf), fn)
	}

	// Initialize services.

	if ts, err := NewAuthTokenService(this); err != nil {
		return nil, err
	} else {
		this.AuthTokenSvc = ts
	}

	if cs, err := NewAuthCookieService(this); err != nil {
		return nil, err
	} else {
		this.AuthCookieSvc = cs
	}

	if ss, err := NewSerialNumService(this); err != nil {
		return nil, err
	} else {
		this.SerialNumSvc = ss
	}

	return this, nil
}
