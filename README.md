## Setup

### DB

#### Fedora

```
sudo dnf install mariadb-server
systemctl start mariadb         # start
mysql_secure_installation       # setup
mysql -u root -ppassword        # connect
```

### Server

```
export APP_ENV=dev      # for local development
# set DB uname/password
go run .
```
