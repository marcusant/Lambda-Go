package migrate

import (
	gsql "database/sql"
	"lambda.sx/marcus/lambdago/models"
	"lambda.sx/marcus/lambdago/sql"
	"log"
	"time"
	"upper.io/db"
	"upper.io/db/mysql"
)

// Migrates the Lambda database from Django to Go. DELETES THE GO DATABASE!!!
func MigrateDB() {
	// Drop the tables
	driver := sql.Connection().Driver().(*gsql.DB)
	cons, _ := sql.Connection().Collections()
	for _, s := range cons {
		driver.Query("TRUNCATE " + s)
	}

	dbsettings := mysql.ConnectionURL{
		Address:  db.Socket("/var/run/mysqld/mysqld.sock"),
		Database: "djlambda",
		User:     "lambda",
		Password: "lambda", // CHANGE FOR PRODUCTION
	}
	djsess, err := db.Open(mysql.Adapter, dbsettings)
	if err != nil {
		log.Fatalf("SQL connection failed! %q\n", err)
		return
	} else {
		type DjUser struct {
			ID         int       `db:"id"`
			Password   string    `db:"password"`
			Username   string    `db:"username"`
			DateJoined time.Time `db:"date_joined"`
		}
		type DjLambdaUser struct {
			ID                int    `db:"id"`
			UserID            int    `db:"user_id"`
			EncryptionEnabled bool   `db:"encryption_enabled"`
			ThemeName         string `db:"theme_name"`
			ApiKey            string `db:"apikey"`
		}

		djusers, _ := djsess.Collection("auth_user")
		djlambdausers, _ := djsess.Collection("djlambda_lambdauser")

		var djlambdauserlist []DjLambdaUser
		djlambdausers.Find(db.Cond{}).All(&djlambdauserlist)

		for _, u := range djlambdauserlist {
			var djuser DjUser
			djusers.Find(db.Cond{"id": u.UserID}).One(&djuser)

			gouser := models.User{
				Username:          djuser.Username,
				Password:          djuser.Password,
				CreationDate:      djuser.DateJoined,
				ApiKey:            u.ApiKey,
				EncryptionEnabled: u.EncryptionEnabled,
				ThemeName:         u.ThemeName,
			}
			userCol, _ := sql.Connection().Collection("users")
			userCol.Append(gouser)
		}
	}
}
