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
	"net/http"
	"flag"
	"log"
)

var conf *Config

func init() {

	var err error
	flag.Parse()

	if conf, err = NewConfig(*fWsConfigFile); err != nil {
		log.Fatalf("%v", err)
	}

	if err = conf.Log.Init(); err != nil {
		log.Printf("%v", err)
	}

	if err = conf.Database.Init(); err != nil {
		conf.Log.Writer[System].WriteError(err)
	}

	conf.Log.Writer[System].WriteString(conf.Database.Info())
}

func main() {

	router := NewRouter()
	log.Fatal(http.ListenAndServe(conf.ListenerInfo(), router))

	conf.Database.Close()
	conf.Log.Close()
}
