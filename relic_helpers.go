package main

import (
	"context"
	"log"
	"time"
)

//Relic grade enum

const (
	Intact      = iota
	Exceptional = iota
	Flawless    = iota
	Radiant     = iota
)

//Enum for rarity
const (
	Common   = iota
	Uncommon = iota
	Rare     = iota
)

//PctRarityToInt converts the fractional chance to drop rarity
func PctRarityToInt(rarity float64) int {
	toRet := -1
	switch rarity {
	case 25.33:
		toRet = Common
	case 11:
		toRet = Uncommon
	case 2:
		toRet = Rare
	}
	return toRet
}
func IntRarityToStr(rarity int) string {
	toRet := ""
	switch rarity {
	case Common:
		toRet = "Common"
	case Uncommon:
		toRet = "Uncommon"
	case Rare:
		toRet = "Rare"
	}
	return toRet
}

//FillRelics is called at startup
func FillRelics(mongoURL string) {
	start_time:=time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	GetRelicAPI(ctx, mongoURL)
	log.Println(time.Since(start_time))
	start_time=time.Now()
	GetPrices(ctx, mongoURL)
	log.Println(time.Since(start_time))
}

//GetProbArray takes an integer from the grade constants
//Recommended to use the relic grade constants for readability purposes
//Returns a 3-element array with the probabilities for that grade in the form [common, uncommon, rare]
//Panics if grade is not within its bounds
func GetProbArray(level int) [3]float64 {
	switch level {
	case Intact:
		return [3]float64{.76 / 3, .11, .02}
	case Exceptional:
		return [3]float64{.7 / 3, .13, .04}
	case Flawless:
		return [3]float64{.2, .17, .06}
	case Radiant:
		return [3]float64{.5 / 3, .2, .10}
	default:
		log.Println("Unexpected relic level: " + string(level))
		return [3]float64{.76 / 3, .11, .02}
	}
}
