package main

import (
	`fmt`
	`github.com/jscherff/cmdbd/stores`
)

func main() {

	ds, err := stores.NewMySqlDataStore(`mysql.json`)

	if err != nil {
		fmt.Printf(`%v`, err)
	}

	fmt.Println(ds.Tables())
	fmt.Println(ds.Columns(`usbci_checkins`))

	//queries := make(map[string]Query)

	//ds.Prepare
}
