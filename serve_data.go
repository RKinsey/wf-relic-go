package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func serve() {
	router := mux.NewRouter()
	router.HandleFunc("/{class}/{id}.json", findRelic).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
func findRelic(h http.ResponseWriter, r *http.Request) {

}
func main() {
	if os.Args[1] == "-g" && len(os.Args) == 3 {
		manualFill()
	} else if len(os.Args) == 2 {
		StartReloader()
		serve()

	} else {
		log.Fatalln("Usage: serve_data [-g] mongourl_file")
	}

}
