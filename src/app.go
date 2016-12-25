package main

import (
	"app/controllers"
	"configs"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	// controller instance...
	UrlsController := controllers.UrlsController{}

	// web routes...
	router.HandleFunc("/{short-url}/{broadcast-id}/{subscriber-email}", UrlsController.RedirectToLongUrl).Methods("GET")

	// api routes...
	router.HandleFunc("/api/short-url", UrlsController.Index).Methods("POST")

	// serve...
	log.Fatal(http.ListenAndServe(":"+configs.APP_RUNNING_PORT, router))
}
