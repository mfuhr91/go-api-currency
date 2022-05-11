package services

import (
	"github.com/gocolly/colly/v2"
	"go-api-currency/models"
	"go-api-currency/utils/constants"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

func GetWebCurrencies() (models.Currency, error) {
	var currency models.Currency
	
	buyPrice, sellPrice, err := fromParaleloHoyWeb()
	if err != nil {
		log.Printf("Something went wrong with %s,"+
			" getting the euro currency from %s",
			constants.ParaleloHoyWeb, constants.EuroBlueWeb)
		
		buyPrice, sellPrice, err = fromEuroBlueWeb()
		if err != nil {
			log.Printf("Something went wrong with %s, cannot get the euro blue currency",
				constants.EuroBlueWeb)
		}
	}
	
	currency.BuyPrice = math.Round(buyPrice*100) / 100
	currency.SellPrice = math.Round(sellPrice*100) / 100
	
	currency.Type = constants.EuroBlueType
	currency.Date = time.Now().UTC()
	
	log.Printf("currency: %v", currency)
	return currency, err
}

func fromParaleloHoyWeb() (buyPrice float64, sellPrice float64, err error) {
	
	log.Println()
	log.Println("CONNECTING TO PARALELOHOY.COM.AR...")
	
	c := colly.NewCollector()
	
	var prices []string
	c.OnHTML(".tabla tbody tr:last-child", func(e *colly.HTMLElement) {
		e.ForEach("td", func(i int, e *colly.HTMLElement) {
			if strings.Contains(e.Text, "$") {
				prices = append(prices, strings.Replace(e.Text, "$", "", 1))
			}
			
		})
		log.Printf("result: %s", prices)
		
		if len(prices) == 2 {
			buyPrice, err = strconv.ParseFloat(prices[0], 64)
			if err != nil {
				log.Printf("error when parsing the value: %s to float", prices[0])
				return
			}
			
			sellPrice, err = strconv.ParseFloat(prices[1], 64)
			if err != nil {
				log.Printf("error when parsing the value: %s to float", prices[1])
				return
			}
		} else {
			log.Printf("index out of range on prices with length %v", len(prices))
		}
		
	})
	
	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting %s", r.URL)
	})
	
	err = c.Visit(constants.ParaleloHoyWeb)
	if err != nil {
		log.Printf(err.Error())
		return buyPrice, sellPrice, err
	}
	if buyPrice > sellPrice {
		tempValue := sellPrice
		sellPrice = buyPrice
		buyPrice = tempValue
	}
	
	/*
		func (s Sequence) Swap(i, j int) {
		    s[i], s[j] = s[j], s[i]
		}
	*/
	
	return buyPrice, sellPrice, nil
}

func fromEuroBlueWeb() (buyPrice float64, sellPrice float64, err error) {
	
	log.Println()
	log.Println("CONNECTING TO EUROBLUE.COM.AR...")
	
	c := colly.NewCollector()
	
	var prices []string
	c.OnHTML(".elementor-text-editor table tbody tr:last-child", func(e *colly.HTMLElement) {
		if len(prices) == 2 {
			return
		}
		e.ForEach("td", func(i int, e *colly.HTMLElement) {
			prices = append(prices, strings.Replace(e.Text, "$", "", 1))
		})
		log.Printf("result: %s", prices)
		
		if len(prices) == 2 {
			buyPrice, err = strconv.ParseFloat(prices[0], 64)
			if err != nil {
				log.Printf("error when parsing the value: %s to float", prices[0])
				return
			}
			
			sellPrice, err = strconv.ParseFloat(prices[1], 64)
			if err != nil {
				log.Printf("error when parsing the value: %s to float", prices[1])
				return
			}
		} else {
			log.Printf("index out of range on prices with length %v", len(prices))
		}
		
	})
	
	c.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting %s", r.URL)
	})
	
	err = c.Visit(constants.EuroBlueWeb)
	if err != nil {
		log.Printf(err.Error())
		return buyPrice, sellPrice, err
	}
	
	return buyPrice, sellPrice, nil
}
