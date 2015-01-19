package sql

import (
	"lambda.sx/marcus/lambdago/settings"
	"log"
	"upper.io/db"
	"upper.io/db/mysql"
)

var sqlConn db.Database = nil

func Init() {
	sess, err := db.Open(mysql.Adapter, settings.DBSettings())
	if err != nil {
		log.Fatalf("SQL connection failed! %q\n", err)
		defer Shutdown()
	} else {
		sqlConn = sess
		createTables()
	}
}

func Shutdown() {
	sqlConn.Close()
}

func Connection() db.Database {
	return sqlConn
}

func createTables() {
	driver = sqlConn.Driver().(*sql.DB)
	driver.Query("CREATE TABLE IF NOT EXISTS users (username VARCHAR(32), password VARCHAR(64), creation_date Date)")
}
