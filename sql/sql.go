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
	}
}

func Shutdown() {
	sqlConn.Close()
}

func Connection() db.Database {
	return sqlConn
}
