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

package service

import (
	`path/filepath`
	`time`

	`github.com/jscherff/cmdbd/common`
)

// Constants.
const (
	priKeyName	string = `PriKey`
	pubKeyName	string = `PubKey`
)

// Shared configurations and services.
type Config struct {
	AuthMaxAge	time.Duration
	AuthTokenSvc	AuthTokenService
	AuthCookieSvc	AuthCookieService
	SerialNumSvc	SerialNumService
	ConfigFile	map[string]string
	CryptoFile	map[string]string
	SerialFormat	map[string]string
}

// NewConfig creates a new Config object and reads its config
// from the provided JSON configuration file.
func NewConfig(cf string) (*Config, error) {

	// Load the base config needed to load remaining configs.

	this := &Config{}

	if err := common.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	// Set the maximum age of auth cookies and tokens.

	this.AuthMaxAge *= time.Minute

	// Prepend the master config directory to other filenames.

	for key, fn := range this.ConfigFile {
		this.ConfigFile[key] = filepath.Join(filepath.Dir(cf), fn)
	}

	for key, fn := range this.CryptoFile {
		this.CryptoFile[key] = filepath.Join(filepath.Dir(cf), fn)
	}

	// Initialize services.

	if ts, err := NewAuthTokenService(this); err != nil {
		return nil, err
	} else {
		this.AuthTokenSvc = ts
	}

	if cs, err := NewAuthCookieService(this); err != nil {
		return nil, err
	} else {
		this.AuthCookieSvc = cs
	}

	if ss, err := NewSerialNumService(this); err != nil {
		return nil, err
	} else {
		this.SerialNumSvc = ss
	}

	return this, nil
}
