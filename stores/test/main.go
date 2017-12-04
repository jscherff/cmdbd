package main

import (
	`fmt`
	`github.com/jscherff/cmdbd/stores`
)

func main() {

	ds, err := stores.NewMysqlDataStore(`mysql.json`)

	if err != nil {
		panic(err)
	}

	fmt.Println(ds.Version())

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
}
