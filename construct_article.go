package main

import (
	"time"
	"fmt"
)

const ArticleIDLenght = 10

type Article struct {
	PageNumber int
	PageId string // AUTO
	PageUserId string  //OK
	PageGUID string // OK
	PageTitle string //OK
	PageContent string //OK
	PageImageId string //OK
	CreatedAt time.Time
	ArticlePagination
	CommentsId int
	PageTags string
}

type ArticlePagination struct {

	FirstPage int
	PrevPrevPage int
	PrevPage int
	CurrentPage int
	NextPage int
	NextNextPage int
	LastPage int
}
func(article *Article) ShowArticleRoute() string {
	return "/articles/" + article.PageGUID
}

func NewArticle(user *User) *Article {
	return  &Article{
		PageGUID: GenerateID("article", ArticleIDLenght),
		PageUserId: user.Id,
		CreatedAt: time.Now(),
	}
}

func (article *Article) CreateArticle() error  {
err := globalArticleStoreMysql.Save(article)
return err
}

func(article *Article) DeleteArticle() error {
image, err := globalImageStore.FindByLocation(article.PageImageId)
	fmt.Println(image, err)
	if image != nil {
		err = image.DeleteImageFromHdd(image.ID)
		fmt.Println(err)
	}
return globalArticleStoreMysql.Delete(article)
}

func (article *Article) Pagination(allPages int, page  int, articles []Article)  []Article {

	// define  new  variable from const ArticlesPerPage
	curArticlePerPage := ArticlesPerPage
	// calculate  number of the pages
	apr := (allPages-1)/ArticlesPerPage + 1
	// calculate start page articles id
	perPageId := page*ArticlesPerPage - ArticlesPerPage
	// calculate  array length from current page
	articlesLen := len(articles[perPageId:])
	// if array length less  than number articles per page then assigned number of  articles equal to articles array
	if curArticlePerPage > articlesLen {
		curArticlePerPage = articlesLen
	}
	// Find pages for current page
	articles = articles[perPageId:perPageId+curArticlePerPage]
	// Define pages to struct
	article.CurrentPage = page

	if apr >= page +1  {
		article.NextPage = page + 1
	}
	if apr >= page+2 {
		article.NextNextPage = page + 2
	}
	if page - 1 > 0 {
		article.PrevPage = page - 1
	}
	if page -2 > 0 {
		article.PrevPrevPage = page - 2
	}
	if page - 1 > 0 {
		article.FirstPage = 1
	}
	if page + 1 <= apr {
		article.LastPage = apr
	}
	// Return Articles
	return articles
}