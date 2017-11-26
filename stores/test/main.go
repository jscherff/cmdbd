package main

import (
	`fmt`
	`github.com/jscherff/cmdbd/stores`
)

func main() {

	//ds, err := stores.NewMySqlDataStore(`mysql.json`)
	var ds stores.DataStore

	factory, err := stores.Factory(`mysql`)

	if err != nil {
		panic(err)
	}

	ds, err = factory(`mysql.json`)

	if err != nil {
		panic(err)
	}

	fmt.Println(ds.Version())
	fmt.Println(ds.Tables())
	fmt.Println(ds.Columns(`usbci_checkins`))

	if err := ds.Prepare(`queries.json`); err != nil {
		panic(err)
	}

	type User struct {
		Username string `db:"username"`
		Password string `db:"password"`
	}

	u := &User{Username: `clubpc`}

	if err := ds.Get(`SelectPassword`, u, u); err != nil {
		fmt.Println(err)
	}

	fmt.Println(u.Username, u.Password)



	//queries := make(map[string]Query)

	//ds.Prepare
}
