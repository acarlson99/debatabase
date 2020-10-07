# This is a debatabase

## Setup

### Dependencies

* go version go1.13.11
* mysql/mariadb
* [swaggo](https://github.com/swaggo/swag)
* node/npm

### DB

#### Fedora

```
$ sudo dnf install mariadb-server
$ systemctl start mariadb                         # start
$ mysql_secure_installation                       # setup
$ mysql -u root -ppassword                        # connect
Welcome to the MariaDB monitor.  Commands end with ; or \g.
Your MariaDB connection id is 28
Server version: 10.3.22-MariaDB MariaDB Server

Copyright (c) 2000, 2018, Oracle, MariaDB Corporation Ab and others.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

MariaDB [(none)]> CREATE DATABASE debatabaseDB;
Query OK, 1 row affected (0.000 sec)

MariaDB [(none)]> 
```

## Deploy

```
cd frontend
npm run build
cd ..

swag init -g server.go

export APP_ENV=dev      # for local development
export MYSQL_PASSWORD=password
export MYSQL_USER=root
export MYSQL_DBNAME=db_name
export HOST_ADDRESS=localhost
export HOST_PORT=9000
go run .
```

## Dev notes

#### source `.env`

`source <(sed 's/^/export /' .env)`

#### CURL

```bash
# upload tags
curl -L -i localhost:9000/api/upload/tag --data '{"name":"engine","description":"a thing that does"}'
curl -L -i localhost:9000/api/upload/tag --data '{"name":"search","description":"a thing that finds"}'
curl -L -i localhost:9000/api/upload/tag --data '{"name":"tank"}'
# upload article
curl -L -i localhost:9000/api/upload/article --data '{"name":"googel","url":"google.com","tags":["engine","search"]}'
# search for 'engine' tag
curl -L -i localhost:9000/api/search/tag?tags=engine
> [{"id":24,"name":"engine","description":"a thing that does"}]
# search for all articles tagged 'engine'
curl -L -i localhost:9000/api/search/article?tags=engine
> [{"id":1,"name":"googel","url":"google.com","description":"","tags":["engine","search"]}]

# upload from CSV -- NOTE: UNDOCUMENTED NOT INTENDED FOR ACTUAL USE
curl -L -i localhost:9000/api/upload/tag/csv --data "`cat resources/tags.csv`"
curl -L -i localhost:9000/api/upload/article/csv --data "`cat resources/articles.csv`"
```

# Endpoints

swagger API http://localhost:9000/swagger/index.html

## Search

All searches handle arguments identically

### URL params

* tags - search for specific tag names
* limit - return at most `limit` results
* offset - skip first `offset` results
* lookslike - filter for name/description matching `lookslike`
* orderby - order results by field.  Supported args are `name`, `description`, `id`(default)
* reverse - reverse results.  `true` or `false`

### Examples

```
Search for articles
GET /api/search/article
curl -L -i localhost:9000/api/search/article?tags=search&limit=1&offset=1&lookslike=ooooogle&orderby=name&reverse=true

Search for tags
GET /api/search/tag
curl -L -i localhost:9000/api/search/tag?tags=search&limit=1&offset=1&lookslike=ooooogle
```

## Upload

### Tag

`name` field required

```
POST /api/upload/tag

{
    "name": "engine",
    "description": "a machine designed to convert one form of energy into mechanical energy"
}
```

### Tag CSV

able to upload multiple tags in CSV format delimited by a single `'\n'`

```
POST /api/upload/tag/csv

name,description
engine,a machine
search,something designed to find things
```

### Article

`name` field required

```
POST /api/upload/article

{
    "name": "googel",
    "url": "google.com",
    "description": "a biiiiiiggg boy search engine",
    "tags": ["engine", "search"]
}
```

### Article CSV

able to upload multiple articles in CSV format delimited by a single `'\n'`

```
POST /api/upload/article/csv

name,url,description,tagsCSV
googel,google.com,a search engine,"engine,search"
```
