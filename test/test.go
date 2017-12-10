package main

import (
	`github.com/jscherff/cmdbd/model/cmdb`
	//`github.com/jscherff/cmdbd/model/cmdb/usbci`
	//`github.com/jscherff/cmdbd/model/cmdb/usbmeta`
	`github.com/jscherff/cmdbd/store`
)

func main() {

	if ds, err := store.New(`mysql`, `config/store/mysql.json`); err != nil {
		panic(err)
	} else if stmts, err := ds.Prepare(`config/model/queries.json`); err != nil {
		panic(err)
	} else {
		cmdb.Model.Init(stmts)
	}

	user := &cmdb.User{}
	user.Username = `clubpc`

	user.Read(user)

	println(user.Username, user.Password)

	if err := user.Verify(`test`); err != nil {
		println(err.Error())
	} else {
		println(`verified`)
	}

	//fmt.Println(cmdb.Model.Stmts().Statement(
}
