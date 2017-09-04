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

// Source: https://thenewstack.io/make-a-restful-json-api-go/

package main

import (
	"net/http"
	"io/ioutil"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	db *Database
	config *Config
	systemLog *MultiWriter
	accessLog *MultiWriter
	errorLog *MultiWriter
)

func init() {

	var err error
	flag.Parse()

	config, err = NewConfig(*fWsConfigFile)

	if err != nil {
		log.Fatalf("%v", err)
	}

	systemLog = NewMultiWriter()
	accessLog = NewMultiWriter()
	errorLog = NewMultiWriter()

	if config.EnableLogFiles || *fEnableLogFiles {

		if slf, alf, elf, err := config.LogFileInfo(); err == nil {
			systemLog.AddFiles(slf)
			accessLog.AddFiles(alf)
			errorLog.AddFiles(elf)
		} else {
			log.Printf("%v", err)
		}
	}

	if config.EnableSyslog || *fEnableSyslog {
		proto, raddr, tag := config.SyslogInfo()
		systemLog.AddSyslog(proto, raddr, tag, LogSystem)
		accessLog.AddSyslog(proto, raddr, tag, LogAccess)
		errorLog.AddSyslog(proto, raddr, tag, LogError)
	}

	if config.EnableConsole || *fEnableConsole {
		systemLog.Add(os.Stdout)
		accessLog.Add(os.Stdout)
		errorLog.Add(os.Stderr)
	}

	if systemLog.Count() == 0 {
		systemLog.Add(ioutil.Discard)
	}

	if accessLog.Count() == 0 {
		accessLog.Add(ioutil.Discard)
	}

	if errorLog.Count() == 0 {
		errorLog.Add(ioutil.Discard)
	}

	if len(*fDbConfigFile) > 0 {
		db, err = NewDatabase(config.Database.Driver, *fDbConfigFile)
	} else {
		db, err = NewDatabase(config.Database.Driver, config.Database.Config)
	}

	if err != nil {
		systemLog.WriteString(fmt.Sprintf("%v", err))
	}

	systemLog.WriteString(db.Info)
}

func main() {

	router := NewRouter()
	listener := fmt.Sprintf("%s:%s", config.Server.ListenerAddress, config.Server.ListenerPort)
	log.Fatal(http.ListenAndServe(listener, router))

	if err := systemLog.Close(); err != nil {
		log.Printf("%v", err)
	}

	if err := accessLog.Close(); err != nil {
		log.Printf("%v", err)
	}

	if err := errorLog.Close(); err != nil {
		log.Printf("%v", err)
	}

	db.Close()
}
