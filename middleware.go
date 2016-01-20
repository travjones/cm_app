package main

import (
	"net/http"

	"github.com/goincremental/negroni-sessions"
)

func RequireAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	s := sessions.GetSession(r)

	if s.Get("user_id") == nil {
		http.Redirect(w, r, "/account/login", http.StatusFound)
		return
	}

	next(w, r)
}
