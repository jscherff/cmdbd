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
	`net/http`
	`time`
	`github.com/jscherff/cmdbd/utils`
)

// server extends the http.Server object.
type Server struct {
	*http.Server
}

// NewServer creates and initializes a new Server instance.
func NewServer(cf string, handler http.Handler) (*Server, error) {

	// -----------------------------------
	// Create a Server with sane defaults.
	// -----------------------------------

	this := &Server{}

	// ------------------------------
	// Load the server configuration.
	// ------------------------------

	if err := utils.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	this.ReadTimeout *= time.Second
	this.WriteTimeout *= time.Second
	this.Handler = handler

	return this, nil
}
