package services

import (
	"encoding/xml"
	"fmt"
	"go-api-currency/models"
	"go-api-currency/utils/constants"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetUSD() (currencies []models.Currency, err error) {
	
	log.Println()
	log.Println("CONNECTING TO DOLARSI.COM...")
	
	url := constants.DolarSiUrl
	
	resp, err := http.Get(url)
	var currency models.Currency
	
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return currencies, err
	}
	
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("error getting the currencies from: %s - statusCode: %v",
			constants.ApiBluelyticsUrl, resp.Status)
		log.Printf(err.Error())
		return currencies, err
	}
	
	var dolarsiRes models.DolarsiResp
	err = xml.Unmarshal(responseData, &dolarsiRes)
	
	if err != nil {
		log.Println(err.Error())
		return currencies, err
	}
	
	currency.Date = time.Now().UTC()
	
	// DOLAR OFICIAL
	sellPrice, buyPrice := parseCurrencies(dolarsiRes.ValoresPrincipales.DolarOficial)
	
	currency.Type = constants.DolarOficialType
	currency.SellPrice = math.Round(sellPrice*100) / 100
	currency.BuyPrice = math.Round(buyPrice*100) / 100
	currencies = append(currencies, currency)
	
	// DOLAR BLUE
	sellPrice, buyPrice = parseCurrencies(dolarsiRes.ValoresPrincipales.DolarBlue)
	
	currency.Type = constants.DolarBlueType
	currency.SellPrice = math.Round(sellPrice*100) / 100
	currency.BuyPrice = math.Round(buyPrice*100) / 100
	currencies = append(currencies, currency)
	
	// DOLAR CCL
	sellPrice, buyPrice = parseCurrencies(dolarsiRes.ValoresPrincipales.DolarCCL)
	
	currency.Type = constants.DolarCCLType
	currency.SellPrice = math.Round(sellPrice*100) / 100
	currency.BuyPrice = math.Round(buyPrice*100) / 100
	currencies = append(currencies, currency)
	
	//DOLAR MEP - BOLSA
	sellPrice, buyPrice = parseCurrencies(dolarsiRes.ValoresPrincipales.DolarMEP)
	
	currency.Type = constants.DolareMEPType
	currency.SellPrice = math.Round(sellPrice*100) / 100
	currency.BuyPrice = math.Round(buyPrice*100) / 100
	currencies = append(currencies, currency)
	
	log.Println(currencies)
	return currencies, nil
}

func parseCurrencies(dolar models.Casa) (sellPrice float64, buyPrice float64) {
	
	dolarSell := strings.Replace(dolar.Venta, ",", ".", 1)
	dolarBuy := strings.Replace(dolar.Compra, ",", ".", 1)
	
	sellPrice, _ = strconv.ParseFloat(dolarSell, 64)
	buyPrice, _ = strconv.ParseFloat(dolarBuy, 64)
	
	if buyPrice > sellPrice {
		tempValue := sellPrice
		sellPrice = buyPrice
		buyPrice = tempValue
	}
	
	return sellPrice, buyPrice
}

/*
func (s Sequence) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
*/
