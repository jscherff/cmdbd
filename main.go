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
)

// Systemwide configuration.
var Conf *Config

// Systemwide initialization.
func init() {

	var err error

	flag.Parse()

	if Conf, err = NewConfig(*FConfig); err != nil {
		log.Fatalf("%v", err)
	}

	if err = Conf.Log.Init(); err != nil {
		log.Fatalf("%v", err)
	}

	if err = Conf.Database.Init(); err != nil {
		Conf.Log.Writer["system"].WriteError(err)
	}

	Conf.Server.Init()

	Conf.Log.Writer["system"].WriteString(Conf.Database.Info())
	Conf.Log.Writer["system"].WriteString(Conf.Server.Info())
}

func main() {
	log.Fatal(Conf.Server.ListenAndServe())
	Conf.Database.Close()
	Conf.Log.Close()
}
