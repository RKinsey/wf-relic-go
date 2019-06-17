package main

import (
	"encoding/json"
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
type MarketStats struct {
	Payload struct {
		Data struct {
			StatArray []struct {
				Volume   int `json:"volume"`
				AvgPrice int `json:"avg_price"`
			} `json:"90days"`
		} `json:"statistics_closed"`
	} `json:"payload"`
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
func GetPrices(mongoURL string) {
	/*urlMap := findItemNames()
	client, _ := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Connect(ctx)
	//rColl := client.Database("warframe").Collection("relics")
	iColl := client.Database("warframe").Collection("items")
	cur, _ := iColl.Find(ctx, bson.D{})
	defer cur.Close(ctx)
	//used := make([]string, 50)
	//wg := new(sync.WaitGroup)
	var result bson.M
	for cur.Next(ctx) {
		cur.Decode(&result)

	}
	//wg.Wait()

	for item, url := range urlMap {
		var priceDat MarketStats
		body := GetBytesFromURL(marketURL + "/item/"+url+"/statistics")
		json.Unmarshal(body, &priceDat)

	}
	client.Disconnect(ctx)*/
}
