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
	RecoveryStack bool
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

// Init performs initialization tasks for each log.
func (this *Logger) Init() {

	for tag, l := range this.Logs {

		flags := l.LoggerFlags(l.LogFlags...)
		file := filepath.Join(this.LogDir, l.LogFile)

		stdout := l.Stdout || this.Stdout || *FStdout
		stderr := l.Stderr || this.Stderr || *FStderr
		syslog := l.Syslog || this.Syslog || *FSyslog

		l.MLogger = log.NewMLogger(tag, flags, stdout, stderr, file)

		if syslog && conf.Syslog.Writer != nil {
			l.AddWriter(conf.Syslog.Writer)
		}
	}
}

// Sync and/or close writers within each logger as necessary.
func (this *Logger) Close() {
	for _, l := range this.Logs {
		l.Close()
	}
}
