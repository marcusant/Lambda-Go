package migrate

import (
	gsql "database/sql"
	"lambda.sx/marcus/lambdago/models"
	"lambda.sx/marcus/lambdago/settings"
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
		Password: settings.DBSettings().Password,
	}
	djsess, err := db.Open(mysql.Adapter, dbsettings)
	if err != nil {
		log.Printf("SQL connection failed! %q\n", err)
		return
	} else {
		defer djsess.Close()
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
		djusers, err := djsess.Collection("auth_user")
		if err != nil {
			return // There was no django database with users, let's just bail out
		}
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

		type DjLambdaImage struct {
			ID         int       `db:"id"`
			OwnerID    int       `db:"owner"`
			Name       string    `db:"name"`
			Extension  string    `db:"extension"`
			UploadDate time.Time `db:"upload_date"`
			Encrypted  bool      `db:"encrypted"`
		}
		djimages, _ := djsess.Collection("djlambda_image")
		var djimagelist []DjLambdaImage
		djimages.Find(db.Cond{}).All(&djimagelist)
		for _, i := range djimagelist {
			goimage := models.File{
				Owner:      uint(i.OwnerID),
				Name:       i.Name,
				Extension:  i.Extension,
				UploadDate: i.UploadDate,
				Encrypted:  i.Encrypted,
				LocalName:  "N/A",
			}
			fileCol, _ := sql.Connection().Collection("files")
			fileCol.Append(goimage)
		}

		type DjLambdaPaste struct {
			ID         int       `db:"id"`
			OwnerID    int       `db:"owner"`
			Name       string    `db:"name"`
			ReqJson    string    `db:"req_json"`
			UploadDate time.Time `db:"creation_date"`
		}
		djpastes, _ := djsess.Collection("djlambda_paste")
		var djpastelist []DjLambdaPaste
		djpastes.Find(db.Cond{}).All(&djpastelist)
		for _, p := range djpastelist {
			gopaste := models.Paste{
				Owner:       uint(p.OwnerID),
				Name:        p.Name,
				ContentJson: p.ReqJson,
				UploadDate:  p.UploadDate,
				IsCode:      true,
			}
			pasteCol, _ := sql.Connection().Collection("pastes")
			pasteCol.Append(gopaste)
		}

	}
}
