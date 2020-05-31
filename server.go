package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Article struct {
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Tags []string `json:"tags"`
}

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CheckEnvVars checks environment variables to make sure they are set
func CheckEnvVars() {
	vars := []string{"APP_ENV", "MYSQL_USER", "MYSQL_PASSWORD", "MYSQL_DBNAME",
		"HOST_ADDRESS", "HOST_PORT"}
	for _, v := range vars {
		if len(os.Getenv(v)) == 0 {
			fmt.Println("WARNING: environment variable `" + v + "` not set")
		}
	}
}

func main() {
	if os.Getenv("APP_ENV") == "production" {
		fmt.Println("Running in production mode")
	} else {
		fmt.Println("Running in development mode")
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	CheckEnvVars()

	uname := os.Getenv("MYSQL_USER")
	passwd := os.Getenv("MYSQL_PASSWORD")
	dbname := os.Getenv("MYSQL_DBNAME")
	db, err := DBConnect(uname, passwd, dbname)
	if err != nil {
		log.Fatal(err)
	}
	db.Init()

	router := CreateRouter(db)
	addr := os.Getenv("HOST_ADDRESS")
	port := os.Getenv("HOST_PORT")
	fmt.Println("Listening and serving `" + addr + ":" + port + "`...")
	http.ListenAndServe(addr+":"+port, router)
}
