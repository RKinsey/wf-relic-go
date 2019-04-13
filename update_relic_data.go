package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"
)

/*StartReloader makes a RelicReloader with its own ticker and quit channel
 *
 */
func StartReloader() chan int {
	mongoURL, err := ioutil.ReadFile(os.Args[1])
	dur, _ := time.ParseDuration("48h")
	ticker := time.NewTicker(dur)
	quit := make(chan int)
	if err != nil {
		log.Fatalln("Argument should be a file with the URL of your mongodb server")
	}

	go RelicReloader(string(mongoURL), ticker, quit)
	return quit
}
func RelicReloader(mongourl string, ticker *time.Ticker, quit chan int) {
	for {
		select {
		case <-ticker.C:
			GetRelicAPI(mongourl, quit)
		case <-quit:
			return
		}
	}
}
