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
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

// Logger contains logger infomration and logging options. It is part of
// the systemwide configuration under Config.Logger.
type Logger struct {

	Writer map[string]*MultiWriter

	LogFile map[string]string

	Console map[string]bool

	Syslog map[string]*Syslog

	Options struct {
		EnableLogFile bool
		EnableConsole bool
		EnableSyslog bool
		RecoveryStack bool
	}

	LogDir struct {
		Windows string
		Linux string
	}
}

func (this *Logger) Close() {
	for _, v := range this.Writer {
		v.Close()
	}
}

func (this *Logger) Init() (err error) {

	if this.Options.EnableLogFile || *FLogFiles {

		var logDir string

		switch runtime.GOOS {

		case "windows":
			logDir = filepath.Dir(os.Args[0])
			logDir = filepath.Join(logDir, this.LogDir.Windows)
		case "linux":
			logDir = this.LogDir.Linux
		}

		for k, v := range this.LogFile {

			if this.Writer[k] == nil {
				continue
			}

			this.Writer[k].AddFile(filepath.Join(logDir, v))
		}
	}

	if this.Options.EnableSyslog || *FSyslog {

		for k, v := range this.Syslog {

			if this.Writer[k] == nil {
				continue
			}

			if err = v.Init(); err != nil {
				continue
			}

			this.Writer[k].AddWriter(v)
		}
	}

	if this.Options.EnableConsole || *FConsole {

		for k, v := range this.Console {

			if this.Writer[k] == nil {
				continue
			}

			if !v {
				continue
			}

			if k == "error" {
				this.Writer[k].AddConsole(os.Stderr)
			} else {
				this.Writer[k].AddConsole(os.Stdout)
			}
		}
	}

	for _, v := range this.Writer {
		if v.Count() == 0 {
			v.AddWriter(ioutil.Discard)
		}
	}

	return err
}


