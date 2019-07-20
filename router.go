package main

import "github.com/gorilla/mux"

//NewAPIRouter s
func NewAPIRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/{tier}/{id:[a-zA-Z][0-9]*}.json", FindRelic).Methods("GET")
	router.HandleFunc("/all", AllRelics).Methods("GET")
	return router
}
