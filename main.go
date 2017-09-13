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
	"flag"
	"log"
	"github.com/jscherff/goutils"
)

// Systemwide configuration.
var (
	conf *Config
	slog, elog, alog *goutils.MultiWriter
	db *Database
	ws *Server
)

// Systemwide initialization.
func init() {

	var err error

	flag.Parse()

	if conf, err = NewConfig(*FConfig); err != nil {
		log.Fatalf("%v", err)
	}

	if err = conf.Log.Init(); err != nil {
		log.Fatalf("%v", err)
	}

	slog = conf.Log.Writer["system"]
	alog = conf.Log.Writer["access"]
	elog = conf.Log.Writer["error"]

	if err = conf.Database.Init(); err != nil {
		slog.WriteError(err)
		log.Fatalf("%v", err)
	}

	conf.Server.Init()

	db = conf.Database
	ws = conf.Server

	slog.WriteString(db.Info())
	slog.WriteString(ws.Info())
}

func main() {
	log.Fatal(conf.Server.ListenAndServe())
	conf.Database.Close()
	conf.Log.Close()
}
