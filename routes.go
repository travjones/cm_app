package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/goincremental/negroni-sessions"
	"golang.org/x/crypto/bcrypt"
)

func (c *ctx) Home(w http.ResponseWriter, r *http.Request) {
	s := sessions.GetSession(r)

	data := struct {
		Email string
	}{
		fmt.Sprintf("%v", s.Get("email")),
	}

	v := View{
		"Home",
		"",
		"home",
		"",
		data,
	}

	c.r.HTML(w, http.StatusOK, "home", v)
}

func (c *ctx) Logout(w http.ResponseWriter, r *http.Request) {
	s := sessions.GetSession(r)
	s.Set("user_id", nil)
	s.Set("email", nil)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (c *ctx) Login(w http.ResponseWriter, r *http.Request) {
	v := View{
		"Log In",
		"",
		"login",
		"",
		struct{}{},
	}

	c.r.HTML(w, http.StatusOK, "account/login", v)
}

func (c *ctx) LoginPost(w http.ResponseWriter, r *http.Request) {
	s := sessions.GetSession(r)

	v := View{
		"Log In",
		"",
		"login",
		"",
		struct{}{},
	}

	email, password := r.FormValue("email"), r.FormValue("password")

	var user_id int
	var dbemail string
	var dbpassword string

	q := "SELECT user_id, email, password FROM users WHERE email=$1"
	err := c.db.QueryRow(q, email).Scan(&user_id, &dbemail, &dbpassword)
	if err != nil {
		log.Println("DB login query failed!")
		c.r.HTML(w, http.StatusOK, "account/login", v)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(password))
	if err != nil {
		v.Error = "Access denied: Incorrect username or password"
		log.Println("Access Denied!")
		c.r.HTML(w, http.StatusOK, "account/login", v)
		return
	} else {
		s.Set("user_id", user_id)
		s.Set("email", dbemail)

		log.Println("Login Successful: " + dbemail)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (c *ctx) Signup(w http.ResponseWriter, r *http.Request) {
	v := View{
		"Sign up for Authex",
		"",
		"signup",
		"",
		struct{}{},
	}

	c.r.HTML(w, http.StatusOK, "account/signup", v)
}

func (c *ctx) SignupPost(w http.ResponseWriter, r *http.Request) {
	v := View{
		"Sign up for Authex",
		"",
		"signup",
		"",
		struct{}{},
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		log.Println("bcrypt password hashing piled out...")
	}

	name, email := r.FormValue("name"), r.FormValue("email")

	q := "INSERT INTO users VALUES (default, $1, $2, $3, now())"

	_, err = c.db.Exec(q, name, email, hashPass)
	if err != nil {
		log.Println("Signup failed.")
		c.r.HTML(w, http.StatusOK, "account/signup", v)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
