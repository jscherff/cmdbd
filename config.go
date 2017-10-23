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

package main

import (
	`encoding/json`
	`fmt`
	`path/filepath`
	`os`
)

var (
	// Program name and version.

	program = filepath.Base(os.Args[0])
	version = `undefined`

	// Configuration aliases.

	db *Database
	qy *Queries
	ws *Server
	sl, al, el *Log
)

// Config contains infomation about the server process and log writers.
type Config struct {

	SerialFmt string
	Configs   map[string]string

	Database *Database
	Queries  *Queries
	Syslog   *Syslog
	Logger   *Logger
	Router   *Router
	MetaUsb  *MetaUsb
	Server   *Server
}

// NewConfig creates a new Config object and reads its config
// from the provided JSON configuration file.
func NewConfig(cf string) (this *Config, err error) {

	// Load the base config needed to load remaining configs.

	this = &Config{}

	if err := loadConfig(this, cf); err != nil {
		return nil, err
	}

	// Prepend the base config directory to other config filenames.

	for key, fn := range this.Configs {
		this.Configs[key] = filepath.Join(filepath.Dir(cf), fn)
	}

	// Create and initialize Database object.

	if database, err := NewDatabase(this.Configs[`Database`]); err != nil {
		return nil, err
	} else {
		this.Database = database
	}

	db = this.Database

	// Create and initialize Queries object.

	if queries, err := NewQueries(this.Configs[`Queries`], db); err != nil {
		return nil, err
	} else {
		this.Queries = queries
	}

	qy = this.Queries

	// Create and initialize Syslog object.

	if syslog, err := NewSyslog(this.Configs[`Syslog`]); err != nil {
		return nil, err
	} else {
		this.Syslog = syslog
	}

	// Create and initialize Logger object.

	if logger, err := NewLogger(this.Configs[`Logger`], *FConsole, this.Syslog); err != nil {
		return nil, err
	} else {
		this.Logger = logger
	}

	sl = this.Logger.Logs[`system`]
	al = this.Logger.Logs[`access`]
	el = this.Logger.Logs[`error`]

	// Create and initialize Router object.

	if router, err := NewRouter(this.Configs[`Router`], al, el); err != nil {
		return nil, err
	} else {
		router.AddRoutes(usbCiRoutes).AddRoutes(usbMetaRoutes)
		this.Router = router
	}

	// Create and initialize MetaUsb object.

	if metausb, err := NewMetaUsb(this.Configs[`MetaUsb`], *FRefresh); err != nil {
		return nil, err
	} else {
		this.MetaUsb = metausb
	}

	// Create and initialize Server object.

	if server, err := NewServer(this.Configs[`Server`]); err != nil {
		return nil, err
	} else {
		server.Handler = this.Router
		this.Server = server
	}

	ws = this.Server

	return this, nil
}

// loadConfig loads a JSON configuration file into an object.
func loadConfig(t interface{}, cf string) error {

	if fh, err := os.Open(cf); err != nil {
		return err
	} else {
		defer fh.Close()
		jd := json.NewDecoder(fh)
		err = jd.Decode(&t)
		return err
	}
}

// displayVersion displays the program version.
func displayVersion() {
	fmt.Fprintf(os.Stderr, "%s version %s\n", program, version)
}
