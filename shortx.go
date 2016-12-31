package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iftekhersunny/shortx/configs"
	"github.com/iftekhersunny/shortx/controllers"
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
