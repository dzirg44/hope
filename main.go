package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"os"
)



func init() {
        folder := []string{"data/images","data/images/preview", "date/images/thumbnail"}
        for _, path := range folder {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				os.Mkdir(path, 0750)
				if err != nil {
					panic(err)
				}
			}
		}

	db, err := MysqlDB("gophr:gophr@tcp(192.168.1.66:3306)/gophr")
	if err != nil {
		panic(err)
	}
	DB = db



	dbs, err := MysqlDB("gophr:gophr@tcp(192.168.1.66:3306)/gophr")
	if err != nil {
		panic(err)
	}
	DBS = dbs

	globalImageStore = NewDBImageStore()
	globalSessionStoreMysql =  NewDBSessionStore()
    globalArticleStoreMysql =  NewDBArticleStoreMysql()
	globalUserStore = NewDBUserStore()


}
func NewApp() Middleware {
// Public Router
	router := NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/"))))
	//router.PathPrefix("/im/*filepath").Handler(http.FileServer(http.Dir("./data/images")))
	router.PathPrefix("/im/").Handler(http.StripPrefix("/im/", http.FileServer(http.Dir("./data/images/"))))
	router.HandleFunc("/", HandleHome).Methods("GET")
	router.HandleFunc("/about", HandleAbout).Methods("Get")
	router.HandleFunc("/register", HandleUserNew).Methods("GET")
	router.HandleFunc("/register", HandleUserCreate).Methods("POST")
	router.HandleFunc("/login", HandleSessionNew).Methods("GET")
	router.HandleFunc("/login", HandleSessionCreate).Methods("POST")
	router.HandleFunc("/image/{imageID}", HandleImageShow).Methods("GET")
	router.HandleFunc("/articles/{articleID}", HandleArticleShow).Methods("GET")
	router.HandleFunc("/user/{userID}", HandleUserShow).Methods("GET")
	router.HandleFunc("/allimages",HandleAllImageShow).Methods("GET")
	router.HandleFunc("/contacts", HandleContacts).Methods("GET")
	// Secure Router
	secureRouter := NewRouter()
	secureRouter.HandleFunc("/sign-out", HandleSessionDestroy).Methods("GET")
	secureRouter.HandleFunc("/account", HandleUserEdit).Methods("GET")
	secureRouter.HandleFunc("/account", HandleUserUpdate).Methods("POST")
	secureRouter.HandleFunc("/images/new", HandleImageNew).Methods("GET")
	secureRouter.HandleFunc("/images/new", HandleImageCreate).Methods("POST")
	router.HandleFunc("/image/{imageID}", HandleImageDelete).Methods("DELETE")
	router.HandleFunc("/articles/{articleID}", HandleArticleDelete).Methods("DELETE")
	secureRouter.HandleFunc("/article/new",HandleArticleNew).Methods("GET")
	secureRouter.HandleFunc("/article/new",HandleArticleCreate).Methods("POST")
// API
     router.HandleFunc("/", HandleApiVersion1)

//Middleware
	middleware := Middleware{}
	middleware.Add(router)
	middleware.Add(http.HandlerFunc(RequireLogin))
	middleware.Add(secureRouter)
	return middleware
}
func main() {

	log.Fatal(http.ListenAndServe(":3000", NewApp()))
}
func NewRouter() *mux.Router {
	router := mux.NewRouter()


	router.NotFoundHandler = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	return router
}