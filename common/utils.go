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

package common

import (
	`encoding/json`
	`os`
	`path/filepath`
	`runtime`
)

// loadConfig loads a JSON configuration file into an object.
func LoadConfig(t interface{}, cf string) (error) {

	if fh, err := os.Open(cf); err != nil {
		return err
	} else {
		defer fh.Close()
		jd := json.NewDecoder(fh)
		err = jd.Decode(&t)
		return err
	}
}

// callerInfo encapsulates path and function elements of caller information.
type callerInfo struct {
	Base string
	Path string
	File string
	Line int
	Func string
}

// CallerInfo returns information about the caller to the caller.
func CallerInfo() (*callerInfo) {

	// Need space only for PC of caller.
	pc := make([]uintptr, 1)

	// Skip PC of Callers() and this function.
	n := runtime.Callers(2, pc)

	// Return nil if no PCs available.
	if n == 0 { return nil }

	// Obtain the caller frame.
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	// Obtain the path elements.
	ci := &callerInfo {
		File: filepath.Base(frame.File),
		Line: frame.Line,
		Func: frame.Function,
	}

	if base, err := filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		ci.Base = filepath.Dir(os.Args[0])
	} else {
		ci.Base = base
	}

	if path, err := filepath.Rel(ci.Base, filepath.Dir(frame.File)); err != nil {
		ci.Path = filepath.Dir(frame.File)
	} else {
		ci.Path = path
	}

	// Return caller info
	return ci
}
