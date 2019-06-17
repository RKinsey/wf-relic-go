package main

import "github.com/gorilla/mux"

//NewAPIRouter s
func NewAPIRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/{class}/{id:[a-zA-Z][0-9]*}.json", findRelic).Methods("GET")
	router.HandleFunc("/all", allRelics).Methods("GET")
	return router
}
