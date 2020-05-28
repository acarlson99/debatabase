package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func generateQueryHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(mux.Vars(r)["query"])
	}
}

// CreateRouter returns a new mux.Router with appropriately registered paths
func CreateRouter(db *DB) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/query/{query}", generateQueryHandler(db))
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/index.html")
	})
	return router
}

// CheckEnvVars checks environment variables to make sure they are set
func CheckEnvVars() {
	vars := []string{"APP_ENV", "MYSQL_USER", "MYSQL_PASSWORD", "HOST_ADDRESS", "HOST_PORT"}
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
	db, err := DBConnect(uname, passwd, "radancomDB")
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
