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
	"strings"
	"fmt"
	"os"
)

// Config contains infomation about the server process and log writers.
type Config struct {

	Server struct {
		ListenerAddress string
		ListenerPort string
	}

	Database struct {
		Driver string
		Config string
	}

	Syslog struct {
		Tag string
		Port string
		Protocol string
		Hostname string
	}

	LogFiles struct {
		Windows struct {
			LogDir string
			AccessLog string
			ErrorLog string
		}
		Linux struct {
			LogDir string
			AccessLog string
			ErrorLog string
		}
	}

	EnableSyslog bool
	EnableLogFiles bool
	EnableConsole bool
	UsbDbUrl string
}

// NewConfig creates a new Config object and reads its configuration from
// the provided JSON configuration file.
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

// LogFileInfo builds and returns the full log file path based on the
// operating system and configuration information.
func (this *Config) LogFileInfo() (afn, efn string, err error) {

	appDir := filepath.Dir(os.Args[0])

	switch runtime.GOOS {

	case "windows":
		dir := filepath.Join(appDir, this.LogFiles.Windows.LogDir)
		afn = filepath.Join(dir, this.LogFiles.Windows.AccessLog)
		efn = filepath.Join(dir, this.LogFiles.Windows.ErrorLog)

	case "linux":
		dir := this.LogFiles.Linux.LogDir
		afn = filepath.Join(dir, this.LogFiles.Linux.AccessLog)
		efn = filepath.Join(dir, this.LogFiles.Linux.ErrorLog)

	default:
		err = fmt.Errorf("operating system '%v' not supported", runtime.GOOS)
	}

	return afn, efn, err
}

// SyslogInfo builds and returns syslog parameters from the configuration
// information.
func (this *Config) SyslogInfo() (proto, raddr, tag string) {

	raddr = strings.Join([]string{this.Syslog.Hostname, this.Syslog.Port}, ":")
	return this.Syslog.Protocol, raddr, this.Syslog.Tag
}
