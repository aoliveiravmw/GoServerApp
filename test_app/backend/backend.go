package backend

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	DB     *sql.DB
	Router *mux.Router
	Port   string
}

func (a *App) Initialize() {

	a.Router = mux.NewRouter()
	a.initalizeRoutes()
}
func (a *App) Run() {
	fmt.Println("Server started and listening on port ", a.Port)
	log.Fatal(http.ListenAndServe(a.Port, a.Router))
	//log.Fatal(http.ListenAndServeTLS(":4443", "../server.crt", "../server.key", nil))
}

func (a *App) initalizeRoutes() {
	a.Router.HandleFunc("/", getReq).Methods("GET")
	a.Router.HandleFunc("/", postReq).Methods("POST")
}

func postReq(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is a post\n")
}

func getReq(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is a get\n")
}
