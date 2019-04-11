package main

import (
	"encoding/json"
	"net/http"
)

const market_url string = "https://api.warframe.market/v1/items/"
const relic_url string = "https://drops.warframestat.us/data/relics"

j:=json.NewDecoder()
jd:=http.Get("google.com")

func get_relics()(){

}