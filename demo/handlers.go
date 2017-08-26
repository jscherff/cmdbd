package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"fmt"
	"io"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo Show:", todoId)
}

func TodoCreate(w http.ResponseWriter, r *http.Request) {

	var todo Todo

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &todo); err != nil {

		w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")
		w.WriteHeader(422) // unprocessable entity

		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t := RepoCreateTodo(todo)

	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}
