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

// Source: https://thenewstack.io/make-a-restful-json-api-go/

import (
	"net/http"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
}

func main() {

	defer db.Close()
	defer auditInsertStmt.Close()
	defer checkinInsertStmt.Close()
	defer serialInsertStmt.Close()
	defer serialUpdateStmt.Close()

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
