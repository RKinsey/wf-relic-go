package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const marketURL string = "https://api.warframe.market/v1/items"

//Struct for unmarshalling the market item call
type MarketItems struct {
	Payload struct{
		Items ItemT`json:"items"`
	}
}
type ItemT struct {
	EN []struct{
		Item_name string `json:"item_name" bson:"item_name"`
		Url_name  string `json:"url_name" bson:"url_name"`
	} `json:"en"`
}
type MarketStats struct {
	Payload struct {
		Data DataT `json:"statistics_closed"`
	} `json:"payload"`
}
type DataT struct {
	StatArray []struct {
		volume   int `json:"volume"`
		Avg_price int `json:"avg_price"`
	}`json:"90days"`
}
func findItemNames() map[string]string {
	body := GetBytesFromURL(marketURL)

	fmt.Println()
	items:= MarketItems{}
	//todo: Find better way to do this
	err :=json.Unmarshal(body, &items)
	if err!=nil{
		log.Println(err)
	}
	toret := make(map[string]string)
	for _, tem := range items.Payload.Items.EN {
		toret[tem.Item_name] = tem.Url_name
	}
	return toret
}
func GetPrices(mongoURL string) {
	urlMap := findItemNames()
	client, _ := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Second)
	defer cancel()
	client.Connect(ctx)
	//rColl := client.Database("warframe").Collection("relics")
	iColl := client.Database("warframe").Collection("items")
	cur,_ := iColl.Find(ctx,bson.D{})
	//used := make([]string, 50)
	//wg := new(sync.WaitGroup)

	itemnames:=make([]string,50)
	for cur.Next(ctx) {
		var result struct{
			Item_name string `bson:"item_name"`
		}
		err:=cur.Decode(&result)
		if err!=nil{
			log.Println(err)
		}
		itemnames=append(itemnames,result.Item_name)
	}
	//cur.Close(ctx)
	//wg.Wait()

	for _,item := range itemnames {
		priceDat:= MarketStats{}
		url:=marketURL +urlMap[item]+"/statistics"
		body := GetBytesFromURL(url)
		err:=json.Unmarshal(body, &priceDat)
		if err!=nil{
			log.Println(err)
		}
		length:=len(priceDat.Payload.Data.StatArray)

		average:=priceDat.Payload.Data.StatArray[length-1].Avg_price
		volume:=priceDat.Payload.Data.StatArray[length-1].volume
		update:=bson.D{{"$set", bson.D{{"avg",average},{"vol",volume}}}}
		iColl.UpdateOne(ctx,bson.D{{"itemName",item}},update,options.Update().SetUpsert(true))

	}
	client.Disconnect(ctx)
}
