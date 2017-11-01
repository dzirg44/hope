package main

import (
	"net/http"
	"testing"
	"net/http/httptest"
	"time"
)

type MockSessionStore struct {
	Session *Session
}
func (store MockSessionStore) Find(string) (*Session, error) {
	return store.Session, nil
}

func (store MockSessionStore) Save(*Session) error {
	return nil
}

func (store MockSessionStore) Delete(*Session) error {
	return nil
}

func TestRequestNewImageUnauthenticated(t *testing.T) {
	request, _ := http.NewRequest("GET", "/images/new", nil)
	recorder := httptest.NewRecorder()
	app := NewApp()
	app.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusFound {
		t.Error("Expected  a  redirect code, but  got", recorder.Code)
	}
	loc := recorder.HeaderMap.Get("Location")
	if loc != "/login?next=%252Fimages%252new"{
		t.Error("Expected Location to redirect to sign, but got", loc)
	}

}


func TestRequestNewImageAuthenticated(t *testing.T){
	oldUserStore := globalUserStore
	defer func() {
		globalUserStore = oldUserStore
	}()
	globalUserStore = &MockUserStore{
		findUser: &User{},
	}
	expiry := time.Now().Add(time.Hour)
	oldSessionStore := globalSessionStore
	defer func() {
	 globalSessionStore = oldSessionStore
	}()
	globalSessionStore = &MockSessionStore{
		Session: &Session{
			ID: "session_123",
			UserID: "user_123",
			Expiry: expiry,

		},
	}
	authCookie := &http.Cookie{
		Name: sessionCookieName,
		Value: "session_123",
		Expires:expiry,
	}
	request, _ := http.NewRequest("GET", "/images/new", nil)
	request.AddCookie(authCookie)
	recorder := httptest.NewRecorder()
	app := NewApp()
	app.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Error("Expected a redirect code, but got", recorder.Code)
	}
}