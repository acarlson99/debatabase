## Setup

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

MariaDB [(none)]> CREATE DATABASE praxisDB;
Query OK, 1 row affected (0.000 sec)

MariaDB [(none)]> 
```

### Server

```
export APP_ENV=dev      # for local development
# set DB uname/password
go run .
```
