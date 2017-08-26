package main

import (
	"net/http"
)

type Route struct {
	Name		string
	Method		string
	Pattern		string
	HandlerFunc	http.HandlerFunc
}

type Routes []Route

var routes = Routes {
	Route {
		Name:		"Index",
		Method:		"GET",
		Pattern:	"/",
		HandlerFunc:	Index,
	},

	Route {
		Name:		"TodoIndex",
		Method:		"GET",
		Pattern:	"/todos",
		HandlerFunc:	TodoIndex,
	},

	Route {
		Name:		"TodoSHow",
		Method:		"GET",
		Pattern:	"/todos/{todoId}",
		HandlerFunc:	TodoShow,
	},

	Route {
		Name:		"TodoCreate",
		Method:		"POST",
		Pattern:	"/todos",
		HandlerFunc:	TodoCreate,
	},
}
