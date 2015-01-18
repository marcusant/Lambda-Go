package views

import (
	"github.com/flosch/pongo2"
	"net/http"
)

//Compile the templates on startup for a speed boost
var indexTpl = pongo2.Must(pongo2.FromFile("templates/index.html"))

func HandleIndex(r *http.Request) (error, string) {
	rendered_index, err := indexTpl.Execute(pongo2.Context{
	//Whatever context
	})
	if err != nil {
		return err, ""
	}
	return nil, rendered_index
}
