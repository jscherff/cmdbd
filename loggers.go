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
	`os`
	`path/filepath`
	`runtime`
	`github.com/jscherff/gox/log`
)

// Logger contains logger information and logging options. It is part of
// the systemwide configuration under Config.Loggers.
type Logger struct {
	*log.MLogger
	LogFile string
	LogFlags []string
	Stdout bool
	Stderr bool
	Syslog bool
}

// Loggers is the collection of application loggers. It is part of the
// systemwide configuration under Config.Loggers.
type Loggers map[string]*Logger

// Init performs initialization tasks for each logger in the systemwide
// configuration file.
func (this Loggers) Init() {

	var logDir string

	switch runtime.GOOS {

	case `windows`:
		logDir = filepath.Dir(os.Args[0])
		logDir = filepath.Join(logDir, conf.LogDir.Windows)
	case `linux`:
		logDir = conf.LogDir.Linux
	}

	for tag, logger := range this {

		flags := log.LoggerFlags(logger.LogFlags...)
		file := filepath.Join(logDir, logger.LogFile)

		stdout := logger.Stdout || conf.Options.Stdout || *FStdout
		stderr := logger.Stderr || conf.Options.Stderr || *FStderr

		logger.MLogger = log.NewMLogger(tag, flags, stdout, stderr, file)

		if logger.Syslog && conf.Syslog.Writer != nil {
			logger.AddWriter(conf.Syslog.Writer)
		}
	}
}

// Sync and/or close writers within each logger as necessary.
func (this Loggers) Close() {
	for _, logger := range this {
		logger.Close()
	}
}
