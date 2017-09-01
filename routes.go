package main

import "net/http"

type Route struct {
	Name		string
	Method		string
	Pattern		string
	HandlerFunc	http.HandlerFunc
}

type Routes []Route

var routes = Routes {

	Route {
		Name:		"Serial",
		Method:		"POST",
		Pattern:	"/serial/{objectType}",
		HandlerFunc:	Serial,
	},

	Route {
		Name:		"Checkin",
		Method:		"POST",
		Pattern:	"/checkin/{objectType}",
		HandlerFunc:	Checkin,
	},

	Route {
		Name:		"Audit",
		Method:		"POST",
		Pattern:	"/audit/{serialNum}",
		HandlerFunc:	Audit,
	},
}
