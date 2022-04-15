package models

import (
	"time"
)

type Currency struct {
	ID        string    `json:"id" firestore:"ID"`
	Type      string    `json:"type" firestore:"type"`
	BuyPrice  float64   `json:"buyPrice" firestore:"buyPrice"`
	SellPrice float64   `json:"sellPrice" firestore:"sellPrice"`
	Date      time.Time `json:"date" firestore:"date"`
}

type CurrencyType struct {
	Pair   string
	Type   string
	Listed bool
}

func GetCurrencyTypes() []CurrencyType {
	var currencyTypes []CurrencyType

	currencyTypes = append(currencyTypes, CurrencyType{
		Type:   "bitcoin",
		Pair:   "btc_ars",
		Listed: false,
	})
	currencyTypes = append(currencyTypes, CurrencyType{
		Type:   "ethereum",
		Pair:   "eth_ars",
		Listed: false,
	})
	currencyTypes = append(currencyTypes, CurrencyType{
		Type:   "tether",
		Pair:   "tusd_btc",
		Listed: false,
	})
	currencyTypes = append(currencyTypes, CurrencyType{
		Type:   "euroBlue",
		Pair:   "euro_ars",
		Listed: false,
	})
	currencyTypes = append(currencyTypes, CurrencyType{
		Type:   "dolarOficial",
		Pair:   "usdo_ars",
		Listed: false,
	})
	currencyTypes = append(currencyTypes, CurrencyType{
		Type:   "dolarBlue",
		Pair:   "usdb_ars",
		Listed: false,
	})
	currencyTypes = append(currencyTypes, CurrencyType{
		Type:   "dolarCCL",
		Pair:   "usdccl_ars",
		Listed: false,
	})
	currencyTypes = append(currencyTypes, CurrencyType{
		Type:   "dolarMEP",
		Pair:   "dolarmep_ars",
		Listed: false,
	})

	return currencyTypes
}
