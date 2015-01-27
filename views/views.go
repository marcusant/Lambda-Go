package views

import (
	"github.com/flosch/pongo2"
	"lambda.sx/marcus/lambdago/session"
	"net/http"
	"strings"
)

// Compile the templates on startup for a speed boost
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

// Try to find a paste or image, or serve 404
func HandleDefault(r *http.Request, w http.ResponseWriter) (error, string) {
	url := strings.Split(r.URL.String(), "?")[0]
	if url != "" {
		url = url[1:] // Remove "/" from before request
	}
	for _, ext := range allowedTypes {
		path := "uploads/" + url + "." + ext
		if fileExists(path) {
			http.ServeFile(w, r, path)
			return nil, ""
		}
	}
	return nil, "404"
}
