package settings

import (
	"upper.io/db"
	"upper.io/db/mysql"
)

var dbsettings = mysql.ConnectionURL{
	Address:  db.Socket("/var/run/mysqld/mysqld.sock"),
	Database: "lambda_go",
	User:     "lambda",
	Password: "lambda", // CHANGE FOR PRODUCTION
}

func DBSettings() mysql.ConnectionURL {
	return dbsettings
}
