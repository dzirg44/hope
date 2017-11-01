package main

import (
	"database/sql"
)

var globalArticleStoreMysql ArticleStore

type ArticleStore interface {
	Save(article *Article) error
	Find(id string) (*Article, error)
	Delete(article *Article) error
	FindAll(offset int, page int) ([]Article, Article, error)
	FindAllByUser(user *User, offset int) ([]Article, error)
	//Pagination (i int, page  int, articles []Article) *Article
}

const Articles = 100
const ArticlesPerPage = 5

type DBArticleStore struct {
	db *sql.DB
}

func NewDBArticleStoreMysql() ArticleStore {
	return &DBArticleStore{
		db: DB,
	}
}

func (store *DBArticleStore) Save(article *Article) error {
	_, err := store.db.Exec(
		`
		REPLACE INTO articles
		(page_user_id, page_guid, page_title, page_content, page_image_id, page_date)
		VALUES
		(?,?,?,?,?,?)
		`,
		article.PageUserId,
		article.PageGUID,
		article.PageTitle,
		article.PageContent,
		article.PageImageId,
		article.CreatedAt,
	)
	return err
}

func (store *DBArticleStore) Find(id string) (*Article, error) {
	row := store.db.QueryRow(
		`
SELECT id, page_user_id, page_guid, page_title, page_content, page_image_id, page_date
FROM articles
WHERE page_guid = ?`,
		id,
	)
	article := Article{}
	err := row.Scan(
		&article.PageId,
		&article.PageUserId,
		&article.PageGUID,
		&article.PageTitle,
		&article.PageContent,
		&article.PageImageId,
		&article.CreatedAt,
	)
	if article.PageGUID == "" {
		return nil, nil
	}
	return &article, err

}

func (store *DBArticleStore) FindAll(offset int, page  int) ([]Article,Article, error) {
	pagination := Article{}
	articles := []Article{}
	rows, err := store.db.Query(
		`
SELECT  id, page_user_id, page_guid, page_title, page_content, page_image_id, page_date
FROM articles
ORDER BY page_date DESC
LIMIT ?`,
		Articles,
	)
	if err != nil {
		return nil, pagination, err
	}
	var i int
	for rows.Next() {
		i++
		article := Article{}
		err := rows.Scan(
			&article.PageId,
			&article.PageUserId,
			&article.PageGUID,
			&article.PageTitle,
			&article.PageContent,
			&article.PageImageId,
			&article.CreatedAt,
		)
		article.PageNumber = i
		if err != nil {
			return nil, pagination, err
		}
		articles = append(articles, article)

	}
	articles = pagination.Pagination(i,page, articles)

	return articles, pagination, nil
}



func (store *DBArticleStore) FindAllByUser(user *User, offset int) ([]Article, error) {
	rows, err := store.db.Query(
		`
SELECT id, page_user_id, page_guid, page_title, page_content, page_image_id, page_date
FROM images
WHERE page_user_id = ?
ORDER BY created_at DESC
LIMIT ?
OFFSET ?`,
		user.Id,
		Articles,
		offset,
	)
	if err != nil {
		return nil, err
	}
	articles := []Article{}
	for rows.Next() {
		article := Article{}
		err := rows.Scan(
			&article.PageId,
			&article.PageUserId,
			&article.PageGUID,
			&article.PageTitle,
			&article.PageContent,
			&article.PageImageId,
			&article.CreatedAt,

		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func (store *DBArticleStore) Delete(article *Article) error {
	_, err := store.db.Exec(
		`
		DELETE
		FROM articles
		WHERE page_guid = ?
		`, article.PageGUID)
	return err
}

