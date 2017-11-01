package main

import (
	"time"
	"net/http"
	"net/url"
)

type Session struct {
	Id string
	SessionName string
	SessionId string
	UserId string
	ExpiryString string
	SessionExpiry time.Time
}

const (
	sessionLenght = 24 * 3 * time.Hour
	sessionCookieName = "GophrSession"
	sessionIDLenght = 20
)
func FindOrCreateSession(w http.ResponseWriter, r *http.Request) *Session {
	session := RequestSession(r)
	if session == nil {
		session = NewSession(w)
	}
	return session
}

func NewSession(w http.ResponseWriter) *Session {
	expiry := time.Now().Add(sessionLenght)

	session := &Session {
		SessionId: GenerateID("sess", sessionIDLenght),
		SessionName: GenerateID("sess", sessionIDLenght),
			SessionExpiry: expiry,
	}
	cookie := http.Cookie{
		Name: sessionCookieName,
		Value: session.SessionId,
		Expires: expiry,
	}
	http.SetCookie(w, &cookie)
	return session
}
func RequestSession(r *http.Request) *Session {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return  nil
	}
	session, err := globalSessionStoreMysql.Find(cookie.Value)
	if err != nil {
		panic(err)
	}
	if session == nil {
		return nil
	}
	if session.Expired() {
		globalSessionStoreMysql.Delete(session)
		return  nil
	}
	return session
}
func RequestUser(r *http.Request) *User {
	session := RequestSession(r)
	if session == nil || session.UserId == "" {
		return nil
	}
	user, err := globalUserStore.Find(session.UserId)
	if err != nil {
		panic(err)
	}
	return user
}
func RequireLogin(w http.ResponseWriter, r *http.Request) {
	if RequestUser(r) != nil {
		return
	}
	query := url.Values{}
	query.Add("next", url.QueryEscape(r.URL.String()))
	http.Redirect(w, r, "/login?"+query.Encode(), http.StatusFound)
}
func (session *Session) Expired() bool {
	return session.SessionExpiry.Before(time.Now())
}