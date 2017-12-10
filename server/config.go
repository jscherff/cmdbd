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


package service

import (
	`path/filepath`
	`time`

	`github.com/jscherff/cmdbd/common`
)

const (
	priKeyName	string = `PriKey`
	pubKeyName	string = `PubKey`
)


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

// Master configuration settings.
type Master struct {

	AuthMaxAge	time.Duration

	Console		bool
	Refresh		bool
	LogDir		string

	ConfigFile	map[string]string
	CryptoFile	map[string]string
	SerialFormat	map[string]string

	PublicKey	[]byte
	PrivateKey	[]byte

	AuthTokenSvc	AuthTokenService
	AuthCookieSvc	AuthCookieService
	SerialNumSvc	SerialNumService

	DataStore	store.DataStore
	MiddleWare	MiddleWare
	SystemLog	Logger
	AccessLog	Logger
	ErrorLog	Logger
	Syslog		Syslog

	//MetaUsb	*legacy.MetaUsb
	//Database	*legacy.Database
	//Queries	*legacy.Queries
	Router		*Router
	Server		*Server
}

// NewConfig creates a new master configuration object and reads its config
// from the provided JSON configuration file. 
func NewConfig(cf string, console, refresh bool) (*Config, error) {

	// ------------------------------
	// Load the master configuration.
	// ------------------------------

	this := &Config{Console: console, Refresh: refresh}

	if err := common.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	// Prepend the master config directory to other filenames.

	for key, fn := range this.ConfigFile {
		this.ConfigFile[key] = filepath.Join(filepath.Dir(cf), fn)
	}

	for key, fn := range this.CryptoFile {
		this.CryptoFile[key] = filepath.Join(filepath.Dir(cf), fn)
	}

	// Set the maximum age of auth cookies and tokens.

	this.AuthMaxAge *= time.Minute

	// --------------------------------------
	// Load the public and private key files.
	// --------------------------------------

	if pubKeyFile, ok := this.ConfigFile[pubKeyName]; !ok {
		return nil, fmt.Errorf(`public key config %q not found`, pubKeyName)
	} else if pemKey, err := ioutil.ReadFile(pubKeyFile); err != nil {
		return nil, err
	} else {
		this.PubKey = pemKey
	}

	if priKeyFile, ok := this.ConfigFile[priKeyName]; !ok {
		return nil, fmt.Errorf(`private key config %q not found`, priKeyName)
	} else if pemKey, err := ioutil.ReadFile(priKeyFile); err != nil {
		return nil, err
	} else {
		this.PriKey = pemKey
	}

	// -------------------------------
	// Create and initialize services.
	// -------------------------------

	if ts, err := NewAuthTokenService(this.pubKey, this.priKey, this.AuthMaxAge); err != nil {
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

	// ------------------------------------
	// Create and initialize the DataStore.
	// ------------------------------------

	if ds, err := store.New(`mysql`, this.ConfigFile[`DataStore`]); err != nil {
		return nil, err
	} else {
		this.DataStore = ds
	}

	// ----------------------------------------
	// Create and initialize the Syslog client.
	// ----------------------------------------

	if sl, err := NewSyslog(this); err != nil {
		return nil, err
	} else {
		this.Syslog = sl
	}

	// -----------------------------------------------------
	// Create and initialize the loggers and create aliases.
	// -----------------------------------------------------

	if ls, err := NewLoggers(this); err != nil {
		return nil, err
	} else {
		this.Loggers = ls
	}

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

	// ---------------------------------
	// Create and initialize MiddleWare.
	// ---------------------------------

	if mw, err := NewMiddleWare(this); err != nil {
		return nil, err
	} else {
		this.MiddleWare = mw
	}

	// -----------------------------
	// Create and initialize Router.
	// -----------------------------

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
	// -----------------------------
	// Create and initialize Server.
	// -----------------------------

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
