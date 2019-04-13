package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const usage string = "Usage: serve_data [-g] mongourl_file"

func serve() {
	router := NewAPIRouter()
	log.Fatal(http.ListenAndServe(":8000", router))
}
func findRelic(h http.ResponseWriter, r *http.Request) {

}
func allRelics(h http.ResponseWriter, r *http.Request) {

}
func manualFill(mongoURL string) {

	GetRelicAPI(mongoURL)
}
func main() {
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		log.Fatalln(usage)
	}
	mongoURL, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatalln(usage + "\nmongourl_file should be a file with the URL of your mongodb server")
	}
	if os.Args[1] == "-u" && len(os.Args) == 3 {
		manualFill(string(mongoURL))
	} else if len(os.Args) == 2 {
		StartReloader()
		serve()
	}

}
