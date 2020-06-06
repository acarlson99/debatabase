package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	// (*w).Header().Set("Access-Control-Allow-Methods", "HEAD, GET, POST, PUT, PATCH, DELETE, OPTIONS")
	// (*w).Header().Set("Access-Control-Allow-Headers", "X-API-KEY, Origin, X-Requested-With, Content-Type, Accept, Access-Control-Request-Method,Access-Control-Request-Headers, Authorization")
	// (*w).Header().Set("Content-Type", "application/json")
	// $method = $_SERVER["REQUEST_METHOD"];
	// if ($method == "OPTIONS") {
	// header("Access-Control-Allow-Origin: *");
	// header("Access-Control-Allow-Headers: X-API-KEY, Origin, X-Requested-With, Content-Type, Accept, Access-Control-Request-Method,Access-Control-Request-Headers, Authorization");
	// header("HTTP/1.1 200 OK");
	// die();
	// }
}

func generateArticleSearchHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		parts := make(map[string]string)
		for k, v := range r.URL.Query() {
			parts[k] = v[0]
		}
		tags := parts["tags"]
		limit, _ := strconv.Atoi(parts["limit"])
		offset, _ := strconv.Atoi(parts["offset"])
		lookslike := parts["lookslike"]

		sp := []string{}
		if len(tags) > 0 {
			sp = strings.Split(tags, ",")
		}

		articles, err := db.ArticlesWithTagsSearch(sp, lookslike, limit, offset)
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
			log.Println("Error marshalling response:", err)
		} else {
			w.Write(resp)
		}
	}
}

func generateTagSearchHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		parts := make(map[string]string)
		for k, v := range r.URL.Query() {
			parts[k] = v[0]
		}
		tagStr := parts["tags"]
		limit, _ := strconv.Atoi(parts["limit"])
		offset, _ := strconv.Atoi(parts["offset"]) // NOTE: does nothing unless `limit` is specified
		lookslike := parts["lookslike"]

		sp := []string{}
		if len(tagStr) > 0 {
			sp = strings.Split(tagStr, ",")
		}

		tags, err := db.TagSearch(sp, lookslike, limit, offset)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			log.Println("Error querying tags:", err)
			return
		}

		resp, err := json.Marshal(tags)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("%v", err)))
			log.Println("Error marshalling response:", err)
		} else {
			w.Write(resp)
		}
	}
}

func generateArticleCSVHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		reader := csv.NewReader(r.Body)
		for {
			// name,url,description,tags
			fields, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Println("Error parsing CSV:", err)
			}
			if len(fields) != 4 {
				log.Println("BAD NUMNER OF FIELDS")
				continue
			}
			a := Article{
				Name:        fields[0],
				URL:         fields[1],
				Description: fields[2],
				Tags:        strings.Split(fields[3], ","),
			}
			_, err = db.InsertArticle(a)
			if err != nil {
				log.Println("Error inserting article:", err)
			}
		}
		err := r.Body.Close()
		if err != nil {
			log.Println("Error closing http.Request body:", err)
		}
	}
}

func generateTagCSVHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		reader := csv.NewReader(r.Body)
		for {
			// name,description
			fields, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Println("Error parsing CSV:", err)
			}
			if len(fields) != 2 {
				log.Println("BAD NUMNER OF FIELDS")
				continue
			}
			t := Tag{
				Name:        fields[0],
				Description: fields[1],
			}
			_, err = db.InsertTag(t)
			if err != nil {
				log.Println("Error inserting article:", err)
			}
		}
		err := r.Body.Close()
		if err != nil {
			log.Println("Error closing http.Request body:", err)
		}
	}
}

func generateArticleHandler(db *DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading body:", err)
			return
		}
		a := Article{}
		err = json.Unmarshal(body, &a)
		if err != nil || len(a.Name) == 0 {
			if err != nil {
				log.Println("Error unmarshalling data:", err)
			}
			w.WriteHeader(400)
			return
		}
		fmt.Printf("%+v\n", a)

		err = r.Body.Close()
		if err != nil {
			log.Println("Error closing http.Request body:", err)
		}

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
		enableCors(&w)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading body:", err)
			return
		}
		a := Tag{}
		err = json.Unmarshal(body, &a)
		if err != nil || len(a.Name) == 0 {
			w.WriteHeader(400)
			return
		}
		fmt.Printf("%+v\n", a)

		err = r.Body.Close()
		if err != nil {
			log.Println("Error closing http.Request body:", err)
		}

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
	router.HandleFunc("/api/search/article", generateArticleSearchHandler(db))
	router.HandleFunc("/api/search/tag", generateTagSearchHandler(db))

	// upload
	router.HandleFunc("/api/upload/article/csv", generateArticleCSVHandler(db)).Methods("POST") // create new article
	router.HandleFunc("/api/upload/article", generateArticleHandler(db)).Methods("POST")        // create new article
	router.HandleFunc("/api/upload/tag/csv", generateTagCSVHandler(db)).Methods("POST")         // create new tag
	router.HandleFunc("/api/upload/tag", generateTagHandler(db)).Methods("POST")                // create new tag

	// serve
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/index.html")
	})
	return router
}
