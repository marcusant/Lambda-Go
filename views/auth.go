package views

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/flosch/pongo2"
	"golang.org/x/crypto/pbkdf2"
	"lambda.sx/marcus/lambdago/models"
	"lambda.sx/marcus/lambdago/sql"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"upper.io/db"
)

var registerTpl = pongo2.Must(pongo2.FromFile("templates/register.html"))

func HandleRegister(r *http.Request) (error, string) {
	if r.Method == "POST" {
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		passwordTwo := r.PostFormValue("password2")

		col, err := sql.Connection().Collection("users")

		if err != nil {
			//TODO we probably don't want to actually output the error in production
			msg := fmt.Sprintf("SQL connection failed: %s", err)
			rendered_tpl, _ := registerTpl.Execute(pongo2.Context{
				"messages": [...]string{msg},
			})
			return nil, rendered_tpl
		}

		//Validate username input
		usernameLength := len([]rune(username))
		re := regexp.MustCompile("^[a-zA-Z0-9_]*$") //alphanumeric test
		var messages []string
		if usernameLength < 4 || !re.MatchString(username) {
			messages = append(messages, "Usernames must be longer than 3 characters, alphanumeric, and have no spaces")
		} else {
			cnt, _ := col.Find(db.Cond{"username": username}).Count()
			if cnt > 0 {
				messages = append(messages, "Username already in use")
			}
		}

		//Validate password input
		passwordLength := len([]rune(password))
		if passwordLength < 6 || strings.Contains(password, " ") {
			messages = append(messages, "Passwords must be longer than 6 characters and contain no spaces")
		}
		if password != passwordTwo {
			messages = append(messages, "Two passwords do not match")
		}

		if len(messages) > 0 { //We had an error
			rendered_tpl, err := registerTpl.Execute(pongo2.Context{
				"messages": messages,
			})
			if err != nil {
				return err, ""
			}
			return nil, rendered_tpl
		} else {
			iter := 12000
			salt := "fdsa" //TODO set to random base64 byte array
			encpass := pbkdf2.Key([]byte(password), []byte(salt), iter, 64, sha256.New)
			hash := base64.StdEncoding.EncodeToString(encpass)
			// From django docs: <algorithm>$<iterations>$<salt>$<hash>
			passentry := fmt.Sprintf("%s$%s$%s$%s", "pbkdf2_sha256", strconv.Itoa(iter), salt, hash)
			col.Append(models.User{
				Username:     username,
				Password:     passentry,
				CreationDate: time.Now(),
			})
			//TODO redirect user to home page
		}
	}
	rendered_tpl, err := registerTpl.Execute(pongo2.Context{
	//Whatever context
	})
	if err != nil {
		return err, ""
	}
	return nil, rendered_tpl
}
