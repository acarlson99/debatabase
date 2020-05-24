package main

import (
	"fmt"
	"log"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	fmt.Println(router)

	db, err := DBConnect("root", "password", "radancomDB", "disable")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SHOW TABLES;")
	if err != nil {
		log.Fatal(err)
	}
	printRows(rows)
	rows.Close()

	db.Init()
	// db.populate()

	rows, err = db.Query("SHOW TABLES;")
	if err != nil {
		log.Fatal(err)
	}
	printRows(rows)
	rows.Close()

	rows, err = db.Query("SELECT * FROM articles")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var a string
		var b string
		// var c sql.NullString
		err := rows.Scan(&a, &b)
		if err != nil {
			fmt.Println("ERROR", err)
		}
		fmt.Println(a, b)
	}
	rows.Close()
}
