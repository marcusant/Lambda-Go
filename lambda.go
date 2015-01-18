package main

import (
	"fmt"
	"lambda.sx/marcus/lambdago/models"
	"lambda.sx/marcus/lambdago/sql"
	"lambda.sx/marcus/lambdago/views"
	"net/http"
	"time"
)

type ViewFunc func(*http.Request) (error, string)

var urlMap = map[string]ViewFunc{
	"/": views.HandleIndex,
}

func main() {
	staticServer := http.FileServer(http.Dir("static"))
	sql.Init()
	testUser := models.User{Username: "test", Password: "testPass", CreationDate: time.Now()}
	models.Save(testUser)
	http.Handle("/static/", http.StripPrefix("/static/", staticServer))
	http.HandleFunc("/", handler)

	http.ListenAndServe(":9000", nil)
}

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
	fmt.Fprint(w, responseHtml)
}
