package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	//"os"
	//"log"


)

func HandleImageNew(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w,r, "images/new", nil)
}


func HandleImageCreate(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("url") != "" {
		HandleImageCreateFromURL(w, r)
		return
	}
	HandleImageCreateFromFile(w,r)
}

func HandleImageCreateFromURL(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)
	image := NewImage(user)
	image.Description = r.FormValue("description")
	err := image.CreateFromURL(r.FormValue("url"))
	if err != nil {
		if isValidationError(err) {
			RenderTemplate(w, r, "images/new", map[string]interface{}{
				"Error": err,
				"ImageURL": r.FormValue("url"),
				"Image": image,
			})
			return
		}
		panic(err)
	}
	http.Redirect(w, r, "/?flash=Image+Upload+Successfully", http.StatusFound)

}
func HandleAllImageShow(w http.ResponseWriter, r *http.Request) {
	images, err := globalImageStore.FindAll(0)
	if err != nil {
		panic(err)
	}

	RenderTemplate(w, r, "images/home", map[string]interface{}{
		"Images": images,
	})

}

func HandleImageCreateFromFile(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)
	image := NewImage(user)
	image.Description = r.FormValue("description")
	file, headers, err := r.FormFile("file")
	if file == nil {
		RenderTemplate(w,r, "images/new", map[string]interface{}{
			"Error": errNoImage,
			"Image": image,
		})
		return
	}
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = image.CreateFromFile(file, headers)
	if err != nil {
		RenderTemplate(w, r, "images/new", map[string]interface{}{
			"Error": err,
			"image": image,
		})
		return
	}
	http.Redirect(w,r, "/?flash=Image+Uploaded+Successfully", http.StatusFound)
}
func HandleImageShow(w http.ResponseWriter, r *http.Request) {
vars := mux.Vars(r)
image, err := globalImageStore.Find(vars["imageID"])
if image == nil {
		http.NotFound(w,r)
		return
}

if err != nil {
	panic(err)
}
//404

user, err := globalUserStore.Find(image.UserID)
if err != nil {
	panic(err)
}
if user == nil {
	panic(fmt.Errorf("Could not find user %s", image.UserID))
}
RenderTemplate(w,r, "images/show", map[string]interface{}{
	"Image": image,
	"User": user,
})
}

func HandleImageDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	image, err := globalImageStore.Find(vars["imageID"])
	if image == nil {
		http.NotFound(w,r)
		return
	}
	err = image.DeleteImageFromHdd(vars["imageID"])
	if err != nil {
		panic(err)
	}
	/*
	err = globalImageStore.Delete(vars["imageID"])
	if err != nil {
		panic(err)
	}
*/


	http.Redirect(w,r, "/?flash=Image+Uploaded+Successfully", http.StatusFound)
}