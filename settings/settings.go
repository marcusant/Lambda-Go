package settings

import (
	"io/ioutil"
	"strings"
	"upper.io/db"
	"upper.io/db/mysql"
)

// CHANGE IN PRODUCTION
// KEEP SECRET!!!
var SecretKey = "dongLyfe420"

var RecaptchaPrivateKey = ""
var RecaptchaPublicKey = ""

var UseCDN = false

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

func Init() {
	sqlInfoContents, err := ioutil.ReadFile("../mysqlauth")
	if err == nil {
		dbsettings.Password = strings.TrimSpace(string(sqlInfoContents))
	}
	secretKey, err := ioutil.ReadFile("../secretkey")
	if err == nil {
		SecretKey = strings.TrimSpace(string(secretKey))
	}

	rPvtKey, err := ioutil.ReadFile("../rcpvtkey")
	if err == nil {
		RecaptchaPrivateKey = strings.TrimSpace(string(rPvtKey))
	}
	rPubKey, err := ioutil.ReadFile("../rcpubkey")
	if err == nil {
		RecaptchaPublicKey = strings.TrimSpace(string(rPubKey))
	}

	_, err = ioutil.ReadFile("../usecdn")
	if err == nil { //The file exists
		UseCDN = true
	}
}
