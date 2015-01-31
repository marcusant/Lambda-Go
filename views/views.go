package views

import (
	"github.com/flosch/pongo2"
	"lambda.sx/marcus/lambdago/models"
	"lambda.sx/marcus/lambdago/session"
	"lambda.sx/marcus/lambdago/settings"
	"lambda.sx/marcus/lambdago/sql"
	"mime"
	"net/http"
	"strings"
	"upper.io/db"
)

// Compile the templates on startup for a speed boost
var indexTpl = pongo2.Must(pongo2.FromFile("templates/index.html"))
var aboutTpl = pongo2.Must(pongo2.FromFile("templates/about.html"))
var fourohfourTpl = pongo2.Must(pongo2.FromFile("templates/404.html"))

func HandleIndex(r *http.Request, w http.ResponseWriter) (error, string) {
	user := session.GetUser(r, w)
	rendered_index, err := indexTpl.Execute(pongo2.Context{
		"user":  user,
		"nocdn": !settings.UseCDN,
	})
	if err != nil {
		return err, ""
	}
	return nil, rendered_index
}

func HandleAbout(r *http.Request, w http.ResponseWriter) (error, string) {
	user := session.GetUser(r, w)
	rendered_about, err := aboutTpl.Execute(pongo2.Context{
		"user":  user,
		"nocdn": !settings.UseCDN,
	})
	if err != nil {
		return err, ""
	}
	return nil, rendered_about
}

// Try to find a paste or image, or serve 404
func HandleDefault(r *http.Request, w http.ResponseWriter) (error, string) {
	url := strings.Split(r.URL.String(), "?")[0]
	if url != "" {
		url = url[1:] // Remove "/" from before request
	}
	url = strings.Split(url, ".")[0] //Ignore any extensions
	for _, ext := range allowedTypes {
		path := "uploads/" + url + "." + ext
		if fileExists(path) {
			mimetype := mime.TypeByExtension("." + ext)
			w.Header().Set("Content-Type", mimetype)
			w.Header().Set("Cache-Control", "public, max-age=259200")
			// Don't let browsers guess a mime type other than the one we say.
			// Prevents users from serving files with a different extension from their true type.
			w.Header().Set("X-Content-Type-Options", "nosniff")
			http.ServeFile(w, r, path)
			return nil, ""
		}
	}

	col, _ := sql.Connection().Collection("pastes")
	result := col.Find(db.Cond{"name": url})
	cnt, err := result.Count()
	if err == nil && cnt > 0 {
		var paste models.Paste
		result.One(&paste)
		return HandleViewPaste(r, w, paste.ContentJson)
	}

	user := session.GetUser(r, w)
	rendered_404, err := fourohfourTpl.Execute(pongo2.Context{
		"user":  user,
		"nocdn": !settings.UseCDN,
	})
	if err != nil {
		return err, ""
	}
	w.WriteHeader(404)
	return nil, rendered_404
}
