package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const relicURL string = "https://drops.warframestat.us/data/relics.json"

//Relic Struct for pulling from the relic API
type Relic struct {
	Tier      string `json:"tier"`
	RelicName string `json:"relicName"`
	Rewards   []struct {
		ID       string `json:"_id"`
		ItemName string `json:"itemName"`
		Rarity   string `json:"rarity"`
	} `json:"rewards"`
	ID string `json:"_id"`
}

//RelicPage is a struct for the Warframestat Relic API JSON
type RelicPage struct {
	Relics []Relic `json:"relics"`
}

//GetBytesFromURL makes a GET request on a URL and returns the body
func GetBytesFromURL(URL string) []byte {
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, readerr := ioutil.ReadAll(resp.Body)
	if readerr != nil {
		log.Fatal(readerr)
	}
	return body
}

func GetRelicAPI(mongourl string) {
	body := GetBytesFromURL(relicURL)
	relicPage := RelicPage{}

	jsonerr := json.Unmarshal(body, &relicPage)
	if jsonerr != nil {
		log.Fatal(jsonerr)
		//return err
	}
	//client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	client, _ := mongo.NewClient(options.Client().ApplyURI(mongourl))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Connect(ctx)
	rColl := client.Database("warframe").Collection("relics")
	iColl := client.Database("warframe").Collection("items")
	inserted := new(sync.Map)
	wg := new(sync.WaitGroup)
	for i := 0; i < len(relicPage.Relics); i += 4 {
		wg.Add(1)
		handleRelic(ctx, rColl, iColl, &relicPage.Relics[i], inserted, wg)
	}
	wg.Wait()
	client.Disconnect(ctx)
}

//RelicToBSON
func RelicToBSON(relic *Relic) bson.D {
	itemIDs := make([]string, len(relic.Rewards))
	for i, rel := range relic.Rewards {
		itemIDs[i] = rel.ID
	}
	return bson.D{
		{Key: "relicid", Value: relic.ID},
		{Key: "Tier", Value: relic.Tier},
		{Key: "relicName", Value: relic.RelicName},
		{Key: "rewardIDs", Value: itemIDs},
	}
}

func handleRelic(ctx context.Context, relicCollection, itemCollection *mongo.Collection, relic *Relic, inserted *sync.Map, wg *sync.WaitGroup) (err error) {
	r := RelicToBSON(relic)
	relicCollection.InsertOne(ctx, r)
	for _, item := range relic.Rewards {
		_, loaded := inserted.LoadOrStore(item.ID, 1)
		if !loaded {
			itemCollection.InsertOne(ctx, bson.D{{Key: "itemid", Value: item.ID}, {Key: "itemName", Value: item.ItemName}, {Key: "rarity", Value: item.Rarity}})
		}
	}
	wg.Done()
	return nil
}
