package views

import (
	"github.com/flosch/pongo2"
	"lambda.sx/marcus/lambdago/models"
	"lambda.sx/marcus/lambdago/session"
	"lambda.sx/marcus/lambdago/settings"
	"lambda.sx/marcus/lambdago/sql"
	"net/http"
	"os"
	"strings"
	"upper.io/db"
)

//Compile the templates on startup for a speed boost
var usercpTpl = pongo2.Must(pongo2.FromFile("templates/usercp.html"))
var manageUploadsTpl = pongo2.Must(pongo2.FromFile("templates/manageuploads.html"))

func HandleUserCP(r *http.Request, w http.ResponseWriter) (error, string) {
	user := session.GetUser(r, w)
	if user.ID == 0 { // Not logged in
		http.Redirect(w, r, "/login", 302)
		return nil, ""
	}
	rendered_user_cp, err := usercpTpl.Execute(pongo2.Context{
		"user":   user,
		"themes": settings.Themes,
		"nocdn":  !settings.UseCDN,
	})
	if err != nil {
		return err, ""
	}
	return nil, rendered_user_cp
}

func HandleManageUploads(r *http.Request, w http.ResponseWriter) (error, string) {
	user := session.GetUser(r, w)
	if user.ID == 0 { // Not logged in
		http.Redirect(w, r, "/login", 302)
		return nil, ""
	}

	pasteCol, _ := sql.Connection().Collection("pastes")
	var pastes []models.Paste
	pasteCol.Find(db.Cond{"owner": user.ID}).Limit(10).All(&pastes)

	fileCol, _ := sql.Connection().Collection("files")
	var files []models.File
	fileCol.Find(db.Cond{"owner": user.ID}).Sort("-upload_date").Limit(10).All(&files)

	rendered_manage_uploads, err := manageUploadsTpl.Execute(pongo2.Context{
		"user":           user,
		"pastes":         pastes,
		"images":         files,
		"img_extensions": []string{".png", ".jpg"},
		"nocdn":          !settings.UseCDN,
	})
	if err != nil {
		return err, ""
	}
	return nil, rendered_manage_uploads
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

func HandleDelete(r *http.Request, w http.ResponseWriter) (error, string) {
	user := session.GetUser(r, w)
	if user.ID == 0 { // Not logged in
		http.Redirect(w, r, "/login", 302)
		return nil, ""
	}

	filename := r.FormValue("file")
	if filename == "" {
		return nil, "No filename sent!"
	}
	col, _ := sql.Connection().Collection("files")
	var file models.File
	result := col.Find(db.Cond{"owner": user.ID, "name": filename})
	result.One(&file)
	if file.ID == 0 {
		// Try to find a paste
		col, _ := sql.Connection().Collection("pastes")
		result = col.Find(db.Cond{"owner": user.ID, "name": filename})
		result.One(&file)
		if file.ID == 0 {
			return nil, "File not found!"
		}
	}
	os.Remove("uploads/" + file.Name + file.Extension)
	result.Remove()

	// Bring back to manage uploads
	http.Redirect(w, r, "/manageuploads", 302)
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
