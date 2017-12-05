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

package legacy

import (
	`crypto/rsa`
	`fmt`
	`io/ioutil`
	`os`
	`path/filepath`
	`time`

	jwt `github.com/dgrijalva/jwt-go`

	`github.com/jscherff/cmdbd/server`
	`github.com/jscherff/cmdbd/service`
	`github.com/jscherff/cmdbd/utils`
	`github.com/jscherff/gox/log`
)

var (
	// Program name and version.

	program = filepath.Base(os.Args[0])
	version = `undefined`

	// Encryption and signing keys.

	privateKey *rsa.PrivateKey
	publicKey *rsa.PublicKey

	// Configuration aliases.

	conf *Config

	db *Database
	dq *Queries
	ws *server.Server
)

var (
	authTokenSvc	service.AuthTokenService
	authCookieSvc	service.AuthCookieService
	serialNumSvc	service.SerialNumService
	middleWare	server.MiddleWare
	sysLog		log.MLogger
	accLog		log.MLogger
	errLog		log.MLogger
)

// Config contains infomation about the server process and log writers.
type Config struct {

	AuthMaxAge	time.Duration

	SerialFmt	map[string]string
	Configs		map[string]string
	KeyFiles	map[string]string

	Database	*Database
	Queries		*Queries
	Syslog		*server.Syslog
	Logger		*server.Logger
	Router		*server.Router
	MetaUsb		*MetaUsb
	Server		*server.Server
}

// NewConfig creates a new Config object and reads its config
// from the provided JSON configuration file.
func NewConfig(cf string, console, refresh bool) (this *Config, err error) {

	// Load the base config needed to load remaining configs.

	this = &Config{}

	if err := utils.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	// Prepend the base config directory to other config filenames.

	for key, fn := range this.Configs {
		this.Configs[key] = filepath.Join(filepath.Dir(cf), fn)
	}

	for key, fn := range this.KeyFiles{
		this.KeyFiles[key] = filepath.Join(filepath.Dir(cf), fn)
	}

	// Initialize services.

	if ts, err := service.NewAuthTokenService(this.KeyFiles, this.AuthMaxAge); err != nil {
		return nil, err
	} else {
		authTokenSvc = ts
	}

	if cs, err := service.NewAuthCookieService(this.AuthMaxAge); err != nil {
		return nil, err
	} else {
		authCookieSvc = cs
	}

	if ss, err := service.NewSerialNumService(this.SerialFmt); err != nil {
		return nil, err
	} else {
		serialNumSvc = ss
	}

	// Create and initialize MiddleWare object.

	if mw, err := server.NewMiddleWare(authTokenSvc, authCookieSvc); err != nil {
		return nil, err
	} else {
		middleWare = mw
	}

	// Create and initialize Database object.

	if database, err := NewDatabase(this.Configs[`Database`]); err != nil {
		return nil, err
	} else {
		this.Database = database
	}

	db = this.Database

	// Create and initialize Queries object.

	if queries, err := NewQueries(this.Configs[`Queries`], db); err != nil {
		return nil, err
	} else {
		this.Queries = queries
	}

	dq = this.Queries

	// Create and initialize Syslog object.

	if syslog, err := server.NewSyslog(this.Configs[`Syslog`]); err != nil {
		return nil, err
	} else {
		this.Syslog = syslog
	}

	// Create and initialize Logger object.

	if logger, err := server.NewLogger(this.Configs[`Logger`], console, this.Syslog); err != nil {
		return nil, err
	} else {
		this.Logger = logger
	}

	// Ensure required loggers are present and create aliases.

	var ok bool

	if sysLog, ok = this.Logger.Logs[`system`]; !ok {
		return nil, fmt.Errorf(`missing "system" log config`)
	}
	if accLog, ok = this.Logger.Logs[`access`]; !ok {
		return nil, fmt.Errorf(`missing "access" log config`)
	}
	if errLog, ok = this.Logger.Logs[`error`]; !ok {
		return nil, fmt.Errorf(`missing "error" log config`)
	}

	// Create and initialize Router object.

	if router, err := server.NewRouter(this.Configs[`Router`], middleWare, accLog, errLog); err != nil {
		return nil, err
	} else {
		this.Router = router.
			AddRoutes(usbCiRoutes).
			AddRoutes(usbMetaRoutes).
			AddRoutes(cmdbAuthRoutes)
	}

	// Create and initialize MetaUsb object.

	if metausb, err := NewMetaUsb(this.Configs[`MetaUsb`], refresh); err != nil {
		return nil, err
	} else {
		this.MetaUsb = metausb
	}

	// Create and initialize Server object.

	if server, err := server.NewServer(this.Configs[`Server`]); err != nil {
		return nil, err
	} else {
		server.Handler = this.Router
		this.Server = server
	}

	ws = this.Server

	// Read and store RSA private key.

	if pemKey, err := ioutil.ReadFile(this.KeyFiles[`PrivateRSA`]); err != nil {
		return nil, err
	} else if rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemKey); err != nil {
		return nil, err
	} else {
		privateKey = rsaKey
	}

	// Read and store RSA public key.

	if pemKey, err := ioutil.ReadFile(this.KeyFiles[`PublicRSA`]); err != nil {
		return nil, err
	} else if rsaKey, err := jwt.ParseRSAPublicKeyFromPEM(pemKey); err != nil {
		return nil, err
	} else {
		publicKey = rsaKey
	}

	conf = this
	return this, nil
}

// displayVersion displays the program version.
func DisplayVersion() {
	fmt.Fprintf(os.Stderr, "%s version %s\n", program, version)
}
