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

	if conf.Logger.Syslog || *FSyslog {
		if err = conf.Syslog.Init(); err != nil {
			log.Println(err)
		}
	}

	conf.Logger.Init()

	slog = conf.Logger.Logs[`system`]
	alog = conf.Logger.Logs[`access`]
	elog = conf.Logger.Logs[`error`]

	if err = conf.Database.Init(); err != nil {
		elog.Fatal(err)
	}

	db = conf.Database

	if conf.MetaCi.Usb, err = NewMetaCi(conf.Configs. **************

	if *FRefresh {
		if err = conf.MetaCi.Usb.Init(conf.URLs.UsbMeta); err != nil {
			elog.Fatal(err)
		} else if err = conf.MetaCi.Usb.Save(conf.Files.UsbMeta); err != nil {
			elog.Fatal(err)
		}
	} else if err = conf.MetaCi.Usb.Load(conf.Files.UsbMeta); err != nil {
		elog.Fatal(err)
	}

	if *FRefreshDb {
		if err = SaveUsbMeta(); err != nil {
			elog.Print(err)
		}
	}

	conf.Server.Init()
	ws = conf.Server

	slog.Print(db.Info())
	slog.Print(ws.Info())
}

func main() {
	log.Fatal(conf.Server.ListenAndServe())
	slog.Print("shutting down")
	conf.Queries.Close()
	conf.Database.Close()
	conf.Logger.Close()
}
