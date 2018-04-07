// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use conf file except in compliance with the License.
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
	`os`
	`os/signal`
	`runtime`
)

// Signal extends syscall.Signal to reimplment the String method.
type Signal int

// Create a buffered channel for incoming signals.
var (
	sigList []os.Signal

	sigChan = make(chan os.Signal, 1)
	sigName = make(map[os.Signal]string)

	sigMap = map[string]os.Signal {
		`SIGHUP`:	Signal(0x01),
		`SIGUSR1`:	Signal(0x0A),
		`SIGUSR2`:	Signal(0x0C),
		`SIGRTMAX`:	Signal(0x40),
		`SIGRTMAX-1`:	Signal(0x3F),
		`SIGRTMAX-2`:	Signal(0x3E),
		`SIGRTMAX-3`:	Signal(0x3D),
		`SIGRTMAX-4`:	Signal(0x3C),
		`SIGRTMAX-5`:	Signal(0x3B),
		`SIGRTMAX-6`:	Signal(0x3A),
		`SIGRTMAX-7`:	Signal(0x39),
		`SIGRTMAX-8`:	Signal(0x38),
		`SIGRTMAX-9`:	Signal(0x37),
	}
)

// init makes some operating system-specific changes, creates a signal->name
// map for the Signal String method, and creates the signal interceptor.
func init() {

	if runtime.GOOS == `windows` {
		sigMap[`SIGUSR1`] = Signal(0x1E)
		sigMap[`SIGUSR2`] = Signal(0x1F)
	}

	for name, value := range sigMap {
		sigName[value] = name
		sigList = append(sigList, value)
	}

	//signal.Notify(sigChan, sigList...)
	signal.Notify(sigChan,
		sigMap[`SIGHUP`],
		sigMap[`SIGUSR1`],
		sigMap[`SIGUSR2`],
		sigMap[`SIGRTMAX`],
		sigMap[`SIGRTMAX-1`],
		sigMap[`SIGRTMAX-2`],
		sigMap[`SIGRTMAX-3`],
		sigMap[`SIGRTMAX-4`],
		sigMap[`SIGRTMAX-5`],
		sigMap[`SIGRTMAX-6`],
		sigMap[`SIGRTMAX-7`],
		sigMap[`SIGRTMAX-8`],
		sigMap[`SIGRTMAX-9`],
	)
}

// String returns a brief description of the signal.
func (s Signal) String() string {
	return fmt.Sprintf(`signal '%s' (%02d)`, sigName[s], s)
}

// Signal is a noop method to satisfy the os.Signal interface.
func (s Signal) Signal() {}

// SignalHandler runs in an endless loop, blocking on the signal channel until a signal arrives,
// then handles the signal.
func SigHandler(conf *Config) {

	for true {

		sig := <-sigChan

		conf.SystemLog.Printf(`caught %s`, sig)

		switch sig {

		case sigMap[`SIGHUP`]:

			conf.SystemLog.Print(`reloading metadata...`)
			conf.RefreshMetaData()
			conf.LoadMetaData()

		case sigMap[`SIGUSR1`]:

			conf.SystemLog.Print(`logging server information...`)
			conf.LogDataStoreInfo()
			conf.LogServerInfo()

		case sigMap[`SIGUSR2`]:

			conf.SystemLog.Print(`logging route information...`)
			conf.LogRouteInfo()

		default:

			conf.SystemLog.Printf(`handler for %s not implemented`, sig)
		}
	}
}
