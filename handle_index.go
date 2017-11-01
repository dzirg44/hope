package main

import (
	"net/http"
	"strconv"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	pageId := r.FormValue("page")
	if pageId == "" {
		pageId = "1"
	}
	k, err := strconv.Atoi(pageId)
	if err != nil {
		panic(err)
	}
	articles, pagination, err := globalArticleStoreMysql.FindAll(0, k)
	if err != nil {
		panic(err)
	}

	RenderTemplate(w, r, "index/home", map[string]interface{}{
		"Articles": articles,
		"Pagination": pagination,
	})
}
/*
func HandleHome2(w http.ResponseWriter, r *http.Request) {
	articles, err := globalArticleStoreMysql.FindAll(0)
	if err != nil {
		panic(err)
	}
	fmt.Println(articles)


	RenderTemplate(w, r, "index/home", map[string]interface{}{
		"Articles": articles,
	})
}
*/