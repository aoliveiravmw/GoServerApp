package main

import "net/http"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request)) {
		http.ServeFile(w, r, "index.php")	
	}
	http.ListenAndServe(":8080", nil)
}