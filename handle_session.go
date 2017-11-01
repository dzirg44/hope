package main

import (
	"net/http"
)

func HandleSessionDestroy(w http.ResponseWriter, r *http.Request) {
	session := RequestSession(r)
	if session != nil {
		err := globalSessionStoreMysql.Delete(session)
		if err != nil {
			panic(err)
		}
	}
	RenderTemplate(w, r, "sessions/destroy", nil)
}
func HandleSessionNew(w http.ResponseWriter, r *http.Request) {
	next := r.URL.Query().Get("next")
	RenderTemplate(w, r, "sessions/new", map[string]interface{}{
		"Next": next,
	})
}
func HandleSessionCreate(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	next := r.FormValue("next")
	user, err := FindUser(username, password)
	if err != nil {
		if isValidationError(err) {
			RenderTemplate(w, r, "sessions/new", map[string]interface{}{
				"Error": err,
				"User": user,
				"Next": next,
			})
			return
		}
		panic(err)
	}
	session := FindOrCreateSession(w, r)
	session.UserId = user.Id
	err = globalSessionStoreMysql.Save(session)
	if err != nil {
		panic(err)
	}
	if next == "" {
		next = "/"
	}
	http.Redirect(w, r, next+"?flash=Signed+in", http.StatusFound)
}