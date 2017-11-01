package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"

)

func HandleArticleNew(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w,r, "articles/new", nil)
}


func HandleArticleCreate(w http.ResponseWriter, r *http.Request) {
	HandleArticleCreatePost(w,r)
}


func HandleArticleCreatePost(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)
	article := NewArticle(user)
	image := NewImage(user)

	article.PageTitle = r.FormValue("title")
	article.PageContent = r.FormValue("content")
	file, headers, err := r.FormFile("file")

/*
	if err != nil {
		panic(err)
	}*/

	if file != nil {
		err = image.CreateFromFile(file, headers)
		if err != nil {
			RenderTemplate(w, r, "articles/new", map[string]interface{}{
				"Error": err,
				"image": image,
			})
			return
		}
		article.PageImageId = image.Location
		defer file.Close()
	} else {
		article.PageImageId = ""
	}
	err = article.CreateArticle()
	//if err != nil {
	//}
	http.Redirect(w,r, "/?flash=Articles+Created+Successfully", http.StatusFound)
}

func HandleArticleShow(w http.ResponseWriter, r *http.Request) {
vars := mux.Vars(r)
article, err := globalArticleStoreMysql.Find(vars["articleID"])
if article == nil {
		http.NotFound(w,r)
		return
}

if err != nil {
	panic(err)
}
//404

user, err := globalUserStore.Find(article.PageGUID)
if err != nil {
	panic(err)
}

if user == nil {
	panic(fmt.Errorf("Could not find user %s", article.PageUserId))
}
RenderTemplate(w,r, "articles/show", map[string]interface{}{
	"Article": article,
	"User": user,
})
}
func HandleArticleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	article, err := globalArticleStoreMysql.Find(vars["articleID"])
	if article == nil {
		http.NotFound(w,r)
		return
	}
	if err != nil {
		panic(err)
	}
    err = article.DeleteArticle()
    if err != nil {
    	panic(err)
	}
}

func HandleArticleCreateFromURL(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)
	image := NewImage(user)
	image.Description = r.FormValue("description")
	err := image.CreateFromURL(r.FormValue("url"))
	if err != nil {
		if isValidationError(err) {
			RenderTemplate(w, r, "articles/new", map[string]interface{}{
				"Error": err,
				"ImageURL": r.FormValue("url"),
				"Image": image,
			})
			return
		}
		panic(err)
	}
	http.Redirect(w, r, "/?flash=Article+Created+Successfully", http.StatusFound)

}