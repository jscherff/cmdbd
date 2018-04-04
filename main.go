// Copyright 2017 John Scherff
//
// Licensed under the Apache License, version 2.0 (the "License");
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
	`fmt`
	`log`
	`os`
	`path/filepath`
	`github.com/jscherff/cmdbd/server`
	`github.com/jscherff/cmdbd/model/cmdb/usbmeta`
)

// Systemwide configuration.
var (
	program	string = filepath.Base(os.Args[0])
	version	string = `undefined`
	config	*server.Config
)

// Systemwide initialization.
func init() {

	var err error

	flag.Parse()
	log.SetFlags(log.Flags() | log.Lshortfile)

	if *FVersion {
		fmt.Fprintf(os.Stderr, "%s version %s\n", program, version)
		os.Exit(0)
	}

	if config, err = server.NewConfig(*FConfig, *FConsole, *FRefresh); err != nil {
		log.Fatal(err)
	}

	if *FRefresh {
		if err := usbmeta.Load(config.MetaUsbSvc.Raw()); err != nil {
			config.LoggerSvc.ErrorLog().Fatal(err)
		} else {
			config.SystemLog.Println(`USB Metadata refreshed.`)
			os.Exit(0)
		}
	}

	config.SystemLog.Printf(`%s version %s started`, program, version)
}

func main() {
	defer config.LoggerSvc.Close()
	defer config.DataStore.Close()
	log.Fatal(config.Server.ListenAndServe())
}
