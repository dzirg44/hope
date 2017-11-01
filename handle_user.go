package main

import (
	"net/http"
	"github.com/gorilla/mux"
)

func HandleUserNew(w http.ResponseWriter, r *http.Request) {
	next := r.URL.Query().Get("next")
	//RenderTemplate(w, r, "users/new", nil)
	RenderTemplate(w, r, "users/new", map[string]interface{}{
		"Next": next,
	})
}

func HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	var err error
	user, errors := NewUser(
		r.FormValue("username"),
		r.FormValue("email"),
		r.FormValue("password"),
	)
	next := r.FormValue("next")
	templateData := map[string]interface{}{
		"User": user,
		"Next": next,
	}

	if errors != nil {
		templateErrors := make([]string, len(errors))
		// loop through the error slice to ensure that all errors are validation errors
		// panic, if otherwise
		for i, err := range errors {

			if !isValidationError(err) {
				panic(err)
			}
			if err != nil {
				templateErrors[i] = err.Error()
			}
		}
		templateData["Errors"] = templateErrors
		RenderTemplate(w, r, "users/new", templateData)
		return
	}


	err = globalUserStore.Save(user)
	if err != nil {
		panic(err)
	}
	session := NewSession(w)
	session.UserId = user.Id
	err = globalSessionStoreMysql.Save(session)
	if err != nil {
		panic(err)
	}
	if next == "" {
		next = "/"
	}
	http.Redirect(w, r, next+"/?flash=User+created", http.StatusFound)
}

func HandleUserEdit(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)
	RenderTemplate(w, r, "users/edit", map[string]interface{}{
		"User": user,
	})
}
func HandleUserUpdate(w http.ResponseWriter, r *http.Request) {
	currentUser := RequestUser(r)
	email := r.FormValue("email")
	currentPassword := r.FormValue("currentPassword")
	newPassword := r.FormValue("newPassword")
	user, err := UpdateUser(currentUser, email, currentPassword, newPassword)
	if err != nil {
		if isValidationError(err) {
			RenderTemplate(w,r,"users/edit", map[string]interface{}{
				"Error": err.Error(),
				"User": user,
			})
			return
		}
		panic(err)
	}
	err = globalUserStore.Save(*currentUser)
	if err != nil {
		panic(err)
	}
	http.Redirect(w,r, "/account?flash=User+updated", http.StatusFound)
}
func HandleUserShow(w http.ResponseWriter, r *http.Request) {
vars := mux.Vars(r)
user, err := globalUserStore.Find(vars["userID"])
if err != nil {
	panic(err)
}
if user == nil {
	http.NotFound(w, r)
}
images, err := globalImageStore.FindAllByUser(user,0)
if err != nil {
	panic(err)
}
RenderTemplate(w,r, "users/show", map[string]interface{}{
	"Images": images,
	"User": user,
})
}