package main

import (
	"app/controllers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	// controller intances...
	UrlsController := controllers.UrlsController{}

	// web routes...
	router.HandleFunc("/{short-url}/{broadcast-id}/{subscriber-email}", UrlsController.RedirectToLongUrl).Methods("GET")

	// api routes...
	router.HandleFunc("/api/short-url", UrlsController.Index).Methods("POST")

	// serve...
	log.Fatal(http.ListenAndServe(":8000", router))
}
