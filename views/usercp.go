package views

import (
	"github.com/flosch/pongo2"
	"lambda.sx/marcus/lambdago/session"
	"net/http"
)

//Compile the templates on startup for a speed boost
var usercpTpl = pongo2.Must(pongo2.FromFile("templates/usercp.html"))

func HandleUserCP(r *http.Request, w http.ResponseWriter) (error, string) {
	user := session.GetUser(r, w)
	if user.ID == 0 { // Not logged in
		http.Redirect(w, r, "/login", 302)
		return nil, ""
	}
	rendered_index, err := usercpTpl.Execute(pongo2.Context{
		"user": user,
	})
	if err != nil {
		return err, ""
	}
	return nil, rendered_index
}
