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
	router.HandleFunc("/tests/reports/{broadcast-id}/{email}", func(w http.ResponseWriter, r *http.Request) {
		log.Println(mux.Vars(r))
		que := mux.Vars(r)
		if "123" == que["broadcast-id"] {
			w.WriteHeader(200)
		}
		w.WriteHeader(404)
	}).Methods("POST")

	// api routes...
	router.HandleFunc("/api/short-url", UrlsController.Index).Methods("POST")

	// serve...
	log.Fatal(http.ListenAndServe(":8000", router))
}
