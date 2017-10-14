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
	`path/filepath`
	`github.com/jscherff/gox/log`
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
	*log.MLogger
	LogFile string
	LogFlags []string
	Stdout bool
	Stderr bool
	Syslog bool
}

// NewLogger creates and initializes a new Logger instance.
func NewLogger(cf string, sy *Syslog) (this *Logger, err error) {

	this = &Logger{}

	if err = loadConfig(this, cf); err != nil {
		return nil, err
	}

	for tag, lg := range this.Logs {

		flags := log.LoggerFlags(lg.LogFlags...)
		file := filepath.Join(this.LogDir, lg.LogFile)

		stdout := lg.Stdout || this.Stdout || *FStdout
		stderr := lg.Stderr || this.Stderr || *FStderr
		syslog := lg.Syslog || this.Syslog || *FSyslog

		lg.MLogger = log.NewMLogger(tag, flags, stdout, stderr, file)

		if syslog && sy.Writer != nil {
			lg.AddWriter(sy.Writer)
		}
	}

	return this, nil
}

// Sync and/or close writers within each logger as necessary.
func (this *Logger) Close() {
	for _, lg := range this.Logs {
		lg.Close()
	}
}
