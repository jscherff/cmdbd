package main

// Source: https://thenewstack.io/make-a-restful-json-api-go/

import (
	"net/http"
	"log"
)

func main() {

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
