package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DB internal database struct
type DB struct {
	*sql.DB
}

// DBConnect creates connection to database with credentials
func DBConnect(uname, password, dbname string) (*DB, error) {
	log.Println("Connecting to database...")
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
	log.Println("Initializing database...")
	if !db.tableExists("articles") {
		log.Println("Creating `articles` table...")
		_, err := db.Exec("CREATE TABLE articles( ID INT AUTO_INCREMENT, URL VARCHAR(512) NOT NULL, PRIMARY KEY (ID) );")
		if err != nil {
			log.Fatal(err)
		}
	}
	// check if `tags` exists
	if !db.tableExists("tags") {
		log.Println("Creating `tags` table...")
		_, err := db.Exec("CREATE TABLE tags( ID INT AUTO_INCREMENT, Name VARCHAR(16) UNIQUE, Description VARCHAR(256), PRIMARY KEY (ID, Name) );")
		if err != nil {
			log.Fatal(err)
		}
	}
	// check if `article_to_tag` exists
	if !db.tableExists("article_to_tag") {
		log.Println("Creating `article_to_tag` table...")
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
CREATE TABLE articles( ID INT AUTO_INCREMENT, URL VARCHAR(512) NOT NULL, PRIMARY KEY (ID) );

MariaDB [radancomDB]> DESCRIBE articles;
+-------+--------------+------+-----+---------+----------------+
| Field | Type         | Null | Key | Default | Extra          |
+-------+--------------+------+-----+---------+----------------+
| ID    | int(11)      | NO   | PRI | NULL    | auto_increment |
| URL   | varchar(512) | NO   |     | NULL    |                |
+-------+--------------+------+-----+---------+----------------+
2 rows in set (0.001 sec)



tags table
CREATE TABLE tags( ID INT AUTO_INCREMENT, Name VARCHAR(16) UNIQUE, Description VARCHAR(256), PRIMARY KEY (ID, Name) );

MariaDB [radancomDB]> DESCRIBE tags;
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

MariaDB [radancomDB]> DESCRIBE article_to_tag;
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

*/
