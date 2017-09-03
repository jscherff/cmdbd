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
	"path/filepath"
	"encoding/json"
	"runtime"
	"fmt"
	"os"
)

type Config struct {

	useSyslog bool

	useLogFiles bool

	syslog struct {
		tag string
		port string
		protocol string
		hostname string
	}

	logFiles struct {
		windows struct {
			logDir string
			accessLog string
			errorLog string
		}
		linux struct {
			logDir string
			accessLog string
			errorLog string
		}
	}

	usbDbUrl string
}

func NewConfig(cf string) (this *Config, err error) {

	this = new(Config)
	appDir := filepath.Dir(os.Args[0])

	fh, err := os.Open(filepath.Join(appDir, cf))
	defer fh.Close()

	if err != nil {
		return this, err
	}

	jd := json.NewDecoder(fh)

	err = jd.Decode(&this)

	return this, err
}

func (this *Config) UseLogFiles() (bool) {
	return this.useLogFiles
}

func (this *Config) UseSyslog() (bool) {
	return this.useSyslog
}

func (this *Config) LogFileInfo() (afn, efn string, err error) {

	appDir := filepath.Dir(os.Args[0])

	switch runtime.GOOS {

	case "windows":
		dir := filepath.Join(appDir, this.logFiles.windows.logDir)
		afn = filepath.Join(dir, this.logFiles.windows.accessLog)
		efn = filepath.Join(dir, this.logFiles.windows.errorLog)

	case "linux":
		dir := this.logFiles.linux.logDir
		afn = filepath.Join(dir, this.logFiles.linux.accessLog)
		efn = filepath.Join(dir, this.logFiles.linux.errorLog)

	default:
		err = fmt.Errorf("operating system '%v' not supported", runtime.GOOS)
	}

	return afn, efn, err
}

func (this *Config) SyslogInfo() (proto, raddr, tag string) {

	raddr = fmt.Sprintf("%s:%s", this.syslog.hostname, this.syslog.port)
	return this.syslog.protocol, raddr, this.syslog.tag
}
