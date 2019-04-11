package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const market_url string = "https://api.warframe.market/v1/items/"
const relic_url string = "https://drops.warframestat.us/data/relics"

func get_relics() {
	type relic struct {
		Tier, RelicName, State string
		Rewards                struct {
			_Id, ItemName, Rarity string
			Chance                int
		}
	}
	fmt.Println("1")
	resp, err := http.Get(market_url + ".json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("2")
	defer resp.Body.Close()
	body, readerr := ioutil.ReadAll(resp.Body)
	fmt.Println("3")
	if readerr != nil {
		log.Fatal(readerr)
	}
	relics := relic{}
	jsonerr := json.Unmarshal(body, &relics)
	fmt.Println("4")
	if jsonerr != nil {
		log.Fatal(jsonerr)
	}
	fmt.Println("5")
	fmt.Println(relics)
}
func main() {
	get_relics()
}
