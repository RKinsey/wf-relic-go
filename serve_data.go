package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//Usage string for main()
const usage string = "Usage: serve_data [-u] mongourl_file"

func serve() {
	router := NewAPIRouter()
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}
func findRelic(h http.ResponseWriter, r *http.Request) {
	var b strings.Builder
	vars := mux.Vars(r)

	log.Println(vars["class"])
	for k, v := range mux.Vars(r) {
		fmt.Fprintf(&b, "%s: %s\n", k, v)
	}
	h.Write([]byte(b.String()))
}
func allRelics(h http.ResponseWriter, r *http.Request) {
	h.Write([]byte("Sorry, this isn't implemented yet"))
}

func main() {
	updateOnly := flag.Bool("u", false, "Run only a database update and do not serve")
	skipUpdate := flag.Bool("skip-update", false, "Skip database update on startup, still starts a reloader")
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatalln(usage)
	}
	mongoURL, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatalln(usage + "\nmongourl_file should be a file with the URL of your mongodb server")
	}
	if *updateOnly && len(flag.Args()) == 1 {
		FillRelics(string(mongoURL))
	} else {
		if !*skipUpdate {
			FillRelics(string(mongoURL))
		}
		StartReloader()
		serve()

	}

}
