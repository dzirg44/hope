package main

import "net/http"

func HandleAbout(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r , "about/about", nil)
}
