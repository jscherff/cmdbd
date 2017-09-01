package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"fmt"
	"io"

	"github.com/jscherff/gocmdb/webapi"
	"github.com/gorilla/mux"
)

const sizeLimit int64 = 1048576

func Serial(w http.ResponseWriter, r *http.Request) {

	// Need object type because different types of devices may 
	// have different formats and series of serial numbers.

	vars := mux.Vars(r)
	objectType := vars["objectType"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, sizeLimit))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	device := new(webapi.Device)

	if err := json.Unmarshal(body, &device); err != nil {

		w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")
		w.WriteHeader(http.StatusUnprocessableEntity)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	if len(device.SerialNum) == 0 {

		device.SerialNum = fmt.Sprintf("24F%04x", 1)	//TODO: generate actual serial number
		fmt.Println(device, objectType)			//TODO: remove

		w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(device); err != nil {
			panic(err)
		}

	} else {

		w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")
		w.WriteHeader(http.StatusNoContent)
	}
}

func Checkin(w http.ResponseWriter, r *http.Request) {

	// Need object type in order to instantiate the correct
	// object from the 'gocmdb' package.

	vars := mux.Vars(r)
	objectType := vars["objectType"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, sizeLimit))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	device := new(webapi.Device)

	if err := json.Unmarshal(body, &device); err != nil {

		w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")
		w.WriteHeader(http.StatusUnprocessableEntity)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	fmt.Println(device, objectType)	//TODO: record checkin to database

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")
	w.WriteHeader(http.StatusAccepted)
}

func Audit(w http.ResponseWriter, r *http.Request) {

	// Need only serial number, not object type, because method
	// will only log changes in the form {date, name, old, new}
	// associated with a device serial number. Serial number can
	// be matched to registration record.

	vars := mux.Vars(r)
	serialNum := vars["serialNum"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, sizeLimit))

	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	changes := new(webapi.Changes)

	if err := json.Unmarshal(body, &changes); err != nil {

		w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")
		w.WriteHeader(http.StatusUnprocessableEntity)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	fmt.Println(changes, serialNum)	//TODO: record changes to database

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")
	w.WriteHeader(http.StatusAccepted)
}
