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
	"log"
	"os"
)

var (
	db *Database
	conf *Config
	systemLog *MultiWriter
	accessLog *MultiWriter
	errorLog *MultiWriter
)

func init() {

	var err error
	flag.Parse()

	conf, err = NewConfig(*fWsConfigFile)

	if err != nil {
		log.Fatalf("%v", err)
	}

	systemLog = NewMultiWriter()
	accessLog = NewMultiWriter()
	errorLog = NewMultiWriter()

	if conf.EnableLogFiles || *fEnableLogFiles {

		if slf, alf, elf, err := conf.LogFileInfo(); err == nil {
			systemLog.AddFile(slf)
			accessLog.AddFile(alf)
			errorLog.AddFile(elf)
		} else {
			log.Printf("%v", err)
		}
	}

	if conf.EnableSyslog || *fEnableSyslog {
		proto, raddr, tag := conf.SyslogInfo()
		systemLog.AddSyslog(proto, raddr, tag, LogSystem)
		accessLog.AddSyslog(proto, raddr, tag, LogAccess)
		errorLog.AddSyslog(proto, raddr, tag, LogError)
	}

	if conf.EnableConsole || *fEnableConsole {
		systemLog.AddConsole(os.Stdout)
		accessLog.AddConsole(os.Stdout)
		errorLog.AddConsole(os.Stderr)
	}

	if systemLog.Count() == 0 {
		systemLog.AddWriter(ioutil.Discard)
	}

	if accessLog.Count() == 0 {
		accessLog.AddWriter(ioutil.Discard)
	}

	if errorLog.Count() == 0 {
		errorLog.AddWriter(ioutil.Discard)
	}

	db = conf.Database

	if err = db.Connect(); err != nil {
		systemLog.WriteError(err)
	}

	systemLog.WriteString(db.Info())
}

func main() {

	router := NewRouter()
	log.Fatal(http.ListenAndServe(conf.ListenerInfo(), router))

	systemLog.Close()
	accessLog.Close()
	errorLog.Close()
	db.Close()
}
