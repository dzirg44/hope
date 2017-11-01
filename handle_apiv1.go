package main

import "net/http"

func HandleApiVersion1(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r , "about/about", nil)
}
