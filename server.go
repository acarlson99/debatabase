package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
)

type Article struct {
	ID          int64    `json:"id"` // NOTE: only updated when removing from DB
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type Tag struct {
	ID          int64  `json:"id"` // NOTE: only updated when removing from DB
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
		fmt.Println("Failed to connect to MySQL DB.  Is DB running?")
		fmt.Println(err)
		os.Exit(1)
	}
	db.Init()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		sig := <-c

		fmt.Println("Received", sig, "signal.  Shutting down...")
		db.Close()
		os.Exit(0)
	}()

	router := CreateRouter(db)
	addr := os.Getenv("HOST_ADDRESS")
	port := os.Getenv("HOST_PORT")
	fmt.Println("Listening and serving `" + addr + ":" + port + "`...")
	http.ListenAndServe(addr+":"+port, router)
}
