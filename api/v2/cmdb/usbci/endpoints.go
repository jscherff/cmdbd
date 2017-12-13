// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS).
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package usbci

import `github.com/gorilla/mux`

// NewRoutes returns a collection of REST APIv2 endpoints providing CMDB CI data.
func (this *handlers) AddRoutes(router *mux.Router) *mux.Router {

	subRouter := router.PathPrefix(`/v2/cmdb/usbci`).Subrouter()

	subRouter.NewRoute().
		Name(`USB CI Checkin Handler`).
		Methods(`POST`).
		Path(`/checkin/{host}/{vid}/{pid}`).
		HandlerFunc(this.Checkin)

	subRouter.NewRoute().
		Name(`USB CI Checkout Handler`).
		Methods(`GET`).
		Path(`/checkout/{host}/{vid}/{pid}/{sn}`).
		HandlerFunc(this.CheckOut)

	subRouter.NewRoute().
		Name(`USB CI NewSn Handler`).
		Methods(`POST`).
		Path(`/newsn/{host}/{vid}/{pid}`).
		HandlerFunc(this.NewSn)

	subRouter.NewRoute().
		Name(`USB CI Audit Handler`).
		Methods(`POST`).
		Path(`/audit/{host}/{vid}/{pid}/{sn}`).
		HandlerFunc(this.Audit)
}
