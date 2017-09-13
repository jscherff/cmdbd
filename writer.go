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
	"bytes"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"os"
)

// MultiWriter is an io.Writer that sends output to multiple destinations.
type MultiWriter struct {
	writers  []io.Writer
	consoles []*os.File
	files    []*os.File
}

// NewMultiWriter returns an initialized MultiWriter object.
func NewMultiWriter() (this *MultiWriter) {
	return new(MultiWriter)
}

// AddWriter appends a writer to a MultiWriter writer.
func (this *MultiWriter) AddWriter(w io.Writer) {
	this.writers = append(this.writers, w)
}

// AddFile appends a file to a MultiWriter writer.
func (this *MultiWriter) AddFile(f string) {

	var err error
	var h *os.File

	if err = os.MkdirAll(filepath.Dir(f), LogDirMode); err == nil {

		if h, err = os.OpenFile(f, LogFileFlags, LogFileMode); err == nil {
			this.files = append(this.files, h)
		}
	}

	if err != nil {
		log.Printf("%v", ErrorDecorator(err))
	}
}

// AddConsole appends a console to a MultiWriter writer. Consoles are
// treated separately as they shouldn't be closed on termination.
func (this *MultiWriter) AddConsole(h *os.File) {
	this.consoles = append(this.consoles, h)
}

// Write writes output to each writer in MultiWriter.
func (this *MultiWriter) Write(b []byte) (n int, err error) {

	var errs int

	b = bytes.TrimSuffix(b, []byte("\n"))
	b = bytes.TrimSuffix(b, []byte("\r"))

	for _, w := range this.writers {
		if n, err = w.Write(b); err != nil { errs++ }
	}

	b = append(b, byte('\n'))

	for _, c := range this.consoles {
		if n, err = c.Write(b); err != nil { errs++ }
	}

	for _, f := range this.files {
		if n, err = f.Write(b); err != nil { errs++ }
	}
	if errs > 0 {
		err = fmt.Errorf("%d write errors", errs)
		log.Printf("%v", ErrorDecorator(err))
	}

	return n, err
}

//WriteString converts string input to []byte and then calls Write.
func (this *MultiWriter) WriteString(s string) (n int, err error) {
	return this.Write([]byte(s))
}

//WriteError converts an error to []byte and then calls Write.
func (this *MultiWriter) WriteError(e error) {
	this.Println(e)
}

// Println writes to each writer in default format with trailing newline.
func (this *MultiWriter) Println(t ...interface{}) {

	for _, w := range this.writers {
		if _, err := fmt.Fprintln(w, t); err != nil {
			log.Printf("%v", ErrorDecorator(err))
		}
	}

	for _, c := range this.consoles {
		if _, err := fmt.Fprintln(c, t); err != nil {
			log.Printf("%v", ErrorDecorator(err))
		}
	}

	for _, f := range this.files {
		if _, err := fmt.Fprintln(f, t); err != nil {
			log.Printf("%v", ErrorDecorator(err))
		}
	}
}

// Count returns the number of writers in MultiWriter.
func (this *MultiWriter) Count() (n int) {
	return len(this.writers) + len(this.consoles) + len(this.files)
}

// Sync syncs underlying file and console writers in MultiWriter.
func (this *MultiWriter) Sync() {
	for _, c := range this.consoles{
		c.Sync()
	}
	for _, f := range this.files {
		f.Sync()
	}
}

// Close syncs and closes underlying file writers in MultiWriter.
func (this *MultiWriter) Close() {
	this.Sync()
	for _, f := range this.files {
		f.Close()
	}
}
