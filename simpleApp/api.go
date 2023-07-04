package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
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
	a.DB, _ = initDB()
	a.Router = mux.NewRouter()
	a.initalizeRoutes()
}
func (a *App) Run() {
	fmt.Println("Server started and listening on port ", a.Port)
	log.Fatal(http.ListenAndServe(a.Port, addCorsHeaders(a.Router)))
	//log.Fatal(http.ListenAndServeTLS(":4443", "../server.crt", "../server.key", nil))
}

func (a *App) initalizeRoutes() {
	a.Router.HandleFunc("/", getReq).Methods("GET")
	a.Router.HandleFunc("/file", a.uploadFile).Methods("POST")
	a.Router.HandleFunc("/file", a.downloadFile).Methods("GET")
	a.Router.HandleFunc("/file", a.deleteFile).Methods("DELETE")
}

func getReq(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Online\n")
}
func (a *App) uploadFile(w http.ResponseWriter, r *http.Request) {
	//verify if there is 1 file already
	count, _ := getCountrows(a.DB)
	if count != 0 {
		http.Error(w, "Cannot upload picture. One file already exists. Delete the existing file", http.StatusBadRequest)
		return
	}
	// Parse the multipart form file
	err := r.ParseMultipartForm(10 << 20) // 10 MB max form size
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}
	// Get the picture data from the form
	file, handler, err := r.FormFile("picture")
	if err != nil {
		http.Error(w, "Failed to read file from form-data", http.StatusBadRequest)
		return
	}
	defer file.Close()
	// print metadata
	fmt.Printf("Uploaded file %+v\n", handler.Filename)
	fmt.Printf("File size %+v\n", handler.Size)
	fmt.Printf("MIME header%+v\n", handler.Header)
	// Read the file data into a byte slice
	picData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read picture data", http.StatusBadRequest)
		return
	}
	// Create a Pictures struct from the file data
	pic := pictures{pic: picData}
	// Store the picture in the database
	err = pic.addFileSql(a.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Return a success response
	w.WriteHeader(http.StatusOK)
}

func (a *App) downloadFile(w http.ResponseWriter, r *http.Request) {
	pic := pictures{}
	err := pic.getFileSql(a.DB)
	if err != nil {
		http.Error(w, fmt.Sprintf("get picture failed with error: %s", err.Error()), http.StatusInternalServerError)
	}
	// Set the content type header to the appropriate image type
	w.Header().Set("Content-Type", "image/png")

	// Serve the image file to the client
	_, err = w.Write(pic.pic)
	if err != nil {
		http.Error(w, fmt.Sprintf("write picture failed with error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (a *App) deleteFile(w http.ResponseWriter, r *http.Request) {
	//verify if there is 1 file already
	count, _ := getCountrows(a.DB)
	if count == 0 {
		http.Error(w, "Cannot delete picture. There are none uploaded.", http.StatusBadRequest)
		return
	}
	err := deleteRows(a.DB)
	if err != nil {
		http.Error(w, fmt.Sprintf("Delete picture failed with error: %s", err.Error()), http.StatusInternalServerError)
	}
	// Return a success response
	w.WriteHeader(http.StatusOK)
}

func addCorsHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow the GET, POST, and OPTIONS methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		// Allow the Content-Type header
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			// Handle preflight requests
			return
		}

		// Call the next handler
		handler.ServeHTTP(w, r)
	})
}
