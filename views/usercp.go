package views

import (
	"github.com/flosch/pongo2"
	"lambda.sx/marcus/lambdago/session"
	"lambda.sx/marcus/lambdago/settings"
	"lambda.sx/marcus/lambdago/sql"
	"net/http"
	"strings"
	"upper.io/db"
)

//Compile the templates on startup for a speed boost
var usercpTpl = pongo2.Must(pongo2.FromFile("templates/usercp.html"))

func HandleUserCP(r *http.Request, w http.ResponseWriter) (error, string) {
	user := session.GetUser(r, w)
	if user.ID == 0 { // Not logged in
		http.Redirect(w, r, "/login", 302)
		return nil, ""
	}
	rendered_user_cp, err := usercpTpl.Execute(pongo2.Context{
		"user":   user,
		"themes": settings.Themes,
	})
	if err != nil {
		return err, ""
	}
	return nil, rendered_user_cp
}

func HandleToggleEncryption(r *http.Request, w http.ResponseWriter) (error, string) {
	user := session.GetUser(r, w)
	if user.ID == 0 { // Not logged in
		http.Redirect(w, r, "/login", 302)
		return nil, ""
	}
	user.EncryptionEnabled = !user.EncryptionEnabled
	col, _ := sql.Connection().Collection("users")
	col.Find(db.Cond{"ID": user.ID}).Update(user)

	// Bring back to user cp
	http.Redirect(w, r, "/usercp", 302)
	return nil, ""
}

func HandleSetTheme(r *http.Request, w http.ResponseWriter) (error, string) {
	if r.Method == "GET" {
		user := session.GetUser(r, w)
		if user.ID == 0 { // Not logged in
			http.Redirect(w, r, "/login", 302)
			return nil, ""
		}

		themeName := strings.ToLower(r.FormValue("name"))

		themeExists := false
		for _, name := range settings.Themes {
			if themeName == strings.ToLower(name) {
				themeExists = true
			}
		}

		if themeExists {
			user.ThemeName = themeName
			col, _ := sql.Connection().Collection("users")
			col.Find(db.Cond{"ID": user.ID}).Update(user)
		}

		// Bring back to user cp
		http.Redirect(w, r, "/usercp", 302)
		return nil, ""
	} else {
		return nil, "POST NOT SUPPORTED"
	}
}
