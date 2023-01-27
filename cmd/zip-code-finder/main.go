package main

import (
	"net/http"

	handler "github.com/jvictore/ZipCodeFinder/internal/handler"
)

func main() {
	http.HandleFunc("/", handler.SearchCepHandler)
	http.HandleFunc("/add", handler.AddCepHandler)
	http.HandleFunc("/get", handler.UpdateCepHandler)
	http.ListenAndServe(":8080", nil)
}
