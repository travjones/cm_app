package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/unrolled/render"

	_ "github.com/lib/pq"
)

type ctx struct {
	db *sqlx.DB
	r  *render.Render
}

type View struct {
	Title  string
	Script string
	Name   string
	Error  string
	Data   interface{}
}

func main() {
	ren := render.New(render.Options{
		Layout:        "shared/layout",
		IsDevelopment: true,
	})

	//sessions
	store := cookiestore.New([]byte("secret"))

	//db
	// db, err := sqlx.Connect("postgres", os.Getenv("CONNECTIONSTRING"))
	db, err := sqlx.Connect("postgres", "user=travisjones dbname=cm_app sslmode=disable")
	if err != nil {
		log.Println("Error initializing database...")
		log.Fatalln(err)
	}

	c := ctx{db, ren}

	n := negroni.New()
	nauth := negroni.New()

	//routers
	router := mux.NewRouter()
	authRouter := mux.NewRouter()

	//AUTHORIZED ROUTES
	authRouter.HandleFunc("/", c.Home)
	router.Handle("/", nauth).Methods("GET")

	//OPEN ROUTES
	router.HandleFunc("/account/login", c.Login).Methods("GET")
	router.HandleFunc("/account/login", c.LoginPost).Methods("POST")
	router.HandleFunc("/account/signup", c.Signup).Methods("GET")
	router.HandleFunc("/account/signup", c.SignupPost).Methods("POST")
	router.HandleFunc("/account/logout", c.Logout).Methods("GET")

	//Middleware
	nauth.Use(negroni.HandlerFunc(RequireAuth))
	nauth.UseHandler(authRouter)

	//Sessions
	n.Use(sessions.Sessions("authex", store))

	n.Use(negroni.NewStatic(http.Dir("public")))

	n.UseHandler(router)

	n.Run(
		fmt.Sprint(":", os.Getenv("PORT")),
	)

}
