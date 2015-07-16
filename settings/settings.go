package settings

import (
	"io/ioutil"
	"log"
	"os"
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
	// MySQL username
	sqlUser := os.Getenv("LMDA_SQL_USER")
	if len(sqlUser) > 0 {
		dbsettings.User = sqlUser
	}

	// MySQL password
	sqlPw := os.Getenv("LMDA_SQL_PASS")
	if len(sqlPw) > 0 {
		dbsettings.Password = sqlPw
	}

	// MySQL address
	sqlAddr := os.Getenv("LMDA_SQL_ADDR")
	if len(sqlAddr) > 0 {
		if strings.HasSuffix(sqlAddr, ".sock") {
			dbsettings.Address = db.Socket(sqlAddr)
		} else {
			dbsettings.Address = db.Host(sqlAddr)
		}
	}

	// MySQL database name
	dbTitle := os.Getenv("LMDA_SQL_DB")
	if len(dbTitle) > 0 {
		dbsettings.Database = dbTitle
	}

	// Secret key
	secretKey := os.Getenv("LMDA_SECRET")
	if len(secretKey) > 0 {
		SecretKey = secretKey
	} else {
		log.Println("[WARNING] Using the default secret key. If in production, " +
			"define the environment variable LMDA_SECRET with your key.")
	}

	// Recapthca keys
	rPvtKey := os.Getenv("LMDA_RECAPTCHA_PRIVATE")
	if len(rPvtKey) > 0 {
		RecaptchaPrivateKey = rPvtKey
	} else {
		log.Println("[WARNING] No recaptcha private key defined. Captcha will not work. " +
			"Define one in the environment variable LMDA_RECAPTCHA_PRIVATE.")
	}
	rPubKey := os.Getenv("LMDA_RECAPTCHA_PUBLIC")
	if len(rPubKey) > 0 {
		RecaptchaPublicKey = rPubKey
	} else {
		log.Println("[WARNING] No recaptcha public key defined. Captcha will not work. " +
			"Define one in the environment variable LMDA_RECAPTCHA_PUBLIC.")
	}

	_, err := ioutil.ReadFile("../usecdn")
	if err == nil { //The file exists
		UseCDN = true
	}
}
