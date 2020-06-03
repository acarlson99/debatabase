package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var (
	nonAlphanumRE = regexp.MustCompile("[^a-zA-Z0-9]+")
)

type ArticleMsg struct {
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Tags []string `json:"tags"`
}

func generateSearchHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tags := mux.Vars(r)["tags"]
		limit, _ := strconv.Atoi(mux.Vars(r)["limit"])
		offset, _ := strconv.Atoi(mux.Vars(r)["offset"])
		lookslike := mux.Vars(r)["lookslike"]

		sp := []string{}
		if len(tags) > 0 {
			sp = strings.Split(tags, ",")
		}

		var articles []Article
		var err error

		articles, err = db.ArticlesWithTagsSearch(sp, lookslike, limit, offset)
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
		if err != nil || len(a.Name) == 0 || len(a.TagNames) == 0 || len(a.URL) == 0 {
			if err != nil {
				log.Println(err)
			}
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
	router := mux.NewRouter().StrictSlash(true)

	// search
	router.HandleFunc("/api/search/tags/{tags}/{limit}/{offset}/{lookslike}", generateSearchHandler(db))
	router.HandleFunc("/api/search/tags/{tags}/{limit}/{offset}", generateSearchHandler(db))
	router.HandleFunc("/api/search/tags/{tags}/{limit}/", generateSearchHandler(db))
	router.HandleFunc("/api/search/tags/{tags}/", generateSearchHandler(db))
	router.HandleFunc("/api/search/{limit}/{offset}/{lookslike}", generateSearchHandler(db))
	router.HandleFunc("/api/search/{limit}/{offset}", generateSearchHandler(db))
	router.HandleFunc("/api/search/{limit}/", generateSearchHandler(db))
	router.HandleFunc("/api/search/", generateSearchHandler(db))

	// upload DB
	router.HandleFunc("/api/upload/article", generateArticleHandler(db)).Methods("POST")
	router.HandleFunc("/api/upload/tag", generateTagHandler(db)).Methods("POST")
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/index.html")
	})
	return router
}
