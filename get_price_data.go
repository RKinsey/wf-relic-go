package main

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const marketURL string = "https://api.warframe.market/v1/items"

//Struct for unmarshalling the market item call
type MarketItems struct {
	Items struct {
		EN []struct {
			ItemName string `json:"item_name"`
			URLName  string `json:"url_name"`
		} `json:"en"`
	} `json:"items"`
}

func findItemNames() map[string]string {
	body := GetBytesFromURL(marketURL)
	items := MarketItems{}
	json.Unmarshal(body, &items)
	toret := make(map[string]string)
	for _, tem := range items.Items.EN {
		toret[tem.ItemName] = tem.URLName
	}
	return toret
}
func getPrices(mongoURL string) {
	urlMap := findItemNames()
	client, _ := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Connect(ctx)
	rColl := client.Database("warframe").Collection("relics")
	iColl := client.Database("warframe").Collection("items")
	cur, _ := iColl.Find(ctx, bson.D{})
	defer cur.Close(ctx)
	used := make([]string,50)
	var result bson.M
	for cur.Next(ctx) {
		cur.Decode(&result)

	}
	wg.Wait()
	client.Disconnect(ctx)
	for item, url := range urlMap {
		iColl.Fin
	}
	body := GetBytesFromURL(marketURL + "/")
}
