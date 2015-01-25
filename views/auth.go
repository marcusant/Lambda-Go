package views

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/pbkdf2"
	"lambda.sx/marcus/lambdago/models"
	"lambda.sx/marcus/lambdago/session"
	"lambda.sx/marcus/lambdago/sql"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"upper.io/db"
)

var registerTpl = pongo2.Must(pongo2.FromFile("templates/register.html"))
var loginTpl = pongo2.Must(pongo2.FromFile("templates/login.html"))

func HandleRegister(r *http.Request, w http.ResponseWriter) (error, string) {
	session, err := session.Store.Get(r, "lambda")
	if err == nil && !session.IsNew { //They're already logged in, take them to the home page
		http.Redirect(w, r, "/", 302)
		return nil, ""
	}

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
			//Add the user
			passentry := hashPasswordDefault(password)
			col.Append(models.User{
				Username:     username,
				Password:     passentry,
				CreationDate: time.Now(),
			})
			var user models.User
			col.Find(db.Cond{"username": username}).One(&user)

			session.Values["userid"] = user.ID
			session.Save(r, w)
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

func HandleLogin(r *http.Request, w http.ResponseWriter) (error, string) {
	session, err := session.Store.Get(r, "lambda")
	if err == nil && !session.IsNew { //They're already logged in, take them to the home page
		http.Redirect(w, r, "/", 302)
		return nil, ""
	}

	if r.Method == "POST" {
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")

		var messages []string

		coll, err := sql.Connection().Collection("users")
		if err != nil {
			//TODO we probably don't want to actually output the error in production
			msg := fmt.Sprintf("SQL connection failed: %s", err)
			rendered_tpl, _ := registerTpl.Execute(pongo2.Context{
				"messages": [...]string{msg},
			})
			return nil, rendered_tpl
		}

		result := coll.Find(db.Cond{"username": username})
		count, err := result.Count()
		if err != nil || count < 1 {
			messages = append(messages, "No user exists with username")
		} else if count > 0 {
			var user models.User
			fmt.Println(user)
			result.One(&user)

			correctPass, _ := checkPassword(user, password)
			if correctPass {
				http.Redirect(w, r, "/", 302)
				return nil, ""
				// TODO give session
			} else {
				messages = append(messages, "Invalid password")
			}
		}

		rendered_tpl, err := loginTpl.Execute(pongo2.Context{
			"messages": messages,
		})
		if err != nil {
			return err, ""
		} else {
			return nil, rendered_tpl
		}

	}
	rendered_tpl, err := loginTpl.Execute(pongo2.Context{
	//Whatever context
	})
	if err != nil {
		return err, ""
	}
	return nil, rendered_tpl
}

func HandleLogout(r *http.Request, w http.ResponseWriter) (error, string) {
	session, _ := session.Store.Get(r, "lambda")

	// Tell the client to expire their session cookie
	session.Options = &sessions.Options{MaxAge: -1}
	session.Save(r, w)

	// Go home
	http.Redirect(w, r, "/", 302)
	return nil, ""
}

// Checks if the specified plaintext password matches the user's password
func checkPassword(user models.User, rawpass string) (bool, error) {
	dollaSplit := strings.Split(user.Password, "$")
	// From django docs: <algorithm>$<iterations>$<salt>$<hash>

	algoritm := dollaSplit[0]
	if algoritm != "pbkdf2_sha256" { // For right now, we only support this algorithm
		return false, errors.New("Algorithm not supported")
	}

	iterations, _ := strconv.Atoi(dollaSplit[1])
	salt := dollaSplit[2]

	hashedInput := hashPassword(rawpass, salt, iterations, 32)

	return hashedInput == user.Password, nil
}

// Creates a random base64 string of the specified length
func genSalt(length int) string {
	rb := make([]byte, length)
	_, err := rand.Read(rb)

	if err != nil {
		fmt.Println(err)
	}

	rs := base64.URLEncoding.EncodeToString(rb)
	return rs
}

// Hashes a password with a new generated salt and the default settings
func hashPasswordDefault(pass string) string {
	salt := genSalt(16)
	return hashPassword(pass, salt, 12000, 32)
}

func hashPassword(pass string, salt string, iterations int, length int) string {
	encpass := pbkdf2.Key([]byte(pass), []byte(salt), iterations, length, sha256.New)
	hashstr := base64.StdEncoding.EncodeToString(encpass)
	// From django docs: <algorithm>$<iterations>$<salt>$<hash>
	return fmt.Sprintf("%s$%s$%s$%s", "pbkdf2_sha256", strconv.Itoa(iterations), salt, hashstr)
}
