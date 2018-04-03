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
	`net`
	`net/http`
	`time`
	`golang.org/x/net/netutil`
	`github.com/jscherff/cmdbd/utils`
)

// tcpKeepAliveListener wraps net.TCPListener and extends the Accept()
// method by implementing TCP Keepalives on the TCP Connection it returns.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept implements the Accept method in the Listener interface; it waits
// for the next call and returns a generic Connection. Extended to enable
// TCP Keepalives.
func (this tcpKeepAliveListener) Accept() (net.Conn, error) {
	if tc, err := this.AcceptTCP(); err != nil {
		return nil, err
	} else {
		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(3 * time.Minute)
		return tc, nil
	}
}

// Server extends the http.Server object with a MaxConnections property and
// LimitListenAndServe method that together limit the number of simultaneous
// connections the server will accept.
type Server struct {
	*http.Server
	MaxConnections int
}

// NewServer creates and initializes a new Server instance.
func NewServer(cf string, handler http.Handler) (*Server, error) {

	this := &Server{MaxConnections: 50} // Sane default

	if err := utils.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	this.Handler = handler
	this.ReadTimeout *= time.Second
	this.WriteTimeout *= time.Second

	return this, nil
}

// SetMaxConnections sets the maximum number of simultaneous connections
// the server will accept.
func (this *Server) SetMaxConnections(maxConnections int) {
	this.MaxConnections = maxConnections
}

// LimitListenAndServe listens on the TCP network address srv.Addr, then
// then calls Serve to handle requests on incoming connections. Accepted
// connections are configured to enable TCP keep-alives. It is identical
// to ListenAndServe except that it uses netutil.LimitListener to limit
// the number of simultaneous connections.
func (this *Server) LimitListenAndServe() error {

	addr := this.Addr
	if addr == `` {
		addr = `:http`
	}

	if ln, err := net.Listen(`tcp`, addr); err != nil {
		return err
	} else {
		lln := netutil.LimitListener(ln, this.MaxConnections)
		return this.Serve(tcpKeepAliveListener{lln.(*net.TCPListener)})
	}
}


// String provides identifying information about the server.
func (this *Server) String() (string) {
	return fmt.Sprintf(`http daemon listening on %q`, this.Addr)
}
