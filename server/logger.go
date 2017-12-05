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

package server

import (
	`path/filepath`
	`github.com/jscherff/gox/log`
	`github.com/jscherff/cmdbd/utils`
)

// Logger contains logging options and a collection of logs.
type Logger struct {
	LogDir string
	Stdout bool
	Stderr bool
	Syslog bool
	Logs map[string]*Log
}

// Log is an instance of MultiLogger with embedded configuration.
type Log struct {
	log.MLogger
	LogFile string
	LogFlags []string
	Stdout bool
	Stderr bool
	Syslog bool
}

// NewLogger creates and initializes a new Logger instance.
func NewLogger(cf string, console bool, syslog *Syslog) (this *Logger, err error) {

	this = &Logger{}

	if err = utils.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	for tag, mlog := range this.Logs {

		flags := log.LoggerFlags(mlog.LogFlags...)
		file := filepath.Join(this.LogDir, mlog.LogFile)

		fStdout := mlog.Stdout || this.Stdout || console
		fStderr := mlog.Stderr || this.Stderr
		fSyslog := mlog.Syslog || this.Syslog

		mlog.MLogger = log.NewMLogger(tag, flags, fStdout, fStderr, file)

		if fSyslog && syslog.Writer != nil {
			mlog.AddWriter(syslog.Writer)
		}
	}

	return this, nil
}

// Sync and/or close writers within each logger as necessary.
func (this *Logger) Close() {
	for _, mlog := range this.Logs {
		mlog.Close()
	}
}
