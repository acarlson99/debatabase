package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// DB internal database struct
type DB struct {
	*sql.DB
}

// DBConnect creates connection to database with credentials
func DBConnect(uname, password, dbname string) (*DB, error) {
	fmt.Println("Connecting to database...")
	connStr := fmt.Sprintf("%s:%s@/%s", uname, password, dbname)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// Init things
func (db *DB) Init() {
	// check if `articles` exists
	fmt.Println("Initializing database...")
	if !db.tableExists("articles") {
		fmt.Println("DB creating table `articles`...")
		_, err := db.Exec("CREATE TABLE articles( ID INT AUTO_INCREMENT, Name VARCHAR(512), URL VARCHAR(512) NOT NULL, PRIMARY KEY (ID) );")
		if err != nil {
			log.Fatal(err)
		}
	}
	// check if `tags` exists
	if !db.tableExists("tags") {
		fmt.Println("DB creating table `tags`...")
		_, err := db.Exec("CREATE TABLE tags( ID INT AUTO_INCREMENT, Name VARCHAR(16) UNIQUE, Description VARCHAR(256), PRIMARY KEY (ID, Name) );")
		if err != nil {
			log.Fatal(err)
		}
	}
	// check if `article_to_tag` exists
	if !db.tableExists("article_to_tag") {
		fmt.Println("DB creating table `article_to_tag`...")
		_, err := db.Exec("CREATE TABLE article_to_tag( ArticleID INT, TagID INT, PRIMARY KEY (ArticleID, TagID) );")
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

// TagExists checks if tag `t` exists in a database, returning the tag ID
func (db *DB) TagExists(t string) (int64, bool) {
	rows, err := db.Query("SELECT ID FROM tags WHERE Name=? LIMIT 1;", t)
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

func UnmarshalArticle(rows *sql.Rows) []Article {
	articles := []Article{}
	if rows == nil {
		return articles
	}
	for rows.Next() {
		id := 0
		name := ""
		url := ""

		rows.Scan(&id, &name, &url)
		articles = append(articles, Article{
			Name: name,
			URL:  url,
		})
	}
	return articles
}

// ArticlesWithTags returns `limit` articles whose tags match all supplied tags, offset by `offset`
func (db *DB) ArticlesWithTags(tags []string, offset, limit int) ([]Article, error) {
	if len(tags) == 0 {
		return []Article{}, nil
	}

	var itags []interface{}
	for _, t := range tags {
		itags = append(itags, t)
	}
	s := "SELECT a.* " +
		"FROM article_to_tag at, articles a, tags t " +
		"WHERE t.ID = at.TagID " +
		"AND a.ID = at.ArticleID " +
		"AND (t.Name IN (?" + strings.Repeat(",?", len(itags)-1) + ")) " +
		"GROUP BY a.ID " +
		"HAVING COUNT(a.ID)=" + strconv.Itoa(len(itags)) + ";"
	rows, err := db.Query(s, itags...)

	// TODO: add tags to article
	if err != nil {
		return []Article{}, err
	}
	defer rows.Close()
	return UnmarshalArticle(rows), nil
}

// InsertArticleTag links an article to a tag
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

// InsertTag inserts a tag into a DB
func (db *DB) InsertTag(t Tag) (int64, error) {
	res, err := db.Exec("INSERT INTO tags (Name, Description) VALUES (?, ?);", t.Name, t.Description)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return id, err
	}
	return id, nil
}

// InsertArticle inserts an article into a DB, linking tags if they exist and returning the article's ID
func (db *DB) InsertArticle(a Article) (int64, error) {
	res, err := db.Exec("INSERT INTO articles (Name, URL) VALUES (?, ?);", a.Name, a.URL)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return id, err
	}
	for _, t := range a.Tags {
		if tagID, ok := db.TagExists(t); ok {
			db.InsertArticleTag(id, tagID)
		}
	}
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
		err := rows.Scan(&line)
		if err != nil {
			log.Println("ERROR", err)
		}
		log.Println(line)
	}
}

/*

create articles table
CREATE TABLE articles( ID INT AUTO_INCREMENT, Name VARCHAR(512), URL VARCHAR(512) NOT NULL, PRIMARY KEY (ID) );

MariaDB [praxisDB]> DESCRIBE articles;
+-------+--------------+------+-----+---------+----------------+
| Field | Type         | Null | Key | Default | Extra          |
+-------+--------------+------+-----+---------+----------------+
| ID    | int(11)      | NO   | PRI | NULL    | auto_increment |
| Name  | varchar(512) | YES  |     | NULL    |                |
| URL   | varchar(512) | NO   |     | NULL    |                |
+-------+--------------+------+-----+---------+----------------+
3 rows in set (0.011 sec)



tags table
CREATE TABLE tags( ID INT AUTO_INCREMENT, Name VARCHAR(16) UNIQUE, Description VARCHAR(256), PRIMARY KEY (ID, Name) );

MariaDB [praxisDB]> DESCRIBE tags;
+-------------+--------------+------+-----+---------+----------------+
| Field       | Type         | Null | Key | Default | Extra          |
+-------------+--------------+------+-----+---------+----------------+
| ID          | int(11)      | NO   | PRI | NULL    | auto_increment |
| Name        | varchar(16)  | NO   | PRI | NULL    |                |
| Description | varchar(256) | YES  |     | NULL    |                |
+-------------+--------------+------+-----+---------+----------------+
3 rows in set (0.003 sec)



create article to tag table
CREATE TABLE article_to_tag( ArticleID INT, TagID INT, PRIMARY KEY (ArticleID, TagID) );

MariaDB [praxisDB]> DESCRIBE article_to_tag;
+-----------+---------+------+-----+---------+-------+
| Field     | Type    | Null | Key | Default | Extra |
+-----------+---------+------+-----+---------+-------+
| ArticleID | int(11) | NO   | PRI | NULL    |       |
| TagID     | int(11) | NO   | PRI | NULL    |       |
+-----------+---------+------+-----+---------+-------+
2 rows in set (0.002 sec)

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
