package session

import (
	"github.com/gorilla/sessions"
	"lambda.sx/marcus/lambdago/models"
	"lambda.sx/marcus/lambdago/settings"
	"lambda.sx/marcus/lambdago/sql"
	"net/http"
	"upper.io/db"
)

var Store = sessions.NewCookieStore([]byte(settings.SecretKey)) // TODO move to using a server-held store (redis?) with keys held by the client

func GetUser(r *http.Request, w http.ResponseWriter) models.User {
	var user models.User
	session, err := Store.Get(r, "lambda")
	if err != nil {
		return user
	}
	id, ok := session.Values["userid"].(uint)
	if !ok {
		return user
	}
	if id > 0 {
		col, err := sql.Connection().Collection("users")
		if err == nil {
			col.Find(db.Cond{"id": id}).One(&user)
		}
	}
	return user
}
