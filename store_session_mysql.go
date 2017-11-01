package main

import (
	"database/sql"
	"time"
)

var globalSessionStoreMysql SessionStoreMysql

type DBSessionStore struct {
	dbSession *sql.DB
}

func NewDBSessionStore() SessionStoreMysql {
	return &DBSessionStore{
		dbSession: DBS,
	}
}

type SessionStoreMysql interface {
	Find(string) (*Session, error)
	Save(*Session) error
	Delete(*Session) error
}

func (s *DBSessionStore) Find(id string) (*Session, error) {

	row := s.dbSession.QueryRow(`
SELECT id, session_name, session_id, user_id, session_expiry
FROM sessions
WHERE session_id = ?`, id)
	session := Session{}
	err := row.Scan(
		&session.Id,
		&session.SessionName,
		&session.SessionId,
		&session.UserId,
		&session.ExpiryString,
	)

	t := "2006-01-02 15:04:05"
	session.SessionExpiry, _ = time.Parse(t, session.ExpiryString)

	if session.Id == "" {
		return nil, nil
	}
	return &session, err

}

func (store *DBSessionStore) Save(session *Session) error {

	_, err := store.dbSession.Exec(`
REPLACE INTO sessions
( id, session_name, session_id, user_id, session_expiry)
VALUES
(?,?,?,?,?)`,
		session.Id,
		session.SessionName,
		session.SessionId,
		&session.UserId,
		&session.SessionExpiry,
	)
	return err
}

func (store *DBSessionStore) Delete(session *Session) error {

	_, err := store.dbSession.Exec(`
DELETE  FROM sessions WHERE session_id=?`,
		session.SessionId,
	)
	return err
}
