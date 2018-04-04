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
	`syscall`
)

// Create a buffered channel for incoming signals.
var (
	SigChan = make(chan os.Signal, 1)
	SigMap = make(map[string]syscall.Signal)
)

// Create a map of signals appropriate for the operating system.
func init() {

	SigMap[`SIGHUP`] = syscall.Signal(0x01)
	SigMap[`SIGINT`] = syscall.Signal(0x02)

	switch runtime.GOOS {

	case `windows`:
		SigMap[`SIGUSR1`] = syscall.Signal(0x1E)
		SigMap[`SIGUSR2`] = syscall.Signal(0x1F)

	case `linux`:
		SigMap[`SIGUSR1`] = syscall.Signal(0x0A)
		SigMap[`SIGUSR2`] = syscall.Signal(0x0C)
	}

	signal.Notify(SigChan,
		SigMap[`SIGHUP`],
		SigMap[`SIGUSR1`],
		SigMap[`SIGUSR2`],
	)
}

// Create the signal handler. The handler runs in an endless
// loop, blocking on the signal channel until a signal arrives,
// then handles the signal.
func SigHandler(conf *Config) {

	for true {

		sig := <-SigChan

		switch sig {

		case SigMap[`SIGHUP`]:

			conf.SystemLog.Print(`caught SIGHUP`)

			if err := conf.RefreshMetaData(); err != nil {
				err = fmt.Errorf(`device metadata refresh failed: %v`, err)
				conf.ErrorLog.Print(err)
			} else {
				conf.SystemLog.Print(`device metadata refresh succeeded`)
			}

			if err := conf.LoadMetaData(); err != nil {
				err = fmt.Errorf(`data model metadata load failed: %v`, err)
				conf.ErrorLog.Print(err)
			} else {
				conf.SystemLog.Print(`data model metadata load succeeded`)
			}

		case SigMap[`SIGUSR1`]:

			conf.SystemLog.Print(`caught SIGUSR1`)
			conf.SystemLog.Printf(`current open database connections: %d`,
				conf.DataStore.GetOpenConns())

		case SigMap[`SIGUSR2`]:

			conf.SystemLog.Print(`caught SIGUSR2`)
			conf.SystemLog.Printf(`device metadata last updated %s`,
				conf.MetaUsbSvc.LastUpdate())
		}
	}
}
