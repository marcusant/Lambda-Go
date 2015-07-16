package sql

import (
	"database/sql"
	"log"

	"lambda.sx/marcus/lambdago/settings"
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
	driver.Query("CREATE TABLE IF NOT EXISTS files (" +
		"id MEDIUMINT UNSIGNED NOT NULL AUTO_INCREMENT," +
		"owner MEDIUMINT UNSIGNED NOT NULL," +
		"name VARCHAR(16) NOT NULL," +
		"extension VARCHAR(5) NOT NULL," +
		"upload_date Date NOT NULL," +
		"encrypted BOOL NOT NULL," +
		"local_name VARCHAR(128) NOT NULL," +
		"primary key(id)" +
		")")
	driver.Query("CREATE TABLE IF NOT EXISTS pastes (" +
		"id MEDIUMINT UNSIGNED NOT NULL AUTO_INCREMENT," +
		"owner MEDIUMINT UNSIGNED NOT NULL," +
		"name VARCHAR(16) NOT NULL," +
		"upload_date Date NOT NULL," +
		"content_json VARCHAR(50000) NOT NULL," +
		"is_code BOOL NOT NULL," +
		"primary key(id)" +
		")")
}
