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

// LoggerService is a collection of LoggerService with a getter method.
type LoggerService interface {
	Create(configFile string) (Logger, error)
	Close()
}

// Logger reimplements the log.MLogger interface.
type Logger interface {
	log.MLogger
}

// loggerService is a collection of MLogger configurations with overrides.
type loggerService struct {
	LogDir string
	Stdout bool
	Syslog io.Writer
	Logger []Logger
}

// logger contains the configuration of a MLogger instance.
type logger struct {
	Stdout bool
	Stderr bool
	Syslog bool
	LogTag string
	LogFile string
	LogFlags []string
}

// NewLoggerService creates and initializes a new collection of LoggerService.
func NewLoggerService(logDir string, console bool, syslog io.Writer) (LoggerService, error) {
	return &loggerService{LogDir: logDir, Stdout: console, Syslog: syslog}, nil
}

// Create creates a new logger from the provided JSON configuration file.
func (this *loggerService) Create(cf string) (Logger, error) {

	conf := &logger{}

	if err := common.LoadConfig(conf, cf); err != nil {
		return nil, err
	}

	newLogger := log.NewMLogger(
		conf.LogTag,
		log.LoggerFlags(conf.LogFlags...),
		this.Stdout || conf.Stdout,
		conf.Stderr,
		filepath.Join(this.LogDir, conf.LogFile),
	)

	if conf.Syslog && this.Syslog != nil {
		newLogger.AddWriter(this.Syslog)
	}

	this.Logger = append(this.Logger, newLogger)

	return newLogger, nil
}

// Sync and/or close writers within each logger as necessary.
func (this *loggerService) Close() {
	for _, logger := range this.Logger {
		logger.Close()
	}
}
