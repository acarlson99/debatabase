package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB internal database struct
type DB struct {
	*sql.DB
}

func makeConnStr(uname, password, hostname, dbname string) string {
	connStr := fmt.Sprintf("%s:%s@", uname, password)
	if len(hostname) > 0 {
		connStr += fmt.Sprintf("tcp(%s)", hostname)
	}
	connStr += "/"
	return connStr
}

// DBConnect creates connection to database (through hostname if it exists) with credentials
func DBConnect(uname, password, hostname, dbname string) (*DB, error) {
	connStr := makeConnStr(uname, password, hostname, dbname)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	// create && use database
	s := "CREATE DATABASE IF NOT EXISTS " + dbname + ";"
	_, err = db.Exec(s)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("USE " + dbname + ";")
	return &DB{db}, err
}

// DBMaintainConnection checks connection to DB every `period` seconds and attempts to recreate DB on failure
func DBMaintainConnection(uname, password, hostname, dbname string, period int) {
	t := time.NewTicker(time.Duration(period) * time.Second)

	for {
		_ = <-t.C
		_, err := db.Exec("SELECT ID FROM users LIMIT 0;")
		if err == nil {
			continue
		}
		log.Println("Error pinging DB:", err)
		newDB, err := DBConnect(uname, password, hostname, dbname)
		if err != nil {
			log.Println("Error reconnecting to DB:", err)
			if newDB != nil {
				newDB.Close()
			}
			continue
		}
		log.Println("Successfully reconnected to DB")
		db.Close()
		db = newDB
	}
}

// Init things
func (db *DB) Init() {
	// make sure `articles` exists
	fmt.Println("Initializing database...")
	if !db.tableExists("articles") {
		fmt.Println("DB creating table `articles`...")
		_, err := db.Exec("CREATE TABLE articles( ID INT AUTO_INCREMENT, Name VARCHAR(512) NOT NULL, URL VARCHAR(512), Description VARCHAR(1024), PRIMARY KEY (ID) );")
		if err != nil {
			log.Fatal(err)
		}
	}
	// make sure `tags` exists
	if !db.tableExists("tags") {
		fmt.Println("DB creating table `tags`...")
		_, err := db.Exec("CREATE TABLE tags( ID INT AUTO_INCREMENT, Name VARCHAR(16) UNIQUE, Description VARCHAR(256), PRIMARY KEY (ID, Name) );")
		if err != nil {
			log.Fatal(err)
		}
	}
	// make sure `article_to_tag` exists
	if !db.tableExists("article_to_tag") {
		fmt.Println("DB creating table `article_to_tag`...")
		_, err := db.Exec("CREATE TABLE article_to_tag( ArticleID INT, TagID INT, PRIMARY KEY (ArticleID, TagID) );")
		if err != nil {
			log.Fatal(err)
		}
	}
	// make sure `users` exists
	if !db.tableExists("users") {
		fmt.Println("DB creating table `users`...")
		_, err := db.Exec("CREATE TABLE users( ID INT AUTO_INCREMENT, Name VARCHAR(64) UNIQUE, Password VARCHAR(256), PRIMARY KEY (ID) );")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (db *DB) populate() {
	db.Exec(`INSERT INTO articles (URL) VALUES ("google.com/");`)
	db.Exec(`INSERT INTO tags (Name, Description) VALUES ("google", "it's google my duderino");`)
	db.Exec(`INSERT INTO tags (Name) VALUES ("search_engine");`)
	db.Exec(`INSERT INTO tags (Name) VALUES ("frogs");`)
	db.Exec(`INSERT INTO article_to_tag (ArticleID, TagID) VALUES (1,1);`)
	db.Exec(`INSERT INTO article_to_tag (ArticleID, TagID) VALUES (1,2);`)
}

// TagNameExists checks if a tag named `s` exists in a database, returning the tag ID && true/false
func (db *DB) TagNameExists(s string) (int64, bool) {
	rows, err := db.Query("SELECT ID FROM tags WHERE Name=? LIMIT 1;", s)
	if err != nil {
		return 0, false
	}
	exists := rows.Next()
	var id int64
	id = 0
	if exists {
		rows.Scan(&id)
	}
	rows.Close()
	return id, exists
}

// TagNamesExist is TagNameExists in a loop
func (db *DB) TagNamesExist(s ...string) ([]int64, bool) {
	res := []int64{}
	for _, t := range s {
		id, exists := db.TagNameExists(t)
		if !exists {
			return []int64{}, false
		}
		res = append(res, id)
	}
	return res, true
}

func stringOrNil(i string) interface{} {
	if len(i) > 0 {
		return i
	}
	return nil
}

func nullStringToString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

// ArticleTags finds all tags associated with an article ID
func (db *DB) ArticleTags(id int64) ([]Tag, error) {
	s := "SELECT t.*" +
		" FROM tags t INNER JOIN article_to_tag at ON t.ID = at.TagID" +
		" WHERE at.ArticleID = ?" +
		" AND at.TagID = t.ID;"
	rows, err := db.Query(s, id)
	if err != nil {
		return []Tag{}, err
	}

	tags := UnmarshalTags(rows)
	rows.Close()
	return tags, nil
}

// PopulateArticleTags adds tags to existing article struct based on article.ID
func (db *DB) PopulateArticleTags(article Article) Article {
	tags, err := db.ArticleTags(article.ID)
	if err != nil {
		return article
	}
	for _, t := range tags {
		article.Tags = append(article.Tags, t.Name)
	}
	return article
}

// UnmarshalArticles takes sql.Rows from the `article` table and parses it into an array of Article structs
// NOTE: does NOT populate `tags` field
func UnmarshalArticles(rows *sql.Rows) []Article {
	articles := []Article{}
	if rows == nil {
		return articles
	}
	for rows.Next() {
		var id int64
		name := ""
		var url sql.NullString
		var desc sql.NullString

		err := rows.Scan(&id, &name, &url, &desc)
		if err != nil {
			log.Println("Error unmarshalling article:", err)
		}
		articles = append(articles, Article{
			ID:          id,
			Name:        name,
			URL:         nullStringToString(url),
			Description: nullStringToString(desc),
		})
	}
	return articles
}

// UnmarshalTags takes sql.Rows from the `tags` table and parses it into an array of Tag structs
func UnmarshalTags(rows *sql.Rows) []Tag {
	tags := []Tag{}
	if rows == nil {
		return tags
	}
	for rows.Next() {
		var id int64
		name := ""
		var description sql.NullString

		err := rows.Scan(&id, &name, &description)
		if err != nil {
			log.Println("Error unmarshalling article:", err)
		}
		tags = append(tags, Tag{
			ID:          id,
			Name:        name,
			Description: nullStringToString(description),
		})
	}
	return tags
}

func findOrderby(s string) string {
	switch s {
	case "name":
		return "Name"
	case "description":
		return "Description"
	case "id":
		return "ID"
	default:
		return "ID"
	}
}

// ArticlesWithTagsSearch returns `limit` articles whose tags match all supplied tags, offset by `offset`, whose names OR description match `lookslike`
func (db *DB) ArticlesWithTagsSearch(tags []string, lookslike, orderby string, reverse bool, limit, offset int) ([]Article, error) {
	var itags []interface{}
	s := "SELECT a.*"

	if len(tags) > 0 {
		for _, t := range tags {
			itags = append(itags, t)
		}
		s += " FROM article_to_tag at INNER JOIN tags t ON at.TagID = t.ID INNER JOIN articles a ON at.ArticleID = a.ID" +
			" WHERE t.Name IN (?" + strings.Repeat(",?", len(tags)-1) + ")"
	} else {
		s += " FROM articles a WHERE TRUE"
	}
	if len(lookslike) > 0 {
		itags = append(itags, lookslike, lookslike)
		s += " AND (a.Name LIKE CONCAT('%',?,'%') OR a.Description LIKE CONCAT('%',?,'%'))"
	}
	s += " GROUP BY a.ID"
	if len(tags) > 0 {
		s += " HAVING COUNT(a.ID)=" + strconv.Itoa(len(tags))
	}
	if len(orderby) > 0 || reverse {
		s += " ORDER BY a." + findOrderby(orderby)
		if reverse {
			s += " DESC"
		} else {
			s += " ASC"
		}
	}
	if limit > 0 {
		itags = append(itags, limit)
		s += " LIMIT ?"
		if offset > 0 {
			itags = append(itags, offset)
			s += " OFFSET ?"
		}
	}
	s += ";"

	rows, err := db.Query(s, itags...)
	if err != nil {
		return []Article{}, err
	}

	articles := UnmarshalArticles(rows)
	rows.Close()

	for ii := range articles {
		articles[ii] = db.PopulateArticleTags(articles[ii])
	}

	return articles, nil
}

// TagSearch returns `limit` tags whose names are in `tags`, offset by `offset`, whose names match `lookslike`
// TagSearch returns a list of Tag structs given an array of tag names
func (db *DB) TagSearch(tags []string, lookslike, orderby string, reverse bool, limit int, offset int) ([]Tag, error) {
	s := "SELECT * FROM tags WHERE"

	var itags []interface{}
	if len(tags) > 0 {
		for _, t := range tags {
			itags = append(itags, t)
		}
		s += " Name IN (?" + strings.Repeat(",?", len(tags)-1) + ")"
	} else {
		s += " TRUE"
	}
	if len(lookslike) > 0 {
		itags = append(itags, lookslike, lookslike)
		s += " AND (Name LIKE CONCAT('%',?,'%') OR Description LIKE CONCAT('%',?,'%'))"
	}
	s += " GROUP BY ID"
	if len(orderby) > 0 || reverse {
		s += " ORDER BY " + findOrderby(orderby)
		if reverse {
			s += " DESC"
		} else {
			s += " ASC"
		}
	}
	if limit > 0 {
		itags = append(itags, limit)
		s += " LIMIT ?"
		if offset > 0 {
			itags = append(itags, offset)
			s += " OFFSET ?"
		}
	}
	s += ";"

	rows, err := db.Query(s, itags...)
	if err != nil {
		return []Tag{}, err
	}

	rtags := UnmarshalTags(rows)
	rows.Close()

	return rtags, nil
}

// ArticleByID searches for an articles with an ID, returning `nil` if not found
func (db *DB) ArticleByID(id int64) (*Article, error) {
	// TODO: make this return single article
	s := "SELECT * FROM articles WHERE ID=?;"
	rows, err := db.Query(s, id)
	if err != nil {
		return nil, err
	}
	articles := UnmarshalArticles(rows)
	rows.Close()

	for ii := range articles {
		articles[ii] = db.PopulateArticleTags(articles[ii])
	}

	if len(articles) >= 1 {
		return &articles[0], nil
	}
	return nil, nil
}

// TagByID searches for all tags with an ID, returning `nil` if not found
func (db *DB) TagByID(id int64) (*Tag, error) {
	// TODO: make this return single tag
	s := "SELECT * FROM tags WHERE ID=?;"
	rows, err := db.Query(s, id)
	if err != nil {
		return nil, err
	}
	tags := UnmarshalTags(rows)
	rows.Close()

	if len(tags) >= 1 {
		return &tags[0], nil
	}
	return nil, nil
}

// InsertArticleTag links an article to a tag, returning ID of inserted element
func (db *DB) InsertArticleTag(articleID int64, tagID int64) (int64, error) {
	res, err := db.Exec("INSERT INTO article_to_tag (ArticleID, TagID) VALUES (?, ?);", articleID, tagID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// InsertArticleTags updates an article's tags
func (db *DB) InsertArticleTags(id int64, tagIDs []int64) error {
	if len(tagIDs) < 1 {
		return nil
	}
	s := "INSERT INTO article_to_tag (ArticleID,TagID) VALUES (?,?)" +
		strings.Repeat(",(?,?)", len(tagIDs)-1) + ";"
	var params []interface{}
	for _, tagID := range tagIDs {
		params = append(params, id, tagID)
	}
	_, err := db.Exec(s, params...)
	return err
}

// InsertArticle inserts an article into a DB, linking tags if they exist and returning the article's ID, returning ID of inserted element
func (db *DB) InsertArticle(a Article) (int64, error) {
	res, err := db.Exec("INSERT INTO articles (Name, URL, Description) VALUES (?, ?, ?);", stringOrNil(a.Name), stringOrNil(a.URL), stringOrNil(a.Description))

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return id, err
	}
	for _, t := range a.Tags {
		if tagID, ok := db.TagNameExists(t); ok {
			db.InsertArticleTag(id, tagID)
		}
	}
	return id, nil
}

// InsertTag inserts a tag into a DB, returning ID of inserted element
func (db *DB) InsertTag(t Tag) (int64, error) {
	res, err := db.Exec("INSERT INTO tags (Name, Description) VALUES (?, ?);", stringOrNil(t.Name), stringOrNil(t.Description))

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return id, err
	}
	return id, nil
}

// RemoveArticleTags removes all article-tag links by articleID
func (db *DB) RemoveArticleTags(articleID int64) error {
	s := "DELETE FROM article_to_tag WHERE ArticleID=?;"
	_, err := db.Exec(s, articleID)
	return err
}

// RemoveTagsFromArticles removes all article-tag links by tagID
func (db *DB) RemoveTagsFromArticles(tagID int64) error {
	s := "DELETE FROM article_to_tag WHERE TagID=?;"
	_, err := db.Exec(s, tagID)
	return err
}

// RemoveArticle removes an article without touching article-tag links
func (db *DB) RemoveArticle(id int64) error {
	s := "DELETE FROM articles WHERE ID=?;"
	_, err := db.Exec(s, id)
	return err
}

// RemoveTag removes a tag without touching article-tag links
func (db *DB) RemoveTag(id int64) error {
	s := "DELETE FROM tags WHERE ID=?;"
	_, err := db.Exec(s, id)
	return err
}

// UpdateArticle updates an article's information, BUT NOT TAGS
func (db *DB) UpdateArticle(id int64, article Article) error {
	s := "UPDATE articles SET Name=?, URL=?, Description=? WHERE ID=?;"
	_, err := db.Exec(s, stringOrNil(article.Name), stringOrNil(article.URL), stringOrNil(article.Description), id)
	return err
}

// UpdateTag updates a tag's information
func (db *DB) UpdateTag(id int64, tag Tag) error {
	s := "UPDATE tags SET Name=?, Description=? WHERE ID=?;"
	_, err := db.Exec(s, stringOrNil(tag.Name), stringOrNil(tag.Description), id)
	return err
}

// UnmarshalUsers unmarshalls users
func UnmarshalUsers(rows *sql.Rows) []User {
	users := []User{}
	if rows == nil {
		return users
	}
	for rows.Next() {
		var id int64
		name := ""
		passwd := ""

		err := rows.Scan(&id, &name, &passwd)
		if err != nil {
			log.Println("Error unmarshalling article:", err)
		}
		users = append(users, User{
			ID:     id,
			Name:   name,
			Passwd: passwd,
		})
	}
	return users
}

// UserByName returns a user with a name, or nil
func (db *DB) UserByName(name string) (*User, error) {
	s := "SELECT * FROM users WHERE Name=?;"
	rows, err := db.Query(s, name)
	if err != nil {
		return nil, err
	}
	users := UnmarshalUsers(rows)
	if len(users) < 1 {
		return nil, nil
	} else if len(users) > 1 {
		log.Println("WARNING: multiple users with name", users)
	}
	return &users[0], nil
}

// InsertUser inserts a user
func (db *DB) InsertUser(user User) (int64, error) {
	s := `INSERT INTO users (Name, Password) VALUES (?, ?);`
	res, err := db.Exec(s, user.Name, user.Passwd)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, nil
}

// WARNING: vulnerable to SQL injection
func (db *DB) tableExists(name string) bool {
	rows, err := db.Query("SELECT 1 FROM " + name + " LIMIT 1;")
	if rows != nil {
		rows.Close()
	}
	return err == nil
}

func printRows(rows *sql.Rows) {
	for rows.Next() {
		var line string
		rows.Scan(&line)
		fmt.Println(line)
	}
}

/*

create articles table
CREATE TABLE articles( ID INT AUTO_INCREMENT, Name VARCHAR(512), URL VARCHAR(512) NOT NULL, PRIMARY KEY (ID) );

MariaDB [praxis_test_DB]> DESCRIBE articles;
+-------------+---------------+------+-----+---------+----------------+
| Field       | Type          | Null | Key | Default | Extra          |
+-------------+---------------+------+-----+---------+----------------+
| ID          | int(11)       | NO   | PRI | NULL    | auto_increment |
| Name        | varchar(512)  | NO   |     | NULL    |                |
| URL         | varchar(512)  | YES  |     | NULL    |                |
| Description | varchar(1024) | YES  |     | NULL    |                |
+-------------+---------------+------+-----+---------+----------------+
4 rows in set (0.008 sec)



tags table
CREATE TABLE tags( ID INT AUTO_INCREMENT, Name VARCHAR(16) UNIQUE, Description VARCHAR(256), PRIMARY KEY (ID, Name) );

MariaDB [praxis_test_DB]> DESCRIBE tags;
+-------------+--------------+------+-----+---------+----------------+
| Field       | Type         | Null | Key | Default | Extra          |
+-------------+--------------+------+-----+---------+----------------+
| ID          | int(11)      | NO   | PRI | NULL    | auto_increment |
| Name        | varchar(16)  | NO   | PRI | NULL    |                |
| Description | varchar(256) | YES  |     | NULL    |                |
+-------------+--------------+------+-----+---------+----------------+
3 rows in set (0.002 sec)



create article to tag table
CREATE TABLE article_to_tag( ArticleID INT, TagID INT, PRIMARY KEY (ArticleID, TagID) );

MariaDB [praxis_test_DB]> DESCRIBE article_to_tag;
+-----------+---------+------+-----+---------+-------+
| Field     | Type    | Null | Key | Default | Extra |
+-----------+---------+------+-----+---------+-------+
| ArticleID | int(11) | NO   | PRI | NULL    |       |
| TagID     | int(11) | NO   | PRI | NULL    |       |
+-----------+---------+------+-----+---------+-------+
2 rows in set (0.001 sec)

# url
INSERT INTO articles (URL) VALUES ("google.com/");

# make tag
INSERT INTO tags (Name, Description) VALUES ("google", "it's google my duderino");
INSERT INTO tags (Name) VALUES ("search_engine");
INSERT INTO tags (Name) VALUES ("frogs");

# relations
INSERT INTO article_to_tag (ArticleID, TagID) VALUES (1,1);
INSERT INTO article_to_tag (ArticleID, TagID) VALUES (1,2);

SELECT Name FROM article_to_tag att INNER JOIN articles a ON a.ID = att.ArticleID INNER JOIN tags t ON t.ID = att.TagID;



SELECT * FROM article_to_tag at, articles a, tags t WHERE t.ID = at.TagID AND a.ID = at.ArticleID AND (t.Name IN ('testtag','testtag3')) GROUP BY a.ID HAVING COUNT(a.ID)=2;

*/
