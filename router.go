package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

var (
	nonAlphanumRE = regexp.MustCompile("[^a-zA-Z0-9]+")
)

func generateQueryHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tags := mux.Vars(r)["tags"]

		// TODO: verify that tags are CSV

		sp := strings.Split(tags, ",")

		filtered := sp[:0]
		for _, t := range sp {
			if len(t) > 0 {
				filtered = append(filtered, nonAlphanumRE.ReplaceAllString(t, ""))
			}
		}

		articles, err := db.ArticlesWithTags(filtered, 0, 10)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			log.Println("Error querying tags:", err)
			return
		}

		resp, err := json.Marshal(articles)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			log.Println("Error marshaling response")
		} else {
			w.Write(resp)
		}

		// w.WriteHeader(501)
	}
}

func generateArticleHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		a := Article{}
		err = json.Unmarshal(body, &a)
		r.Body.Close()
		if err != nil || len(a.Name) == 0 || len(a.Tags) == 0 || len(a.URL) == 0 {
			w.WriteHeader(400)
			return
		}
		fmt.Printf("%+v\n", a)

		_, err = db.InsertArticle(a)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			log.Println("Error inserting article:", err)
		}
	}
}

func generateTagHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		a := Tag{}
		err = json.Unmarshal(body, &a)
		r.Body.Close()
		if err != nil || len(a.Name) == 0 || len(a.Description) == 0 {
			w.WriteHeader(400)
			return
		}
		fmt.Printf("%+v\n", a)

		_, err = db.InsertTag(a)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			log.Println("Error inserting article:", err)
		}
	}
}

// CreateRouter returns a new mux.Router with appropriately registered paths
func CreateRouter(db *DB) *mux.Router {
	router := mux.NewRouter()

	// search
	router.HandleFunc("/api/query/tags/{tags}", generateQueryHandler(db))

	// upload DB
	router.HandleFunc("/api/upload/article", generateArticleHandler(db)).Methods("POST")
	router.HandleFunc("/api/upload/tag", generateTagHandler(db)).Methods("POST")
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/index.html")
	})
	return router
}
