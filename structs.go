package main

//RelicPage is a struct for the Warframestat Relic API JSON used to sever out individual relic entries
type RelicPage struct {
	Relics []APIRelic `json:"relics"`
}
type SendSingleRelic struct {
	ToSend *Relic `json:"relic"`
}
type SendManyRelics struct{
	ToSend []Relic `json:"relics"`
}
type Relic struct{
		//Tier of the Relic
		Tier string `json:"tier" bson:"tier"`
		//Name is the two character identifier (e.g. A2)
		Name string `json:"relicName" bson:"relicName"`
		//Relic's expected value
		RelicEV float64 `json:"relicEV" bson:"-"`
		Rewards []struct {
			ID       string  `json:"-" bson:"itemid"`
			ItemName string  `json:"itemName" bson:"itemName"`
			Rarity   int     `json:"-" bson:"rarity"`
			Chance   float64 `json:"dropChance"`
			Volume   int `json:"recentVolume" bson:"vol"`
			AvgPrice float64 `json:"avgPrice" bson:"avg"`
		} `json:"rewards" bson:"rewards"`
}

//Relic Struct for pulling from the relic API
type APIRelic struct {
	ID        string `json:"_id" bson:"_id"`
	Tier      string `json:"tier" bson:"tier"`
	RelicName string `json:"relicName" bson:"relicName"`
	Rewards   []struct {
		ID         string  `json:"_id" bson:"itemid"`
		ItemName   string  `json:"itemName" bson:"-"`
		RarityFrac float64 `json:"chance" bson:"-"`
		RarityEnum int     `bson:"rarity"`
	} `json:"rewards" bson:"rewards"`
}

//Struct for unmarshalling the market item call
/*type MarketItems struct {
	Payload struct {
		Items struct {
			EN []struct {
				Item_name string `json:"item_name" bson:"itemName"`
				Url_name  string `json:"url_name" bson:"urlName"`
			} `json:"en"`
		} `json:"items"`
	}`json:"payload"`
}*/
type MarketItems struct {
	Payload struct {
		Items []struct {
				Item_name string `json:"item_name" bson:"itemName"`
				Url_name  string `json:"url_name" bson:"urlName"`
		} `json:"items"`
	}`json:"payload"`
}

//type ItemT
type MarketStats struct {
	Payload struct {
		StatClosed struct {
			StatArray []struct {
				Volume    int     `json:"volume"`
				Avg_price float64 `json:"avg_price"`
			} `json:"90days"`
		} `json:"statistics_closed"`
	} `json:"payload"`
}
