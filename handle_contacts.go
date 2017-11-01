package main

import "net/http"

func HandleContacts(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r , "contacts/contacts", nil)
}
