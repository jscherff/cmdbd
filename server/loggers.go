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
	`fmt`
	`path/filepath`
	`github.com/jscherff/gox/log`
	`github.com/jscherff/cmdbd/common`
)

// Loggers is a collection of Loggers with a getter method.
type Loggers interface {
	Logger(name string) (Logger, error)
	Close()
}

// Logger reimplements the log.MLogger interface.
type Logger interface {
	log.MLogger
}

// loggers is a collection of MLogger configurations with overrides.
type loggers struct {
	Stdout bool
	Stderr bool
	Syslog bool
	LogDir string
	LogFlags []string
	Config map[string]*logger
	logger map[string]Logger
}

// logger contains the configuration of a MLogger instance.
type logger struct {
	Stdout bool
	Stderr bool
	Syslog bool
	LogFile string
	LogFlags []string
}

// NewLoggers creates and initializes a new collection of Loggers.
func NewLoggers(conf *Config) (Loggers, error) {

	this := &loggers{}

	if err := common.LoadConfig(this, conf.ConfigFile[`Loggers`]); err != nil {
		return nil, err
	}

	syslog := conf.Syslog
	console := conf.Console

	for name, conf := range this.Config {

		flags :=
			log.LoggerFlags(this.LogFlags...) |
			log.LoggerFlags(conf.LogFlags...)

		file := filepath.Join(this.LogDir, conf.LogFile)

		stdout := this.Stdout || conf.Stdout || console
		stderr := this.Stderr || conf.Stderr

		this.logger[name] = log.NewMLogger(name, flags, stdout, stderr, file)

		if (this.Syslog || conf.Syslog) && (syslog != nil) {
			this.logger[name].AddWriter(syslog)
		}
	}

	return this, nil
}

func (this *loggers) Logger(name string) (Logger, error) {

	if logger, ok := this.logger[name]; !ok {
		return nil, fmt.Errorf(`%s logger not found`, name)
	} else {
		return logger, nil
	}
}

// Sync and/or close writers within each logger as necessary.
func (this *loggers) Close() {
	for _, logger := range this.logger {
		logger.Close()
	}
}
