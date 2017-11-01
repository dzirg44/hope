package main

import (
	"database/sql"
)

var globalUserStore UserStore


type  UserStore interface {
	Find(string) (*User, error)
	FindByEmail(string) (*User, error)
	FindByUsername(string) (*User, error)
	Save(User) error
}

func NewDBUserStore() UserStore {
	return &MysqlUserStore{
		db: DB,
	}
}
type MysqlUserStore struct {
	db *sql.DB
}


func (store MysqlUserStore) Save(user User) error {
	_, err := store.db.Exec(`
INSERT INTO users
SET user_id=?, user_email=?, user_password=?, user_name=?, user_level=0`,
	user.Id,
		user.Email,
			user.HashedPassword,
				user.Username,
	)
	if err != nil {
		return err
	}
	return nil
}
func (store MysqlUserStore) Find(id string) (*User, error) {
	userdb, err := store.db.Query(`
SELECT user_id, user_email, user_password, user_name
FROM users
WHERE user_id = ?
`, id)
    username := User{}
	for userdb.Next() {
		err = userdb.Scan(
			&username.Id,
			&username.Email,
			&username.HashedPassword,
			&username.Username,
		)
	}
	if  err != nil {
		return nil, nil
	}
	return &username, nil

}

func (store MysqlUserStore) FindByUsername(username string) (*User, error) {
	if username == "" {
		return nil, nil
	}
	userdb, err := store.db.Query(`
SELECT user_id, user_email, user_password, user_name
FROM users
WHERE user_name = ?
`, username)
	usernameDb := User{}
	defer userdb.Close()
	for userdb.Next() {
		err =  userdb.Scan(
			&usernameDb.Id,
			&usernameDb.Email,
			&usernameDb.HashedPassword,
			&usernameDb.Username,)
	}
	if  err != nil {
		return nil, nil
	}
	return &usernameDb, nil
}

func (store MysqlUserStore) FindByEmail(email string) (*User, error) {
	if email == "" {
		return nil, nil
	}
	userdb, err := store.db.Query(`
SELECT user_id, user_email, user_password, user_name
FROM users
WHERE user_name = ?
`, email)
	usermailDb := User{}
	for userdb.Next() {
		err = userdb.Scan(
			&usermailDb.Id,
			&usermailDb.Email,
			&usermailDb.HashedPassword,
			&usermailDb.Username,
		)
	}
	if  err != nil {
		return nil, nil
	}
	return &usermailDb, nil
}
