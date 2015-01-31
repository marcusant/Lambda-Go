package views

import (
	"github.com/flosch/pongo2"
	"lambda.sx/marcus/lambdago/models"
	"lambda.sx/marcus/lambdago/session"
	"lambda.sx/marcus/lambdago/settings"
	"lambda.sx/marcus/lambdago/sql"
	"net/http"
	"time"
)

var pasteTpl = pongo2.Must(pongo2.FromFile("templates/paste.html"))
var viewPasteTpl = pongo2.Must(pongo2.FromFile("templates/viewpaste.html"))

func HandlePaste(r *http.Request, w http.ResponseWriter) (error, string) {
	user := session.GetUser(r, w)
	if user.ID == 0 { // Not logged in
		http.Redirect(w, r, "/login", 302)
		return nil, ""
	}
	if r.Method != "POST" {
		rendered_paste_page, _ := pasteTpl.Execute(pongo2.Context{
			"user":  user,
			"nocdn": !settings.UseCDN,
		})
		return nil, rendered_paste_page
	} else {
		encr := r.PostFormValue("encr")
		if encr == "" {
			return nil, "Error: no json posted"
		}
		pasteName := genFilename() //TODO edit genFilename to check pastes along with the images for existing

		col, _ := sql.Connection().Collection("pastes")
		col.Append(models.Paste{
			Owner:       user.ID,
			Name:        pasteName,
			UploadDate:  time.Now(),
			ContentJson: encr,
		})
		return nil, pasteName
	}

}

func HandleViewPaste(r *http.Request, w http.ResponseWriter, json string) (error, string) {
	renderedViewPaste, _ := viewPasteTpl.Execute(pongo2.Context{
		"content": json,
		"nocdn":   !settings.UseCDN,
	})
	return nil, renderedViewPaste
}
