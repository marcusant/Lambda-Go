package sql

import (
	"database/sql"
	"lambda.sx/marcus/lambdago/settings"
	"log"
	"upper.io/db"
	"upper.io/db/mysql"
)

var sqlConn db.Database = nil

func Init() {
	// Start a new SSL session
	sess, err := db.Open(mysql.Adapter, settings.DBSettings())
	if err != nil {
		log.Fatalf("SQL connection failed! %q\n", err)
		defer Shutdown()
	} else {
		sqlConn = sess
		// Create all of the tables for Lambda
		createTables()
	}
}

// Shutdown closes the SQL connection
func Shutdown() {
	sqlConn.Close()
}

// Connection returns the current MySQL connection
func Connection() db.Database {
	return sqlConn
}

func createTables() {
	driver := sqlConn.Driver().(*sql.DB)
	// Create users table
	driver.Query("CREATE TABLE IF NOT EXISTS users (" +
		"id MEDIUMINT UNSIGNED NOT NULL AUTO_INCREMENT," +
		"username VARCHAR(32) NOT NULL," +
		"password VARCHAR(128) NOT NULL," +
		"creation_date Date NOT NULL," +
		"apikey VARCHAR(64) NOT NULL," +
		"encryption_enabled BOOL NOT NULL," +
		"theme_name VARCHAR(32) NOT NULL," +
		"primary key(id)" +
		")")
}
