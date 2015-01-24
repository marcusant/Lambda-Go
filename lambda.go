package main

import (
	"fmt"
	"lambda.sx/marcus/lambdago/sql"
	"lambda.sx/marcus/lambdago/views"
	"net/http"
)

// Function that is called for a view
type ViewFunc func(*http.Request) (error, string)

var urlMap = map[string]ViewFunc{
	"/":         views.HandleIndex,
	"/register": views.HandleRegister,
}

func main() {
	// Start up SQL connection
	sql.Init()

	// Create a static server for serving things in the static/ directory
	staticServer := http.FileServer(http.Dir("static"))
	// Use a static file server for /static/
	http.Handle("/static/", http.StripPrefix("/static/", staticServer))
	// Use handler for everything else
	http.HandleFunc("/", handler)

	// Start listening on port 9000 on every interface
	http.ListenAndServe(":9000", nil)
}

// handler dispatches requests to their view or serves an error
func handler(w http.ResponseWriter, r *http.Request) {
	vfunc := urlMap[r.URL.String()]
	if vfunc == nil { // 404, no view to handle request
		//TODO 404 page
		panic("404 NOT YET IMPLEMENTED")
	}
	err, responseHtml := vfunc(r)
	if err != nil { // 500, server dun goofed
		//TODO 500 page
		fmt.Println("500 at " + r.URL.String())
		panic(err)
	}
	fmt.Fprint(w, responseHtml) // Print the html to the response writer
}
