package main

import (
	//`fmt`
	`log`
	`github.com/jmoiron/sqlx`
	_ `github.com/go-sql-driver/mysql`
)

func main() {

	db, err := sqlx.Open("mysql", "cmdbd:K2Cvg3NeyR@/gocmdb")
	if err != nil {
		log.Fatal(err)
	}

        if err = db.Ping(); err != nil {
                log.Fatal(err)
        }
}
