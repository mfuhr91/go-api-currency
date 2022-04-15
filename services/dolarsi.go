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
	dolarsiRes.ValoresPrincipales.DolarOficial.Venta = strings.Replace(dolarsiRes.ValoresPrincipales.DolarOficial.Venta, ",", ".", 1)
	dolarsiRes.ValoresPrincipales.DolarOficial.Compra = strings.Replace(dolarsiRes.ValoresPrincipales.DolarOficial.Compra, ",", ".", 1)

	sellPrice, _ := strconv.ParseFloat(dolarsiRes.ValoresPrincipales.DolarOficial.Venta, 64)
	buyPrice, _ := strconv.ParseFloat(dolarsiRes.ValoresPrincipales.DolarOficial.Compra, 64)

	currency.Type = constants.DolarOficialType
	currency.SellPrice = math.Round(sellPrice*100) / 100
	currency.BuyPrice = math.Round(buyPrice*100) / 100
	currencies = append(currencies, currency)

	// DOLAR BLUE
	dolarsiRes.ValoresPrincipales.DolarBlue.Venta = strings.Replace(dolarsiRes.ValoresPrincipales.DolarBlue.Venta, ",", ".", 1)
	dolarsiRes.ValoresPrincipales.DolarBlue.Compra = strings.Replace(dolarsiRes.ValoresPrincipales.DolarBlue.Compra, ",", ".", 1)

	sellPrice, _ = strconv.ParseFloat(dolarsiRes.ValoresPrincipales.DolarBlue.Venta, 64)
	buyPrice, _ = strconv.ParseFloat(dolarsiRes.ValoresPrincipales.DolarBlue.Compra, 64)

	currency.Type = constants.DolarBlueType
	currency.SellPrice = math.Round(sellPrice*100) / 100
	currency.BuyPrice = math.Round(buyPrice*100) / 100
	currencies = append(currencies, currency)

	// DOLAR CCL
	dolarsiRes.ValoresPrincipales.DolarCCL.Venta = strings.Replace(dolarsiRes.ValoresPrincipales.DolarCCL.Venta, ",", ".", 1)
	dolarsiRes.ValoresPrincipales.DolarCCL.Compra = strings.Replace(dolarsiRes.ValoresPrincipales.DolarCCL.Compra, ",", ".", 1)

	sellPrice, _ = strconv.ParseFloat(dolarsiRes.ValoresPrincipales.DolarCCL.Venta, 64)
	buyPrice, _ = strconv.ParseFloat(dolarsiRes.ValoresPrincipales.DolarCCL.Compra, 64)

	currency.Type = constants.DolarCCLType
	currency.SellPrice = math.Round(sellPrice*100) / 100
	currency.BuyPrice = math.Round(buyPrice*100) / 100
	currencies = append(currencies, currency)

	//DOLAR MEP - BOLSA
	dolarsiRes.ValoresPrincipales.DolarMEP.Venta = strings.Replace(dolarsiRes.ValoresPrincipales.DolarMEP.Venta, ",", ".", 1)
	dolarsiRes.ValoresPrincipales.DolarMEP.Compra = strings.Replace(dolarsiRes.ValoresPrincipales.DolarMEP.Compra, ",", ".", 1)

	sellPrice, _ = strconv.ParseFloat(dolarsiRes.ValoresPrincipales.DolarMEP.Venta, 64)
	buyPrice, _ = strconv.ParseFloat(dolarsiRes.ValoresPrincipales.DolarMEP.Compra, 64)

	currency.Type = constants.DolareMEPType
	currency.SellPrice = math.Round(sellPrice*100) / 100
	currency.BuyPrice = math.Round(buyPrice*100) / 100
	currencies = append(currencies, currency)

	log.Println(currencies)
	return currencies, nil
}
