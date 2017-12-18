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
	`os`

	`github.com/jscherff/cmdbd/legacy`
)

// Systemwide configuration.
var conf *legacy.Config

// Systemwide initialization.
func init() {

	var err error

	flag.Parse()
	log.SetFlags(log.Flags() | log.Lshortfile)

	if *FVersion {
		legacy.DisplayVersion()
		os.Exit(0)
	}

	if conf, err = legacy.NewConfig(*FConfig, *FConsole, *FRefresh); err != nil {
		log.Fatal(err)
	}

	if *FRefresh {
		if err := legacy.SaveUsbMeta(); err != nil {
			conf.Logger.Error.Fatal(err)
		} else {
			conf.Logger.System.Println(`USB Metadata refreshed.`)
			os.Exit(0)
		}
	}

	conf.Logger.System.Print(conf.Database.Info())
	conf.Logger.System.Print(conf.Server.Info())
}

func main() {
	defer conf.Logger.Close()
	defer conf.Database.Close()
	defer conf.Queries.Close()
	log.Fatal(conf.Server.ListenAndServe())
}
