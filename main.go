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
)

var conf *Config

// Systemwide initialization.
func init() {

	var err error

	log.SetFlags(log.Flags() | log.Lshortfile)

	flag.Parse()

	if conf, err = NewConfig(*FConfig); err != nil {
		log.Fatal(err)
	}

	if *FRefresh {
		if err := SaveUsbMeta(); err != nil {
			el.Fatal(err)
		} else {
			sl.Println(`USB Metadata refreshed.`)
			os.Exit(0)
		}
	}

	sl.Print(conf.Database.Info())
	sl.Print(conf.Server.Info())
}

func main() {
	log.Fatal(conf.Server.ListenAndServe())
	conf.Queries.Close()
	conf.Database.Close()
	conf.Logger.Close()
}
