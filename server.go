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
	// UNameMinLen is min length of username
	UNameMinLen = 3
	// UNameMaxLen is max length of username
	UNameMaxLen = 30
	// UPasswdMinLen is min length of password
	UPasswdMinLen = 5
	// UPasswdMaxLen is max length of password
	UPasswdMaxLen = 50
)

var (
	hostPort string
	hostAddr string
	db       *DB
)

// DBArticle is a representation of an article from MySQL DB
type DBArticle struct {
	ID          int64  `json:"id,omitempty" example:"1"`
	Name        string `json:"name" maximum:"512" example:"google"`
	URL         string `json:"url" maximum:"512" example:"google.com"`
	Description string `json:"description" maximum:"1024" example:"a popular search engine"`
	// List of tag names
	Tags []string `json:"tags" example:"engine,search,browser"`
	// This is a list of filenames to be queried via other endpoint
	Images []string `json:"images" maxItems:"4" example:"a.png, evidence1.png, metal-beams.jpg"` // NOTE: limit of 4
}

// Image is an image format and Base64 representation of image
type Image struct {
	// Base64 encoded image data
	Data string `json:"data" format:"base64" example:"dGhpcyBpcyBhbiBpbWFnZQo="`
	// Image format (PNG,JPG,etc.)
	Format   string `json:"format" example:"png"`
	Filename string `swaggerignore:"true"`
}

// UploadArticle is a representation of an article sent from frontend to be uploaded to MySQL DB
type UploadArticle struct {
	Name        string   `json:"name" maximum:"512" example:"google"`
	URL         string   `json:"url" maximum:"512" example:"google.com"`
	Description string   `json:"description" maximum:"1024" example:"a popular search engine"`
	Tags        []string `json:"tags" example:"engine,search,browser"`
	Images      []Image  `json:"images" maxItems:"4"`
}

// DBTag is a representation of a tag from MySQL DB
type DBTag struct {
	ID          int64  `json:"id,omitempty" example:"1"`
	Name        string `json:"name" maximum:"16" example:"engine"`
	Description string `json:"description" maximum:"256" example:"a machine designed to convert one form of energy into mechanical energy"`
	// Name        string `json:"name"`
	// Description string `json:"description"`
}

// UploadTag is a representation of a tag sent from frontend to be uploaded to MySQL DB
type UploadTag struct {
	Name        string `json:"name" maximum:"16" example:"engine"`
	Description string `json:"description" maximum:"256" example:"a machine designed to convert one form of energy into mechanical energy"`
}

// type User struct {
// 	ID     int64  `json:"id,omitempty"` // NOTE: only updated when removing from DB
// 	Name   string `json:"name" minLength:"3" maxLength:"30"`
// 	Passwd string `json:"password" minLength:"5" maxLength:"50"`
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
// // @securityDefinitions.apikey Bearer
// // @in header
// // @name Authorization
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
		fmt.Println("FILES_TO_SERVE:", serveLocation)
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
