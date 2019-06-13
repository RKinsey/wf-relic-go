package main

import "github.com/gorilla/mux"

//NewAPIRouter s
func NewAPIRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/{class}/{id}.json", findRelic).Methods("GET")
	router.HandleFunc("/all", allRelics).Methods("GET")
	return router
}