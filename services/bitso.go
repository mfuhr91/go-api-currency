package services

import (
	"encoding/json"
	"go-api-currency/models"
	"go-api-currency/utils/constants"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetCrypto(crypto string) models.Currency {

	log.Println()
	log.Println("CONNECTING TO API.BITSO.COM...")

	url := constants.BitsoUrl + constants.BookFilter

	switch crypto {
	case constants.BitcoinType:
		url = url + constants.BitcoinPair
	case constants.EthereumType:
		url = url + constants.EthereumPair
	case constants.TetherType:
		url = url + constants.TetherPair
	default:
		return models.Currency{}
	}

	resp, err := http.Get(url)
	var currency models.Currency

	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return currency
	}

	var bitsoResponse models.BitsoResponse
	err = json.Unmarshal(responseData, &bitsoResponse)

	if err != nil {
		log.Println(err.Error())
		return currency
	}
	books := bitsoResponse.Payload

	buyPriceAdded := false
	sellPriceAdded := false

	for _, book := range books {

		if book.MarkerSide == "buy" && !buyPriceAdded {
			currency.BuyPrice, err = strconv.ParseFloat(book.Price, 64)
			if err != nil {
				log.Println("error trying to parse the buy price")
			}

			buyPriceAdded = true

		}
		if book.MarkerSide == "sell" && !sellPriceAdded {
			currency.SellPrice, err = strconv.ParseFloat(book.Price, 64)
			if err != nil {
				log.Println("error trying to parse the sell price")
			}

			sellPriceAdded = true
		}

	}
	currency.Type = crypto
	currency.Date = time.Now().UTC()

	log.Println(currency)
	return currency
}
