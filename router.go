package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/acarlson99/praxis/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @Summary Search articles by name
// @Param tags query string false "Tag names" collectionFormat(csv)
// @Param limit query integer false "Maximum number of results"
// @Param offset query integer false "Results to skip.  Does nothing unless 'limit' is specified"
// @Param lookslike query string false "Filter for matching names/descriptions"
// @Param orderby query string false "Field by which to order results" Enums(id, name, description)
// @Produce json
// @Success 200 {array} main.Article "All matching articles"
// @Failure 500 {string} string "Internal error"
// @Router /api/search/article?articles=engine,train&limit=5&offset=5&lookslike=american&orderby=name [GET]
func searchArticle(w http.ResponseWriter, r *http.Request) {
	parts := make(map[string]string)
	for k, v := range r.URL.Query() {
		parts[k] = v[0]
	}
	tags := parts["tags"]
	limit, _ := strconv.Atoi(parts["limit"])
	offset, _ := strconv.Atoi(parts["offset"])
	lookslike := parts["lookslike"]
	orderby := parts["orderby"]
	rev := parts["reverse"] == "true"

	sp := []string{}
	if len(tags) > 0 {
		sp = strings.Split(tags, ",")
	}

	articles, err := db.ArticlesWithTagsSearch(sp, lookslike, orderby, rev, limit, offset)
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

// @Summary Search articles by name
// @Param id path integer false "Filter by ID"
// @Produce json
// @Success 200 {array} main.Article "All matching articles"
// @Failure 500 {string} string "Internal error"
// @Router /api/search/article/{id} [GET]
func searchArticleID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	articles, err := db.ArticleByID(id)
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

// @Summary Search tags by name
// @Param tags query string false "Tag names" collectionFormat(csv)
// @Param limit query integer false "Maximum number of results"
// @Param offset query integer false "Results to skip.  Does nothing unless 'limit' is specified"
// @Param lookslike query string false "Filter for matching names/descriptions"
// @Param orderby query string false "Field by which to order results" Enums(id, name, description)
// @Produce json
// @Success 200 {array} main.Tag "All matching tags"
// @Failure 500 {string} string "Internal error"
// @Router /api/search/tag?tags=engine,train&limit=5&offset=5&lookslike=american&orderby=name [GET]
func searchTag(w http.ResponseWriter, r *http.Request) {
	parts := make(map[string]string)
	for k, v := range r.URL.Query() {
		parts[k] = v[0]
	}
	tagStr := parts["tags"]
	limit, _ := strconv.Atoi(parts["limit"])
	offset, _ := strconv.Atoi(parts["offset"]) // NOTE: does nothing unless `limit` is specified
	lookslike := parts["lookslike"]
	orderby := parts["orderby"]
	rev := parts["reverse"] == "true"

	sp := []string{}
	if len(tagStr) > 0 {
		sp = strings.Split(tagStr, ",")
	}

	tags, err := db.TagSearch(sp, lookslike, orderby, rev, limit, offset)
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

// @Summary Search tags by name
// @Param id path integer false "Filter by ID"
// @Produce json
// @Success 200 {array} main.Tag "All matching tags"
// @Failure 500 {string} string "Internal error"
// @Router /api/search/tag/{id} [GET]
func searchTagID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	tags, err := db.TagByID(id)
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

func uploadCSVArticle(w http.ResponseWriter, r *http.Request) {
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

func uploadCSVTag(w http.ResponseWriter, r *http.Request) {
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

// @Summary Create Article
// @Accept json
// @Param tag body main.DocArticle true "Article data"
// @Success 200 "Ok"
// @Failure 400 "Bad request"
// @Failure 422 "Invalid tag(s)"
// @Failure 500 {string} string "Internal error"
// @Router /api/upload/article [POST]
func uploadArticle(w http.ResponseWriter, r *http.Request) {
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

	// check if all tags exist
	if !db.TagNamesExist(a.Tags...) {
		w.WriteHeader(422)
		return
	}

	_, err = db.InsertArticle(a)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		log.Println("Error inserting article:", err)
	}
}

// @Summary Create Tag
// @Accept  json
// @Param tag body main.DocTag true "Tag data"
// @Success 200 "Ok"
// @Failure 400 "Bad request"
// @Failure 403 "Duplicate tag"
// @Failure 500 {string} string "Internal error"
// @Router /api/upload/tag [POST]
func uploadTag(w http.ResponseWriter, r *http.Request) {
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

	// check duplicates
	if _, r := db.TagNameExists(a.Name); r {
		w.WriteHeader(403)
		log.Println("Not inserting tag. Already exists")
		return
	}

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

func generateAuthHandler(w http.ResponseWriter, r *http.Request) {
	values := mux.Vars(r)
	fmt.Println(values["uname"])
	fmt.Println(values["passwd"])
	fmt.Println(values)
}

func enableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
}

// CreateRouter returns a new mux.Router with appropriately registered paths
func CreateRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Use(enableCors)

	// swagger serve
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
		// httpSwagger.DeepLinking(true),
	))

	// search
	r.HandleFunc("/api/search/article/{id}", searchArticleID)
	r.HandleFunc("/api/search/article", searchArticle)
	r.HandleFunc("/api/search/tag/{id}", searchTagID)
	r.HandleFunc("/api/search/tag", searchTag)

	// upload
	// TODO: add `edit` feature for articles
	r.HandleFunc("/api/upload/article/csv", uploadCSVArticle).Methods("POST") // create new article
	r.HandleFunc("/api/upload/article", uploadArticle).Methods("POST")        // create new article
	r.HandleFunc("/api/upload/tag/csv", uploadCSVTag).Methods("POST")         // create new tag
	r.HandleFunc("/api/upload/tag", uploadTag).Methods("POST")                // create new tag

	// user
	// TODO: add users
	r.HandleFunc("/api/user/auth/{uname}/{passwd}", generateAuthHandler)   // sends Json Web Token to client if uname/passwd match DB
	r.HandleFunc("/api/user/create/{uname}/{passwd}", generateAuthHandler) // creates user

	// serve
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/build/")))
	return r
}
