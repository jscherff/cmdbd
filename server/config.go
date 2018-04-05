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
	`net/http`
	`path/filepath`
	`strings`
	`time`
	`github.com/gorilla/handlers`
	`github.com/gorilla/mux`
	`github.com/jscherff/cmdbd/service`
	`github.com/jscherff/cmdbd/store`
	`github.com/jscherff/cmdbd/utils`
	`github.com/jscherff/gox/log`

	model_cmdb	`github.com/jscherff/cmdbd/model/cmdb`
	model_usbci	`github.com/jscherff/cmdbd/model/cmdb/usbci`
	model_usbmeta	`github.com/jscherff/cmdbd/model/cmdb/usbmeta`

	api_cmdb_v1	`github.com/jscherff/cmdbd/api/v1/cmdb`
	api_usbci_v1	`github.com/jscherff/cmdbd/api/v1/cmdb/usbci`
	api_usbmeta_v1	`github.com/jscherff/cmdbd/api/v1/cmdb/usbmeta`

	api_cmdb_v2	`github.com/jscherff/cmdbd/api/v2/cmdb`
	api_usbci_v2	`github.com/jscherff/cmdbd/api/v2/cmdb/usbci`
	api_usbmeta_v2	`github.com/jscherff/cmdbd/api/v2/cmdb/usbmeta`

	api_cmdb_v3	`github.com/jscherff/cmdbd/api/v3/cmdb`
	api_usbci_v3	`github.com/jscherff/cmdbd/api/v3/cmdb/usbci`
	api_usbmeta_v3	`github.com/jscherff/cmdbd/api/v3/cmdb/usbmeta`
)

// Message for server timeout middleware.
const timeoutMessage = `server timed out waiting for available connection`

// Master configuration settings.
type Config struct {

	Console		bool
	Refresh		bool
	RecoveryStack	bool
	MaxConnections	int
	ServerTimeout	time.Duration
	ConfigFile	map[string]string

	AuthSvc		service.AuthSvc
	SerialSvc	service.SerialSvc
	LoggerSvc	service.LoggerSvc
	MetaUsbSvc	service.MetaUsbSvc
	DataStore	store.DataStore

	AccessLog	log.MLogger
	SystemLog	log.MLogger
	ErrorLog	log.MLogger

	Syslog		*Syslog
	Router		*Router
	Server		*Server
}

// NewConfig creates a new master configuration object and reads its config
// from the provided JSON configuration file. 
func NewConfig(cf string, console, refresh bool) (*Config, error) {

	// -----------------------------------
	// Create a Config with sane defaults.
	// -----------------------------------

	this := &Config{
		MaxConnections: 50,
		ServerTimeout: 60,
	}

	// ------------------------------
	// Load the master configuration.
	// ------------------------------

	if err := utils.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	this.ServerTimeout *= time.Second

	console = this.Console || console
	refresh = this.Refresh || refresh

	// -----------------------------------------------
	// Prepend master config directory to other paths.
	// -----------------------------------------------

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

	if ls, err := service.NewLoggerSvc(this.ConfigFile[`LoggerSvc`], console, this.Syslog); err != nil {
		return nil, err
	} else {
		this.LoggerSvc = ls
		this.AccessLog = ls.AccessLog()
		this.SystemLog = ls.SystemLog()
		this.ErrorLog = ls.ErrorLog()
	}

	this.SystemLog.Print(`logging service initialized`)

	if as, err := service.NewAuthSvc(this.ConfigFile[`AuthSvc`]); err != nil {
		return nil, err
	} else {
		this.AuthSvc = as
	}

	this.SystemLog.Print(`authentication service initialized`)

	if ss, err := service.NewSerialSvc(this.ConfigFile[`SerialSvc`]); err != nil {
		return nil, err
	} else {
		this.SerialSvc = ss
	}

	this.SystemLog.Print(`serial number service initialized`)

	if mus, err := service.NewMetaUsbSvc(this.ConfigFile[`MetaUsbSvc`]); err != nil {
		return nil, err
	} else {
		this.MetaUsbSvc = mus
	}

	this.SystemLog.Print(`device metadata service initialized`)

	// -------------------------------------
	// Refresh Device Metadata if Requested.
	// -------------------------------------

	if refresh { this.RefreshMetaData() }
	this.SystemLog.Printf(`device metadata last updated %s`, this.MetaUsbSvc.LastUpdate())

	// ------------------------------------
	// Create and initialize the DataStore.
	// ------------------------------------

	if ds, err := store.NewMysqlDataStore(this.ConfigFile[`DataStore`]); err != nil {
		return nil, err
	} else if err = ds.SetConnPool(this.ConfigFile[`ConnPool`]); err != nil {
		return nil, err
	} else if err = ds.Prepare(this.ConfigFile[`Queries`]); err != nil {
		return nil, err
	} else {
		this.DataStore = ds
	}

	this.SystemLog.Printf(`datastore initialized`)
	this.LogDataStoreInfo()

	// -------------------------------------
	// Create and Initialize Request Router.
	// -------------------------------------

	if r, err := NewRouter(this.AuthSvc); err != nil {
		return nil, err
	} else {
		this.Router = r
	}

	this.SystemLog.Print(`request router initialized`)

	// ------------------
	// Initialize Models.
	// ------------------

	model_cmdb.Init(this.DataStore)
	model_usbci.Init(this.DataStore)
	model_usbmeta.Init(this.DataStore)

	this.SystemLog.Print(`data models initialized`)

	// ----------------------------------
	// Load Device Metadata if Requested.
	// ----------------------------------

	if refresh { this.LoadMetaData() }

	// ----------------------
	// Initialize API Routes.
	// ----------------------

	api_cmdb_v3.Init(this.AuthSvc, this.SystemLog, this.ErrorLog)
	api_usbci_v3.Init(this.AuthSvc, this.SerialSvc, this.SystemLog, this.ErrorLog)
	api_usbmeta_v3.Init(this.MetaUsbSvc, this.SystemLog, this.ErrorLog)

	this.SystemLog.Print(`route endpoints initialized`)

	// -----------------------------
	// Add Routes to Request Router.
	// -----------------------------

	this.Router.
		AddRoutes(api_cmdb_v3.Routes).
		AddRoutes(api_usbci_v3.Routes).
		AddRoutes(api_usbmeta_v3.Routes).
		AddRoutes(api_cmdb_v2.Routes).
		AddRoutes(api_usbci_v2.Routes).
		AddRoutes(api_usbmeta_v2.Routes).
		AddRoutes(api_cmdb_v1.Routes).
		AddRoutes(api_usbci_v1.Routes).
		AddRoutes(api_usbmeta_v1.Routes)

	this.SystemLog.Print(`route endpoints loaded`)
	this.LogRouteInfo()

	// ---------------------------
	// Chain Middleware to Routes.
	// ---------------------------

	var handler http.Handler

	// ------------------------------------------------
	// Prepend Max Connection Handler to Route Handler.
	// ------------------------------------------------

	handler = MaxConnectionHandler(this.Router, this.MaxConnections)

	this.SystemLog.Print(`max connection handler initialized`)

	// ---------------------------------------------------------
	// Prepend Server Timeout Handler to Max Connection Handler.
	// ---------------------------------------------------------

	handler = http.TimeoutHandler(handler, this.ServerTimeout, timeoutMessage)

	this.SystemLog.Print(`connection timeout handler initialized`)

	// ---------------------------------------------------
	// Prepend Recovery Handler to Server Timeout Handler.
	// ---------------------------------------------------

	recoveryHandler := handlers.RecoveryHandler(
		handlers.PrintRecoveryStack(this.RecoveryStack),
		handlers.RecoveryLogger(this.ErrorLog))

	handler = recoveryHandler(handler)

	this.SystemLog.Print(`recovery handler initialized`)

	// --------------------------------------------------------------------
	// Prepend Logging Handler Prepend Logging Handler to Recovery Handler.
	// --------------------------------------------------------------------

	handler = handlers.CombinedLoggingHandler(this.AccessLog, handler)

	this.SystemLog.Print(`access log handler initialized`)

	// -----------------------------
	// Create and initialize Server.
	// -----------------------------

	if server, err := NewServer(this.ConfigFile[`Server`], handler); err != nil {
		return nil, err
	} else {
		this.Server = server
	}

	this.SystemLog.Print(`server initialized`)
	this.LogServerInfo()

	// -------------------------
	// Start the signal handler.
	// -------------------------

	go SigHandler(this)

	return this, nil
}

// RefreshMetaData downloads a fresh copy of the device metadata.
func (this *Config) RefreshMetaData() {
	if err := this.MetaUsbSvc.Refresh(); err != nil {
		err = fmt.Errorf(`device metadata refresh failed: %v`, err)
		this.ErrorLog.Print(err)
	} else if err := this.MetaUsbSvc.Save(); err != nil {
		err = fmt.Errorf(`device metadata save failed: %v`, err)
		this.ErrorLog.Print(err)
	} else {
		this.SystemLog.Print(`device metadata refresh and save succeeded`)
	}
}

// LoadMetaData loads the metadata tables in the datastore.
func (this *Config) LoadMetaData() {
	if err := model_usbmeta.Load(this.MetaUsbSvc.Raw()); err != nil {
		err = fmt.Errorf(`data model metadata load failed: %v`, err)
		this.ErrorLog.Print(err)
	} else {
		this.SystemLog.Print(`data model metadata load succeeded`)
	}
}

// LogDataStoreInfo logs information about the datastore to the system log.
func (this *Config) LogDataStoreInfo() {

	connPool := this.DataStore.GetConnPool()

	this.SystemLog.Printf(`datastore driver %s`, this.DataStore)
	this.SystemLog.Printf(`datastore maximum open connections set to %d`, connPool.MaxOpenConns)
	this.SystemLog.Printf(`datastore maximum idle connections set to %d`, connPool.MaxIdleConns)
	this.SystemLog.Printf(`datastore connection maximum lifetime set to %s`, connPool.ConnMaxLifetime)
	this.SystemLog.Printf(`datastore current open connections: %d`, this.DataStore.GetOpenConns())
}

// LogServerInfo logs information about the server to the system log.
func (this *Config) LogServerInfo() {

	this.SystemLog.Printf(`server listening on %s`, this.Server.Addr)
	this.SystemLog.Printf(`server read timeout set to %s`, this.Server.ReadTimeout)
	this.SystemLog.Printf(`server write timeout set to %s`, this.Server.WriteTimeout)
	this.SystemLog.Printf(`server connection timeout set to %s`, this.ServerTimeout)
	this.SystemLog.Printf(`server maximum connections set to %d`, this.MaxConnections)
}

// LogRouteInfo logs information about API endpoint routes.
func (this *Config) LogRouteInfo() {

	this.Router.Walk(func(rt *mux.Route, rtr *mux.Router, anc []*mux.Route) error {

		var methods, path string

		if m, err := rt.GetMethods(); err != nil {
			methods = ``
		} else {
			methods = strings.Join(m, `|`)
		}
		if p, err := rt.GetPathTemplate(); err != nil {
			path = ``
		} else {
			path = p
		}
		this.SystemLog.Printf(`route '%s' template '%s %s'`,
			rt.GetName(),
			methods,
			path,
		)
		return nil
	})
}

