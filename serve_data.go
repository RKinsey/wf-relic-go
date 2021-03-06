package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

//Usage string for main()
const usage string = "Usage: serve_data [-u] mongourl_file"
//TODO: figure out better way to get mongoURL to FindRelic. Env var?
var MONGOURL string

func serve() {
	router := NewAPIRouter()
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}
func FindRelic(h http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qvar := r.FormValue("qual")
	quality, err := strconv.Atoi(qvar)
	if err!=nil||quality>3{
		quality=0
	}
	log.Println(qvar)
	result:=SendSingleRelic{FillAndCalculate(vars["id"],vars["tier"],quality)}
	marshaledSend,_:=json.Marshal(result)
	h.Header().Set("Content-Type","application/json")
	h.Write(marshaledSend)

}
func AllRelics(h http.ResponseWriter, r *http.Request) {

	qvar := r.FormValue("qual")
	quality, err := strconv.Atoi(qvar)
	if err!=nil||quality>3{
		quality=0
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, _ := mongo.NewClient(options.Client().ApplyURI(MONGOURL))
	client.Connect(ctx)
	rColl := client.Database("warframe").Collection("relics")
	cur,_ := rColl.Find(ctx, bson.D{})


	numRel64,_:=rColl.CountDocuments(ctx,bson.D{})
	numRel:=int(numRel64)

	rc:=make(chan *Relic)
	relics:=make([]*Relic,numRel)

	for cur.Next(ctx){
		var relicInfo struct{
			Tier string `bson:"tier"`
			RelicName string `bson:"relicName"`
		}
		cur.Decode(&(relicInfo))
		go FCChannel(relicInfo.RelicName,relicInfo.Tier,quality,rc)
	}

	for i:=0;i<numRel;i++{
		relics[i] = <-rc
	}
	result:= SendManyRelics{relics}
	marshaledSend,_:=json.Marshal(result)
	h.Header().Set("Content-Type","application/json")
	h.Write(marshaledSend)
}
func FCChannel(id string, tier string, quality int, relChan chan *Relic){
	relChan<-FillAndCalculate(id,tier,quality)
}

func FillAndCalculate(id string, tier string, quality int)(*Relic) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, _ := mongo.NewClient(options.Client().ApplyURI(MONGOURL))
	client.Connect(ctx)
	rColl := client.Database("warframe").Collection("relics")
	cur := rColl.FindOne(ctx, bson.D{{"relicName", id}, {"tier", tier}})

	result:=new(Relic)
	err := cur.Decode(result)
	if err != nil {
		log.Println(err)
	}
	rarityArray := GetProbArray(quality)
	iColl := client.Database("warframe").Collection("items")
	var ev float64
	for i, item := range result.Rewards {
		cur = iColl.FindOne(ctx, bson.D{{"_id", item.ID}})
		chance := rarityArray[item.Rarity]
		//Doesn't decode properly if I use Result.Relic.Rewards[i]
		var rwds struct {
			ItemName string  `bson:"itemName"`
			Avg      float64 `bson:"avg"`
			Vol      int     `bson:"vol"`
		}
		cur.Decode(&rwds)
		result.Rewards[i].ItemName = rwds.ItemName
		result.Rewards[i].AvgPrice = rwds.Avg
		result.Rewards[i].Volume = rwds.Vol
		result.Rewards[i].Chance = chance
		ev += rwds.Avg * chance
	}
	result.RelicEV = ev
	return result
}



func main() {
	updateOnly := flag.Bool("u", false, "Run only a database update and do not serve")
	skipUpdate := flag.Bool("skip-update", false, "Skip database update on startup, still starts a reloader")
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatalln(usage)
	}
	mongoURLBy, err := ioutil.ReadFile(flag.Arg(0))
	mongoURL:=string(mongoURLBy)
	if err != nil {
		log.Fatalln(usage + "\nmongourl_file should be a file with the URL of your mongodb server")
	}
	MONGOURL=mongoURL
	if *updateOnly && len(flag.Args()) == 1 {
		FillRelics(mongoURL)
	} else {
		if !*skipUpdate {
			FillRelics(mongoURL)
		}
		StartReloader()
		serve()

	}

}
