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
	`path/filepath`
	`github.com/jscherff/cmdbd/service`
	`github.com/jscherff/cmdbd/store`
	`github.com/jscherff/cmdbd/utils`

	model_cmdb	`github.com/jscherff/cmdbd/model/cmdb`
	model_usbci	`github.com/jscherff/cmdbd/model/cmdb/usbci`
	model_usbmeta	`github.com/jscherff/cmdbd/model/cmdb/usbmeta`

	api_cmdb_v1	`github.com/jscherff/cmdbd/api/v1/cmdb`
	api_usbci_v1	`github.com/jscherff/cmdbd/api/v1/cmdb/usbci`
	api_usbmeta_v1	`github.com/jscherff/cmdbd/api/v1/cmdb/usbmeta`

	api_cmdb_v2	`github.com/jscherff/cmdbd/api/v2/cmdb`
	api_usbci_v2	`github.com/jscherff/cmdbd/api/v2/cmdb/usbci`
	api_usbmeta_v2	`github.com/jscherff/cmdbd/api/v2/cmdb/usbmeta`
)

// Master configuration settings.
type Config struct {

	Console		bool
	Refresh		bool
	ConfigFile	map[string]string

	AuthSvc		service.AuthSvc
	SerialSvc	service.SerialSvc
	LoggerSvc	service.LoggerSvc
	MetaUsbSvc	service.MetaUsbSvc
	DataStore	store.DataStore

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

	if err := utils.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	this.Console = this.Console || console
	this.Refresh = this.Refresh || refresh

	// Prepend the master config directory to other filenames.

	for key, fn := range this.ConfigFile {
		this.ConfigFile[key] = filepath.Join(filepath.Dir(cf), fn)
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

	if as, err := service.NewAuthSvc(this.ConfigFile[`AuthSvc`]); err != nil {
		return nil, err
	} else {
		this.AuthSvc = as
	}

	if ss, err := service.NewSerialSvc(this.ConfigFile[`SerialSvc`]); err != nil {
		return nil, err
	} else {
		this.SerialSvc = ss
	}

	if ls, err := service.NewLoggerSvc(this.ConfigFile[`LoggerSvc`], this.Console, this.Syslog); err != nil {
		return nil, err
	} else {
		this.LoggerSvc = ls
	}

	if mus, err := service.NewMetaUsbSvc(this.ConfigFile[`MetaUsbSvc`], refresh); err != nil {
		return nil, err
	} else {
		this.MetaUsbSvc = mus
	}

	// ------------------------------------
	// Create and initialize the DataStore.
	// ------------------------------------

	if ds, err := store.NewMysqlDataStore(this.ConfigFile[`DataStore`]); err != nil {
		return nil, err
	} else if err := ds.SetPool(this.ConfigFile[`ConnPool`]); err != nil {
		return nil, err
	} else if err := ds.Prepare(this.ConfigFile[`Queries`]); err != nil {
		return nil, err
	} else {
		this.DataStore = ds
	}

	// ---------------------------------
	// Create and initialize the Router.
	// ---------------------------------

	if rt, err := NewRouter(this.ConfigFile[`Router`], this.AuthSvc, this.LoggerSvc); err != nil {
		return nil, err
	} else {
		this.Router = rt
	}

	// ------------------
	// Initialize Models.
	// ------------------

	model_cmdb.Init(this.DataStore)
	model_usbci.Init(this.DataStore)
	model_usbmeta.Init(this.DataStore)

	// ----------------------
	// Initialize API Routes.
	// ----------------------

	api_cmdb_v2.Init(this.AuthSvc, this.LoggerSvc)
	api_usbci_v2.Init(this.AuthSvc, this.SerialSvc, this.LoggerSvc)
	api_usbmeta_v2.Init(this.MetaUsbSvc, this.LoggerSvc)

	// ------------------------
	// Add Routes to Router.
	// ------------------------

	this.Router.
		AddRoutes(api_cmdb_v2.Routes).
		AddRoutes(api_usbci_v2.Routes).
		AddRoutes(api_usbmeta_v2.Routes)

	this.Router.
		AddRoutes(api_cmdb_v1.Routes).
		AddRoutes(api_usbci_v1.Routes).
		AddRoutes(api_usbmeta_v1.Routes)

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
