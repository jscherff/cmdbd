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
	`net/http`
	`time`
	`github.com/jscherff/cmdbd/common`
)

// server extends the http.Server object.
type Server struct {
	*http.Server
}

// NewServer creates and initializes a new Server instance.
func NewServer(cf string, handler http.Handler) (*Server, error) {

	this := &Server{}

	if err := common.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	this.Handler = handler
	this.ReadTimeout *= time.Second
	this.WriteTimeout *= time.Second

	return this, nil
}

// String provides identifying information about the server.
func (this *Server) String() (string) {
	return fmt.Sprintf(`http daemon listening on %q`, this.Addr)
}
