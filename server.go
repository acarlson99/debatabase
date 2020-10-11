package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

const (
	UNameMinLen   = 3
	UNameMaxLen   = 30
	UPasswdMinLen = 5
	UPasswdMaxLen = 50
)

var (
	hostPort string
	hostAddr string
	db       *DB
)

type Article struct {
	ID          int64    `json:"id,omitempty"` // NOTE: only updated when removing from DB
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type Tag struct {
	ID          int64  `json:"id,omitempty"` // NOTE: only updated when removing from DB
	Name        string `json:"name"`
	Description string `json:"description"`
}

// type User struct {
// 	ID     int64  `json:"id,omitempty"` // NOTE: only updated when removing from DB
// 	Name   string `json:"name"`
// 	Passwd string `json:"password"`
// }

// CheckEnvVars checks environment variables to make sure they are set
func CheckEnvVars() {
	vars := []string{"APP_ENV", "MYSQL_USER", "MYSQL_PASSWORD", "MYSQL_HOSTNAME",
		"MYSQL_DBNAME", "HOST_ADDRESS", "HOST_PORT", "FILES_TO_SERVE"}
	for _, v := range vars {
		if len(os.Getenv(v)) == 0 {
			fmt.Println("WARNING: environment variable `" + v + "` not set")
		}
	}
}

// @title DB
// @version 1.0
// @description Debatabase
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
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
	hostname := os.Getenv("MYSQL_HOSTNAME")
	var err error
	fmt.Println("Connecting to database...")
	db, err = DBConnect(uname, passwd, hostname, dbname)
	if err != nil {
		fmt.Println("Failed to connect to MySQL DB.  Is DB running?")
		fmt.Println(err)
		os.Exit(1)
	}
	db.Init()
	go DBMaintainConnection(uname, passwd, hostname, dbname, 15)

	serveLocation := os.Getenv("FILES_TO_SERVE")
	if len(serveLocation) == 0 {
		serveLocation = "./frontend/build/"
	}
	ok, err := regexp.MatchString("/(build)|(static)/?$", serveLocation)
	if err != nil {
		panic(err)
	}
	if !ok {
		fmt.Println("WARNING: `FILES_TO_SERVE` does not end in `/build/` or `/static/`")
		fmt.Println("Are you sure you would like to continue? [y/n]")
		input := ""
		fmt.Scanln(&input)
		if strings.ToLower(input) != "y" {
			os.Exit(1)
		}
	}
	r := CreateRouter(serveLocation)

	hostAddr = os.Getenv("HOST_ADDRESS")
	hostPort = os.Getenv("HOST_PORT")
	fmt.Println("Listening and serving `" + hostAddr + ":" + hostPort + "`...")
	go http.ListenAndServe(hostAddr+":"+hostPort, r)

	// handle interrupts and shutdown safely
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	sig := <-c

	fmt.Println("Received", sig, "signal.  Shutting down...")
	db.Close()
}
