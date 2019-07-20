package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"
)

//var forcerefreshtime time.Time

/*StartReloader makes a RelicReloader with its own ticker and quit channel
 *quit is currently unused outside of this and is only returned for ease of future modifications
 *Note: this does NOT run when started because of how time.Ticker works
 */
func StartReloader() chan int {
	//TODO: use timestamp checking for scalability (i.e. only one server in a cluster runs the update). Might be better to just use cron job
	mongoURL, err := ioutil.ReadFile(os.Args[1])
	dur, _ := time.ParseDuration("48h")
	ticker := time.NewTicker(dur)
	quit := make(chan int)
	if err != nil {
		log.Fatalln("Argument should be a file with the URL of your mongodb server")
	}
	go relicReloader(string(mongoURL), ticker, quit)
	return quit
}

func relicReloader(mongoURL string, ticker *time.Ticker, quit chan int) {
	for {
		select {
		case <-ticker.C:
			FillRelics(mongoURL)
		case <-quit:
			return
		}
	}
}
