package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("APP_ENV") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	uname := os.Getenv("MYSQL_USER")
	passwd := os.Getenv("MYSQL_PASSWORD")

	db, err := DBConnect(uname, passwd, "radancomDB")
	if err != nil {
		log.Fatal(err)
	}

	db.Init()

	router := mux.NewRouter()
	addr := os.Getenv("HOST_ADDRESS")
	port := os.Getenv("HOST_PORT")

	log.Println("Listening and serving `" + addr + ":" + port + "`...")
	http.ListenAndServe(addr+":"+port, router)
}
