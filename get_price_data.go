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
)

const marketURL string = "https://api.warframe.market/v1/items"

func findItemNames() map[string]string {
	body := GetBytesFromURL(marketURL)

	fmt.Println()
	items := MarketItems{}
	//todo: Find better way to do this
	err := json.Unmarshal(body, &items)
	if err != nil {
		log.Println(err)
	}
	toret := make(map[string]string)
	for _, tem := range items.Payload.Items {
		toret[tem.Item_name] = tem.Url_name
	}
	return toret
}
func GetPrices(ctx context.Context, mongoURL string) {
	urlMap := findItemNames()
	client, _ := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	client.Connect(ctx)
	//rColl := client.Database("warframe").Collection("relics")
	iColl := client.Database("warframe").Collection("items")
	cur, _ := iColl.Find(ctx, bson.D{})
	//used := make([]string, 50)
	//wg := new(sync.WaitGroup)

	itemnames := make([]string, 0)
	for cur.Next(ctx) {
		var result struct {
			Item_name string `bson:"itemName"`
		}
		err := cur.Decode(&result)
		if err != nil {
			log.Println(err)
		}
		itemnames = append(itemnames, result.Item_name)
	}
	//cur.Close(ctx)
	//wg.Wait()

	for _, item := range itemnames {
		var average float64
		var volume int

		if item != "Forma Blueprint" {
			priceDat := MarketStats{}
			if strings.Contains(item, "Kavasa Prime") {
				splitItem := strings.Split(item, " ")
				item = "Kavasa Prime Collar " + splitItem[len(splitItem)-1]
			}
			uname := urlMap[item]
			if uname == "" {
				uname = urlMap[item[:len(item)-10]]
			}
			url := marketURL + "/" + uname + "/statistics"
			body := GetBytesFromURL(url)
			err := json.Unmarshal(body, &priceDat)
			if err != nil {
				log.Println(err)
			}
			length := len(priceDat.Payload.StatClosed.StatArray)

			average = priceDat.Payload.StatClosed.StatArray[length-1].Avg_price
			volume = priceDat.Payload.StatClosed.StatArray[length-1].Volume
		} else {
			average = 11 + 2./3
			volume = 0
		}
		update := bson.D{{"$set", bson.D{{"avg", average}, {"vol", volume}}}}
		iColl.UpdateOne(ctx, bson.D{{"itemName", item}}, update)

	}
	client.Disconnect(ctx)
}
