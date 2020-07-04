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
		internalError("querying tags", w, err)
		return
	}

	resp, err := json.Marshal(articles)
	if err != nil {
		internalError("marshalling response", w, err)
		return
	}
	w.Write(resp)
}

// @Summary Search articles by name
// @Param id path integer false "Filter by ID"
// @Produce json
// @Success 200 {object} main.Article "All matching articles"
// @Failure 400 "Bad request"
// @Failure 404 "Article not found"
// @Failure 500 {string} string "Internal error"
// @Router /api/search/article/{id} [GET]
func searchArticleID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(400)
		return
	}
	articles, err := db.ArticleByID(int64(id))
	if err != nil {
		internalError("querying tags", w, err)
		return
	} else if articles == nil {
		w.WriteHeader(404)
		return
	}

	resp, err := json.Marshal(articles)
	if err != nil {
		internalError("marshalling response", w, err)
		return
	}
	w.Write(resp)
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
		internalError("querying tags", w, err)
		return
	}

	resp, err := json.Marshal(tags)
	if err != nil {
		internalError("marshalling response", w, err)
		return
	}
	w.Write(resp)
}

// @Summary Search tags by name
// @Param id path integer false "Filter by ID"
// @Produce json
// @Success 200 {object} main.Tag "All matching tags"
// @Failure 400 "Bad request"
// @Failure 404 "Tag not found"
// @Failure 500 {string} string "Internal error"
// @Router /api/search/tag/{id} [GET]
func searchTagID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(400)
		return
	}
	tags, err := db.TagByID(int64(id))
	if err != nil {
		internalError("querying tags", w, err)
		return
	} else if tags == nil {
		w.WriteHeader(404)
		return
	}

	resp, err := json.Marshal(tags)
	if err != nil {
		internalError("marshalling response", w, err)
		return
	}
	w.Write(resp)
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
		internalError("reading body", w, err)
		return
	}
	article := Article{}
	err = json.Unmarshal(body, &article)
	if err != nil || len(article.Name) == 0 {
		if err != nil {
			log.Println("Error unmarshalling data:", err)
		}
		w.WriteHeader(400)
		return
	}
	fmt.Printf("%+v\n", article)

	err = r.Body.Close()
	if err != nil {
		log.Println("Error closing http.Request body:", err)
	}

	// check if all tags exist
	if _, e := db.TagNamesExist(article.Tags...); !e {
		w.WriteHeader(422)
		return
	}

	_, err = db.InsertArticle(article)
	if err != nil {
		internalError("inserting article", w, err)
		return
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
		internalError("reading body", w, err)
		return
	}
	tag := Tag{}
	err = json.Unmarshal(body, &tag)
	if err != nil || len(tag.Name) == 0 {
		w.WriteHeader(400)
		return
	}
	fmt.Printf("%+v\n", tag)

	// check duplicates
	if _, r := db.TagNameExists(tag.Name); r {
		w.WriteHeader(403)
		log.Println("Not inserting tag. Already exists")
		return
	}

	err = r.Body.Close()
	if err != nil {
		log.Println("Error closing http.Request body:", err)
	}

	_, err = db.InsertTag(tag)
	if err != nil {
		internalError("inserting article", w, err)
		return
	}
}

// @Summary Modify Article
// @Accept  json
// @Param id path integer true "ID of article to modify"
// @Param article body main.DocArticle true "Updated article data"
// @Success 200 "Ok"
// @Failure 400 "Bad request"
// @Failure 404 "Article does not exist"
// @Failure 422 "Invalid tag(s)"
// @Failure 500 {string} string "Internal error"
// @Router /api/edit/article/{id} [POST]
func editArticle(w http.ResponseWriter, r *http.Request) {
	article := Article{}
	s, err := ioutil.ReadAll(r.Body)
	if err != nil {
		internalError("reading body", w, err)
		return
	}
	r.Body.Close()
	err = json.Unmarshal(s, &article)
	id2, err2 := strconv.Atoi(mux.Vars(r)["id"])
	id := int64(id2)
	if err != nil || err2 != nil || len(article.Name) == 0 {
		w.WriteHeader(400)
		return
	}
	// check if article exists
	res, err := db.ArticleByID(id)
	if err != nil {
		internalError("querying DB", w, err)
		return
	} else if res == nil {
		w.WriteHeader(404)
		return
	}
	// check if tags exist
	tagIDs, exists := db.TagNamesExist(article.Tags...)
	if !exists {
		w.WriteHeader(422)
		return
	}
	// update
	err = db.UpdateArticle(id, article)
	if err != nil {
		internalError("updating article", w, err)
		return
	}
	// update tags
	err = db.RemoveArticleTags(id)
	if err != nil {
		internalError("removing tags", w, err)
		return
	}
	err = db.InsertArticleTags(id, tagIDs)
	if err != nil {
		internalError("updating tags", w, err)
		return
	}
}

// @Summary Modify Tag
// @Accept  json
// @Param id path integer true "ID of tag to modify"
// @Param tag body main.DocTag true "Updated tag data"
// @Success 200 "Ok"
// @Failure 400 "Bad request"
// @Failure 404 "Tag does not exist"
// @Failure 500 {string} string "Internal error"
// @Router /api/edit/tag/{id} [POST]
func editTag(w http.ResponseWriter, r *http.Request) {
	tag := Tag{}
	s, err := ioutil.ReadAll(r.Body)
	if err != nil {
		internalError("reading body", w, err)
		return
	}
	r.Body.Close()
	err = json.Unmarshal(s, &tag)
	id2, err2 := strconv.Atoi(mux.Vars(r)["id"])
	id := int64(id2)
	if err != nil || err2 != nil || len(tag.Name) == 0 {
		w.WriteHeader(400)
		return
	}
	// check if exists
	res, err := db.TagByID(id)
	if err != nil {
		internalError("querying DB", w, err)
		return
	} else if res == nil {
		w.WriteHeader(404)
		return
	}
	// update
	err = db.UpdateTag(id, tag)
	if err != nil {
		internalError("updating tag", w, err)
		return
	}
}

// @Summary Delete Article
// @Accept  json
// @Param id path integer true "ID of article to modify"
// @Success 200 "Ok"
// @Failure 400 "Bad request"
// @Failure 404 "Tag does not exist"
// @Failure 500 {string} string "Internal error"
// @Router /api/del/article/{id} [GET]
func deleteArticle(w http.ResponseWriter, r *http.Request) {
	id2, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(400)
		return
	}
	id := int64(id2)
	res, err := db.ArticleByID(id)
	if err != nil {
		internalError("querying DB", w, err)
		return
	} else if res == nil {
		w.WriteHeader(404)
		return
	}
	err = db.RemoveArticle(id)
	if err != nil {
		internalError("querying DB", w, err)
		return
	}
	err = db.RemoveArticleTags(id)
	if err != nil {
		internalError("querying DB", w, err)
		return
	}
}

// @Summary Delete Tag
// @Accept  json
// @Param id path integer true "ID of tag to modify"
// @Success 200 "Ok"
// @Failure 400 "Bad request"
// @Failure 404 "Tag does not exist"
// @Failure 500 {string} string "Internal error"
// @Router /api/del/tag/{id} [GET]
func deleteTag(w http.ResponseWriter, r *http.Request) {
	id2, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(400)
		return
	}
	id := int64(id2)
	res, err := db.TagByID(id)
	if err != nil {
		internalError("querying DB", w, err)
		return
	} else if res == nil {
		w.WriteHeader(404)
		return
	}
	err = db.RemoveTag(id)
	if err != nil {
		internalError("querying DB", w, err)
		return
	}
	err = db.RemoveTagsFromArticles(id)
	if err != nil {
		internalError("querying DB", w, err)
		return
	}
}

// internalError writes a 500 response to a ResponseWriter and logs an error
func internalError(logMsg string, w http.ResponseWriter, err error) {
	log.Println("Error", logMsg+": ", err)
	w.WriteHeader(500)
	w.Write([]byte(fmt.Sprintf("%v", err)))
}

// @Summary Create User
// @Accept  json
// @Param user body main.DocUser true "User data"
// @Success 200 "Ok"
// @Failure 400 "Bad request"
// @Failure 403 "Duplicate"
// @Failure 500 {string} string "Internal error"
// @Router /api/user/create [POST]
func userCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		internalError("reading body", w, err)
		return
	}
	user := User{}
	err = json.Unmarshal(body, &user)
	if err != nil || len(user.Name) < UNameMinLen || len(user.Name) > UNameMaxLen || len(user.Passwd) < UPasswdMinLen || len(user.Passwd) > UNameMaxLen {
		if err != nil {
			log.Println("Error unmarshalling data:", err)
		}
		w.WriteHeader(400)
		return
	}

	// no duplicates
	u, err := db.UserByName(user.Name)
	if u != nil {
		// user already exists
		// forbidden
		w.WriteHeader(403)
		return
	} else if err != nil {
		internalError("querying users", w, err)
		return
	}

	_, err = db.InsertUser(user)
	if err != nil {
		internalError("querying users", w, err)
		return
	}
}

// @Summary Log in as User
// @Accept  json
// @Param user body main.DocUser true "User data"
// @Success 200 {object} main.User "JWT token"
// @Failure 400 "Bad request"
// @Failure 403 "Invalid credentials"
// @Failure 500 {string} string "Internal error"
// @Router /api/user/auth [POST]
func userAuthHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		internalError("reading body", w, err)
		return
	}
	user := User{}
	err = json.Unmarshal(body, &user)
	if err != nil || len(user.Name) < UNameMinLen || len(user.Name) > UNameMaxLen || len(user.Passwd) < UPasswdMinLen || len(user.Passwd) > UNameMaxLen {
		if err != nil {
			log.Println("Error unmarshalling data:", err)
		}
		w.WriteHeader(400)
		return
	}

	u, err := db.UserByName(user.Name)
	if err != nil {
		internalError("querying users", w, err)
		return
	} else if u == nil || user.Name != u.Name || user.Passwd != u.Passwd {
		// user doesn't exist.  Cannot log in
		// forbidden
		w.WriteHeader(403)
		return
	}

	// TODO: create JWT token and write to connection
	resp, err := json.Marshal(*u)
	if err != nil {
		internalError("querying users", w, err)
		return
	}
	w.Write(resp)
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
		httpSwagger.DocExpansion("list"),
	))

	// search
	r.HandleFunc("/api/search/article/{id}", searchArticleID)
	r.HandleFunc("/api/search/article", searchArticle)
	r.HandleFunc("/api/search/tag/{id}", searchTagID)
	r.HandleFunc("/api/search/tag", searchTag)
	// upload
	r.HandleFunc("/api/upload/article/csv", uploadCSVArticle).Methods("POST") // create new article
	r.HandleFunc("/api/upload/article", uploadArticle).Methods("POST")        // create new article
	r.HandleFunc("/api/upload/tag/csv", uploadCSVTag).Methods("POST")         // create new tag
	r.HandleFunc("/api/upload/tag", uploadTag).Methods("POST")                // create new tag
	// edit
	r.HandleFunc("/api/edit/article/{id}", editArticle).Methods("POST") // modify article by ID
	r.HandleFunc("/api/edit/tag/{id}", editTag).Methods("POST")         // modify tag by ID
	// delete
	r.HandleFunc("/api/del/article/{id}", deleteArticle)
	r.HandleFunc("/api/del/tag/{id}", deleteTag)
	// user
	// TODO: add users
	r.HandleFunc("/api/user/create", userCreateHandler) // creates user
	r.HandleFunc("/api/user/auth", userAuthHandler)     // sends Json Web Token to client if uname/passwd match DB

	// serve
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/build/")))
	return r
}
