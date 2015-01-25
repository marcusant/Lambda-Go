package settings

import (
	"upper.io/db"
	"upper.io/db/mysql"
)

// CHANGE IN PRODUCTION
// KEEP SECRET!!!
var SecretKey = "dongLyfe420"

var Themes = [...]string{"material", "space"}

var dbsettings = mysql.ConnectionURL{
	Address:  db.Socket("/var/run/mysqld/mysqld.sock"),
	Database: "lambda_go",
	User:     "lambda",
	Password: "lambda", // CHANGE FOR PRODUCTION
}

func DBSettings() mysql.ConnectionURL {
	return dbsettings
}
