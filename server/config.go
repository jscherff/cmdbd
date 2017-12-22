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

	`github.com/jscherff/cmdb/meta/peripheral`

	`github.com/jscherff/cmdbd/common`
	`github.com/jscherff/cmdbd/service`
	`github.com/jscherff/cmdbd/store`

	model_cmdb	`github.com/jscherff/cmdbd/model/cmdb`
	model_usbci	`github.com/jscherff/cmdbd/model/cmdb/usbci`
	model_usbmeta	`github.com/jscherff/cmdbd/model/cmdb/usbmeta`

	//v1cmdb `github.com/jscherff/cmdbd/api/v1/cmdb`
	//v1usbci `github.com/jscherff/cmdbd/api/v1/cmdb/usbci`
	//v1usbmeta `github.com/jscherff/cmdbd/api/v1/cmdb/usbmeta`

	api_cmdb_v2	`github.com/jscherff/cmdbd/api/v2/cmdb`
	api_usbci_v2	`github.com/jscherff/cmdbd/api/v2/cmdb/usbci`
	api_usbmeta_v2	`github.com/jscherff/cmdbd/api/v2/cmdb/usbmeta`
)

const (
	pubKeyName	string = `PublicKey`
	priKeyName	string = `PrivateKey`
)

var (
	program		string = filepath.Base(os.Args[0])
	version		string = `undefined`
)

// Master configuration settings.
type Config struct {

	AuthMaxAge	time.Duration

	Console		bool
	Refresh		bool

	ConfigFile	map[string]string
	CryptoFile	map[string]string
	SerialFormat	map[string]string

	PublicKey	[]byte
	PrivateKey	[]byte

	AuthSvc		service.AuthSvc
	SerialSvc	service.SerialSvc
	LoggerSvc	service.LoggerSvc
	DataStore	store.DataStore

	//Database	*legacy.Database
	//Queries	*legacy.Queries

	UsbMeta		*peripheral.Usb
	Syslog		*Syslog
	Router		*Router
	Server		*Server
}

// NewConfig creates a new master configuration object and reads its config
// from the provided JSON configuration file. 
func NewConfig(cf string, console, refresh bool) (*Config, error) {

	// ------------------------------
	// Load the master configuration.
	// ------------------------------

	this := &Config{}

	if err := common.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	this.Console = this.Console || console
	this.Refresh = this.Refresh || refresh

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

	if pubKeyFile, ok := this.CryptoFile[pubKeyName]; !ok {
		return nil, fmt.Errorf(`public key config %q not found`, pubKeyName)
	} else if pemKey, err := ioutil.ReadFile(pubKeyFile); err != nil {
		return nil, err
	} else {
		this.PublicKey = pemKey
	}

	if priKeyFile, ok := this.CryptoFile[priKeyName]; !ok {
		return nil, fmt.Errorf(`private key config %q not found`, priKeyName)
	} else if pemKey, err := ioutil.ReadFile(priKeyFile); err != nil {
		return nil, err
	} else {
		this.PrivateKey = pemKey
	}

	// ----------------------------------------
	// Create and initialize the Syslog client.
	// ----------------------------------------

	if sl, err := NewSyslog(this.ConfigFile[`Syslog`]); err != nil {
		return nil, err
	} else {
		this.Syslog = sl
	}

	// -------------------------------
	// Create and initialize services.
	// -------------------------------

	if as, err := service.NewAuthSvc(this.PublicKey, this.PrivateKey, this.AuthMaxAge); err != nil {
		return nil, err
	} else {
		this.AuthSvc = as
	}

	if ss, err := service.NewSerialSvc(this.SerialFormat); err != nil {
		return nil, err
	} else {
		this.SerialSvc = ss
	}

	if ls, err := service.NewLoggerSvc(this.ConfigFile[`LoggerSvc`], this.Console, this.Syslog); err != nil {
		return nil, err
	} else {
		this.LoggerSvc = ls
	}

	// -----------------------------------------------
	// Create and initialize the DataStore and models.
	// -----------------------------------------------

	if ds, err := store.NewMysqlDataStore(this.ConfigFile[`DataStore`]); err != nil {
		return nil, err
	} else if err := ds.Prepare(this.ConfigFile[`Queries`]); err != nil {
		return nil, err
	} else {
		model_cmdb.Init(ds)
		model_usbci.Init(ds)
		model_usbmeta.Init(ds)
		this.DataStore = ds
	}

	// -----------------------------------
	// Create and initialize USB metadata.
	// -----------------------------------

	if um, err := NewUsbMeta(this.ConfigFile[`UsbMeta`], refresh); err != nil {
		return nil, err
	} else {
		this.UsbMeta = um
	}

	// -----------------------------
	// Create and initialize Router.
	// -----------------------------

	if rt, err := NewRouter(this.ConfigFile[`Router`], this.AuthSvc, this.LoggerSvc); err != nil {
		return nil, err
	} else {
		this.Router = rt
	}

	// -----------------------------------------------------
	// Initialize API Endpoints and add Endpoints to Router.
	// -----------------------------------------------------

	api_cmdb_v2.Init(this.AuthSvc, this.LoggerSvc)
	api_usbci_v2.Init(this.AuthSvc, this.SerialSvc, this.LoggerSvc)
	api_usbmeta_v2.Init(this.UsbMeta, this.LoggerSvc)

	this.Router.
		AddEndpoints(api_cmdb_v2.Endpoints).
		AddEndpoints(api_usbci_v2.Endpoints).
		AddEndpoints(api_usbmeta_v2.Endpoints)

/*
	this.Router.
		AddRoutes(usbCiRoutes).
		AddRoutes(usbMetaRoutes).
		AddRoutes(cmdbAuthRoutes)

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
