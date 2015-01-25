package views

import (
	"github.com/flosch/pongo2"
	"lambda.sx/marcus/lambdago/session"
	"net/http"
)

//Compile the templates on startup for a speed boost
var indexTpl = pongo2.Must(pongo2.FromFile("templates/index.html"))

func HandleIndex(r *http.Request, w http.ResponseWriter) (error, string) {
	user := session.GetUser(r, w)
	rendered_index, err := indexTpl.Execute(pongo2.Context{
		"user": user,
	})
	if err != nil {
		return err, ""
	}
	return nil, rendered_index
}
