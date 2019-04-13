package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func serve() {
	router := mux.NewRouter()
	router.HandleFunc("/{class}/{id}.json", findRelic).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
func findRelic(h http.ResponseWriter, r *http.Request) {

}
