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
	`io/ioutil`
	`os`
	`path/filepath`
	`time`

	//`github.com/jscherff/cmdbd/legacy`
	`github.com/jscherff/cmdbd/common`
	`github.com/jscherff/cmdbd/service`
	`github.com/jscherff/cmdbd/store`

	`github.com/jscherff/cmdbd/model/cmdb`
	`github.com/jscherff/cmdbd/model/cmdb/usbci`
	`github.com/jscherff/cmdbd/model/cmdb/usbmeta`

	//v1cmdb `github.com/jscherff/cmdbd/api/v1/cmdb`
	//v1usbci `github.com/jscherff/cmdbd/api/v1/cmdb/usbci`
	//v1usbmeta `github.com/jscherff/cmdbd/api/v1/cmdb/usbmeta`

	//v2cmdb `github.com/jscherff/cmdbd/api/v2/cmdb`
	//v2usbci `github.com/jscherff/cmdbd/api/v2/cmdb/usbci`
	//v2usbmeta `github.com/jscherff/cmdbd/api/v2/cmdb/usbmeta`
)

const (
	priKeyName	string = `PriKey`
	pubKeyName	string = `PubKey`
)

// Global variables.
var (
	program		string = filepath.Base(os.Args[0])
	version		string = `undefined`
)

// Master configuration settings.
type Config struct {

	AuthMaxAge	time.Duration

	Console		bool
	Refresh		bool
	LogDir		string

	ConfigFile	map[string]string
	CryptoFile	map[string]string
	SerialFormat	map[string]string

	PublicKey	[]byte
	PrivateKey	[]byte

	AuthTokenSvc	service.AuthTokenService
	AuthCookieSvc	service.AuthCookieService
	SerialNumSvc	service.SerialNumService
	LoggerSvc	service.LoggerService

	DataStore	store.DataStore
	MiddleWare	MiddleWare
	SystemLog	service.Logger
	AccessLog	service.Logger
	ErrorLog	service.Logger
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
		this.PublicKey = pemKey
	}

	if priKeyFile, ok := this.ConfigFile[priKeyName]; !ok {
		return nil, fmt.Errorf(`private key config %q not found`, priKeyName)
	} else if pemKey, err := ioutil.ReadFile(priKeyFile); err != nil {
		return nil, err
	} else {
		this.PrivateKey = pemKey
	}

	// ----------------------------------------
	// Create and initialize the Syslog client.
	// ----------------------------------------

	if sl, err := NewSyslog(this); err != nil {
		return nil, err
	} else {
		this.Syslog = sl
	}

	// -------------------------------
	// Create and initialize services.
	// -------------------------------

	if ats, err := service.NewAuthTokenService(this.PublicKey, this.PrivateKey, this.AuthMaxAge); err != nil {
		return nil, err
	} else {
		this.AuthTokenSvc = ats
	}

	if acs, err := service.NewAuthCookieService(this.AuthMaxAge); err != nil {
		return nil, err
	} else {
		this.AuthCookieSvc = acs
	}

	if sns, err := service.NewSerialNumService(this.SerialFormat); err != nil {
		return nil, err
	} else {
		this.SerialNumSvc = sns
	}

	if ls, err := service.NewLoggerService(this.LogDir, this.Console, this.Syslog); err != nil {
		return nil, err
	} else {
		this.LoggerSvc = ls
	}

	// ----------------------------------
	// Create and initialize the Loggers.
	// ----------------------------------

	if sl, err := this.LoggerSvc.Create(this.ConfigFile[`SystemLog`]); err != nil {
		return nil, err
	} else {
		this.SystemLog = sl
	}

	if al, err := this.LoggerSvc.Create(this.ConfigFile[`AccessLog`]); err != nil {
		return nil, err
	} else {
		this.AccessLog = al
	}

	if el, err := this.LoggerSvc.Create(this.ConfigFile[`ErrorLog`]); err != nil {
		return nil, err
	} else {
		this.ErrorLog = el
	}

	// ------------------------------------
	// Create and initialize the DataStore.
	// ------------------------------------

	if ds, err := store.New(`mysql`, this.ConfigFile[`DataStore`]); err != nil {
		return nil, err
	} else {
		this.DataStore = ds
	}

	// --------------------------------
	// Create and initialize the Model.
	// --------------------------------

	if stmts, err := this.DataStore.Prepare(this.ConfigFile[`Queries`]); err != nil {
		return nil, err
	} else {
		cmdb.Init(stmts)
		usbci.Init(stmts)
		usbmeta.Init(stmts)
	}

	// ---------------------------------
	// Create and initialize MiddleWare.
	// ---------------------------------

	if mw, err := NewMiddleWare(this.AuthTokenSvc, this.AuthCookieSvc); err != nil {
		return nil, err
	} else {
		this.MiddleWare = mw
	}

	// -----------------------------
	// Create and initialize Router.
	// -----------------------------

	if rt, err := NewRouter(this.ConfigFile[`Router`], this.MiddleWare, this.AccessLog, this.ErrorLog); err != nil {
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

	if ws, err := NewServer(this.ConfigFile[`Server`], this.Router); err != nil {
		return nil, err
	} else {
		this.Server = ws
	}

	return this, nil
}

// displayVersion displays the program version.
func DisplayVersion() {
	fmt.Fprintf(os.Stderr, "%s version %s\n", program, version)
}
