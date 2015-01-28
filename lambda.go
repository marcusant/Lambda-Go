package main

import (
	"fmt"
	"lambda.sx/marcus/lambdago/migrate"
	"lambda.sx/marcus/lambdago/settings"
	"lambda.sx/marcus/lambdago/sql"
	"lambda.sx/marcus/lambdago/views"
	"net/http"
	"strings"
)

// Function that is called for a view
type ViewFunc func(*http.Request, http.ResponseWriter) (error, string)

var urlMap = map[string]ViewFunc{
	"/":              views.HandleIndex,
	"/register":      views.HandleRegister,
	"/login":         views.HandleLogin,
	"/logout":        views.HandleLogout,
	"/usercp":        views.HandleUserCP,
	"/toggleencrypt": views.HandleToggleEncryption,
	"/settheme":      views.HandleSetTheme,
	"/keycheck":      views.HandleVerifyKey,
	"/getkey":        views.HandleGetKey,
	"/upload":        views.HandleUpload,
	"/paste":         views.HandlePaste,
	"/p":             views.HandlePaste,
	"/about":         views.HandleAbout,
	"/manageuploads": views.HandleManageUploads,
	"/delete":        views.HandleDelete,
}

func main() {
	settings.Init()
	// Start up SQL connection
	sql.Init()
	migrate.MigrateDB()

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
	vfunc := urlMap[strings.Split(strings.ToLower(r.URL.String()), "?")[0]]
	if vfunc == nil { // 404, no view to handle request
		vfunc = views.HandleDefault
	}
	err, responseHtml := vfunc(r, w)
	if err != nil { // 500, server dun goofed
		//TODO 500 page
		fmt.Println("500 at " + r.URL.String())
		panic(err)
	}
	fmt.Fprint(w, responseHtml) // Print the html to the response writer
}
