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

package service

import (
	`io`
	`path/filepath`
	`github.com/jscherff/gox/log`
	`github.com/jscherff/cmdbd/common`
)

// Loggers is a collection of Loggers with a getter method.
type Loggers interface {
	System() log.MLogger
	Access() log.MLogger
	Error() log.MLogger
	Close()
}

// loggers is a collection of MLogger configurations with overrides.
type loggers struct {
	LogDir string
	Stdout bool
	Stderr bool
	Syslog bool
	Logger struct {
		System *logger
		Access *logger
		Error *logger
	}
}

// System returns the system logger.
func (this *loggers) System() (log.MLogger) {
	return this.Logger.System
}

// Access returns the access logger.
func (this *loggers) Access() (log.MLogger) {
	return this.Logger.Access
}

// Error returns the error logger.
func (this *loggers) Error() (log.MLogger) {
	return this.Logger.Error
}

// Close closes the system, access, and error loggers.
func (this *loggers) Close() {
	this.Logger.System.Close()
	this.Logger.Access.Close()
	this.Logger.Error.Close()
}

// logger contains the configuration of a MLogger instance.
type logger struct {
	log.MLogger
	Tag string
	Stdout bool
	Stderr bool
	Syslog bool
	LogFile string
	LogFlags []string
}


// NewLoggers creates and initializes a new collection of Loggers.
func NewLoggers(cf string, console bool, syslog io.Writer) (Loggers, error) {

	this := &loggers{}

	if err := common.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	init := func(l *logger) {

		flags := log.LoggerFlags(l.LogFlags...)
		file := filepath.Join(this.LogDir, l.LogFile)

		l.MLogger = log.NewMLogger(
			l.Tag,
			flags,
			l.Stdout || this.Stdout || console,
			l.Stderr || this.Stderr,
			file,
		)

		if (l.Syslog || this.Syslog) && syslog != nil {
			l.AddWriter(syslog)
		}
	}

	init(this.Logger.System)
	init(this.Logger.Access)
	init(this.Logger.Error)

	return this, nil
}
