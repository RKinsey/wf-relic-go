package main

type APIRelic struct{
	FissureTier string	`json:"fissure_tier"`
	RelicName	string	`json:"relic_name"`
	RelicGrade	string	`json:"relic_grade"`
	Rewards []struct{
		ID 			string	`json:"-" bson:"_id"`
		ItemName	string	`json:"reward_name" bson:"itemName"`
		Rarity		string	`json:"rarity"`
		Chance		int		`json:"drop_chance"`
		RarityEnum	int		`bson:"rarity"`
		Price		float64	`bson:""`
	}
}
//Relic Struct for pulling from the relic API
type Relic struct {
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
