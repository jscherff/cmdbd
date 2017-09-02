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
	"path/filepath"
	"encoding/json"
	"log"
	"os"
)

type WsConfig struct {
	LogDir string `json:"log_dir"`
}

const wsConfigFile string = "wsconfig.json"

var wsConfig *WsConfig

func init() {

	var err error

	fn := filepath.Join(filepath.Dir(os.Args[0]), wsConfigFile)
	fh, err := os.Open(fn)

	if err != nil {
		log.Fatalf("%v", err)
	}

	defer fh.Close()
	jd := json.NewDecoder(fh)

	if err = jd.Decode(&wsConfig); err != nil {
		log.Fatalf("%v", err)
	}
}
