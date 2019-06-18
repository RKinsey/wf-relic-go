package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//relicURL holds the URL to the relic data API
const relicURL string = "https://drops.warframestat.us/data/relics.json"

//Relic Struct for pulling from the relic API
type Relic struct {
	ID string `json:"_id" bson:"_id"`
	Tier      string `json:"tier" bson:"tier"`
	RelicName string `json:"relicName" bson:"relicName"`
	Rewards   []struct {
		ID          string `json:"_id" bson:"itemid"`
		ItemName    string `json:"itemName" bson:"-"`
		RarityFrac  float64 `json:"chance" bson:"-"`
		RarityEnum  int    `bson:"rarity"`
	} `json:"rewards" bson:"rewards"`

}

//RelicPage is a struct for the Warframestat Relic API JSON used to sever out individual relic entries
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

//GetRelicAPI s
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
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	client.Connect(ctx)
	rColl := client.Database("warframe").Collection("relics")
	iColl := client.Database("warframe").Collection("items")
	inserted := new(sync.Map)
	wg := new(sync.WaitGroup)
	//relics:=make(bson.D,len(relicPage.Relics))
	for i := 0; i < len(relicPage.Relics); i += 4 {
		wg.Add(1)
		handleRelic(ctx, rColl, iColl, &relicPage.Relics[i], inserted, wg)
	}
	wg.Wait()
	client.Disconnect(ctx)
}

//handleRelic crunches a a relic struct into BSON form and inserts it into the MongoDB instance
func handleRelic(ctx context.Context, relicCollection, itemCollection *mongo.Collection, relic *Relic, inserted *sync.Map, wg *sync.WaitGroup) {

	for i, item := range relic.Rewards {
		relic.Rewards[i].RarityEnum=PctRarityToInt(item.RarityFrac)
		_, loaded := inserted.LoadOrStore(item.ID, 1)
		if !loaded {
			query := bson.D{{"_id", item.ID}, {"item_name", item.ItemName}}
			ud := bson.D{{"$set", bson.D{{"_id", item.ID}, {"item_name", item.ItemName}}}}
			opt := options.Update()
			opt.SetUpsert(true)
			_, err := itemCollection.UpdateOne(ctx, query, ud, opt)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	query:=bson.D{{"_id",relic.ID},{"tier",relic.Tier},{"relicName",relic.RelicName}}
	//This is extremely hacky but I can't find a way to convert my structs directly to bson.Ds
	umholder:=new(bson.D)
	updata,err:=bson.Marshal(relic)
	if err!=nil{
		log.Println(err)
	}
	err= bson.Unmarshal(updata,umholder)
	if err!=nil{
		log.Println(err)
	}
	update_doc:=bson.D{{"$set",umholder}}
	opt:=options.Update()
	opt.SetUpsert(true)
	_,err=relicCollection.UpdateOne(ctx,query,update_doc,opt)
	if err!=nil{
		log.Println(err)
	}
	wg.Done()
}
