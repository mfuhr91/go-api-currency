package services

import (
	"encoding/json"
	"fmt"
	"go-api-currency/models"
	"go-api-currency/utils/constants"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetUSDBluelytics() (currencies []models.Currency, err error) {

	log.Println()
	log.Println("CONNECTING TO API.BLUELYTICS.COM.AR...")

	url := constants.ApiBluelyticsUrl

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

	var bluelyticsResp models.BluelyticsResp
	err = json.Unmarshal(responseData, &bluelyticsResp)

	if err != nil {
		log.Println(err.Error())
		return currencies, err
	}

	currency.Date = time.Now()

	currency.Type = constants.DolarOficialType
	currency.SellPrice = bluelyticsResp.Oficial.ValueSell
	currency.BuyPrice = bluelyticsResp.Oficial.ValueBuy
	currencies = append(currencies, currency)

	currency.Type = constants.DolarBlueType
	currency.SellPrice = bluelyticsResp.Blue.ValueSell
	currency.BuyPrice = bluelyticsResp.Blue.ValueBuy
	currencies = append(currencies, currency)

	log.Println(currencies)
	return currencies, nil
}
