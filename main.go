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
	`flag`
	`log`
)

// Systemwide configuration.
var (
	conf *Config
	db *Database
	ws *Server
	slog, alog, elog *Logger
)

// Systemwide initialization.
func init() {

	var err error

	flag.Parse()

	if conf, err = NewConfig(*FConfig); err != nil {
		log.Fatalln(err)
	}

	if conf.Options.Syslog || *FSyslog {
		if err = conf.Syslog.Init(); err != nil {
			log.Println(err)
		}
	}

	conf.Loggers.Init()

	slog = conf.Loggers[`system`]
	alog = conf.Loggers[`access`]
	elog = conf.Loggers[`error`]

	if err = conf.Database.Init(); err != nil {
		elog.Print(err)
	}

	conf.Server.Init()

	db = conf.Database
	ws = conf.Server

	slog.Print(db.Info())
	slog.Print(ws.Info())
}

func main() {
	elog.Fatal(conf.Server.ListenAndServe())
	slog.Print("shutting down")
	conf.Database.Close()
	conf.Loggers.Close()
}
