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
	`strings`
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
	Exec     string
	File     string
	Line     int
	Package  string
	Function string
	FilePath string
	ExecPath string
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

	// Build callerInfo attributes.

	exec := os.Args[0]
	filePath := filepath.Dir(frame.File)
	execPath := filepath.Dir(exec)

	if p, err := filepath.Abs(filePath); err == nil {
		filePath = p
	}

	if p, err := filepath.Abs(execPath); err == nil {
		execPath = p
	}

	pkgPath, funcBase := filepath.Split(frame.Function)
	pkgBase := strings.Split(funcBase, `.`)[0]
	pkgName := filepath.Join(pkgPath, pkgBase)

	return &callerInfo{
		Exec:     filepath.ToSlash(exec),
		File:     filepath.ToSlash(frame.File),
		Line:     frame.Line,
		Package:  filepath.ToSlash(pkgName),
		Function: filepath.ToSlash(frame.Function),
		FilePath: filepath.ToSlash(filePath),
		ExecPath: filepath.ToSlash(execPath),
	}
}
