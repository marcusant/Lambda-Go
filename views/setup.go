package views

import (
	"github.com/flosch/pongo2"
	"lambda.sx/marcus/lambdago/session"
	"net/http"
)

// Compile the templates on startup for a speed boost
var setupLambdaTpl = pongo2.Must(pongo2.FromFile("templates/setup/lambda.html"))
var setupSharexTpl = pongo2.Must(pongo2.FromFile("templates/setup/sharex.html"))

func HandleSetupSharex(r *http.Request, w http.ResponseWriter) (error, string) {
	user := session.GetUser(r, w)
	if user.ID == 0 { // Not logged in
		http.Redirect(w, r, "/login", 302)
		return nil, ""
	}
	rendered_sharex, err := setupSharexTpl.Execute(pongo2.Context{
		"user": user,
	})
	if err != nil {
		return err, ""
	}
	return nil, rendered_sharex
}

func HandleSetupLambda(r *http.Request, w http.ResponseWriter) (error, string) {
	user := session.GetUser(r, w)
	if user.ID == 0 { // Not logged in
		http.Redirect(w, r, "/login", 302)
		return nil, ""
	}
	rendered_lambda, err := setupLambdaTpl.Execute(pongo2.Context{
		"user": user,
	})
	if err != nil {
		return err, ""
	}
	return nil, rendered_lambda
}
