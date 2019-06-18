package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
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
		Item_name string `json:"item_name" bson:"itemName"`
		Url_name  string `json:"url_name" bson:"urlName"`
	} `json:"en"`
}
type MarketStats struct {
	Payload struct {
		StatClosed  struct {
			StatArray []struct {
				volume    int `json:"volume"`
				Avg_price float64 `json:"avg_price"`
			} `json:"90days"`
		}`json:"statistics_closed"`
	} `json:"payload"`
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

	itemnames:=make([]string,0)
	for cur.Next(ctx) {
		var result struct{
			Item_name string `bson:"itemName"`
		}
		err:=cur.Decode(&result)
		if err!=nil{
			log.Println(err)
		}
		itemnames=append(itemnames,result.Item_name)
	}
	//cur.Close(ctx)
	//wg.Wait()

	for i,item := range itemnames {
		var average float64
		var volume int
		log.Printf("%d ",i)
		if item !="Forma Blueprint" {
			priceDat := MarketStats{}
			if strings.Contains(item,"Kavasa Prime"){
				 splitItem:=strings.Split(item," ")
				 item="Kavasa Prime Collar "+splitItem[len(splitItem)-1]
			}
			uname := urlMap[item]
			if uname == "" {
				uname = urlMap[item[:len(item)-10]]
			}
			url := marketURL + "/" + uname + "/statistics"
			log.Println(url)
			body := GetBytesFromURL(url)
			err := json.Unmarshal(body, &priceDat)
			if err != nil {
				log.Println(err)
			}
			length := len(priceDat.Payload.StatClosed.StatArray)

			average = priceDat.Payload.StatClosed.StatArray[length-1].Avg_price
			volume = priceDat.Payload.StatClosed.StatArray[length-1].volume
		} else{
			average=11+2./3
			volume = 0
		}
		update:=bson.D{{"$set", bson.D{{"avg",average},{"vol",volume}}}}
		iColl.UpdateOne(ctx,bson.D{{"itemName",item}},update)

	}
	client.Disconnect(ctx)
}
