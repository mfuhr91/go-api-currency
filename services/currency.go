package services

import (
	"go-api-currency/models"
	"go-api-currency/repository"
	"go-api-currency/utils/constants"
	"log"
	"math"
	"sort"
	"strings"
)

var (
	firestoreRepo = repository.NewFirestoreRepository()
)

type CurrencyService interface {
	Save()
	FindAll() ([]models.Currency, error)
	GetLastCurrencies() ([]models.Currency, error)
}

type service struct{}

func NewCurrencyService(repository repository.FirestoreRepository) CurrencyService {
	firestoreRepo = repository
	return &service{}
}

func (*service) GetLastCurrencies() ([]models.Currency, error) {
	currencies, err := firestoreRepo.FindAll()
	if err != nil {
		return currencies, err
	}
	
	if len(currencies) == 0 {
		return currencies, err
	}
	
	var lastCurrencies []models.Currency
	
	var currencyTypes = models.GetCurrencyTypes()
	
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Date.After(currencies[j].Date)
	})
	
	for _, currency := range currencies {
		for i, currType := range currencyTypes {
			if strings.EqualFold(currency.Type, currType.Type) && !currType.Listed {
				lastCurrencies = append(lastCurrencies, currency)
				currencyTypes[i].Listed = true
				break
			}
		}
	}
	return lastCurrencies, nil
}

func (*service) FindAll() ([]models.Currency, error) {
	currencies, err := firestoreRepo.FindAll()
	if err != nil {
		return currencies, err
	}
	
	if len(currencies) == 0 {
		return currencies, err
	}
	
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Date.After(currencies[j].Date)
	})
	
	return currencies, nil
}

func (*service) Save() {
	currencies, _ := getCrypto()
	
	currency, _ := getEuroBlue()
	currencies = append(currencies, currency)
	
	currenciesUSD, _ := getUSD()
	for _, item := range currenciesUSD {
		currencies = append(currencies, item)
	}
	
	result, _ := firestoreRepo.FindAll()
	if len(result) == 0 {
		_, _ = firestoreRepo.Save(currencies)
		log.Println("currencies saved!")
		return
	}
	
	_ = firestoreRepo.UpdateCurrencies(currencies)
	log.Println("currencies updated!")
}

func getCrypto() ([]models.Currency, error) {
	currencyTypes := models.GetCurrencyTypes()
	
	var currencies []models.Currency
	var currencyBitcoin models.Currency
	for _, currType := range currencyTypes {
		
		if currType.Type == constants.EuroBlueType ||
			currType.Type == constants.DolarOficialType ||
			currType.Type == constants.DolarCCLType ||
			currType.Type == constants.DolareMEPType ||
			currType.Type == constants.DolarBlueType {
			continue
		}
		
		currency := GetCrypto(currType.Type)
		
		if currency.Type == constants.BitcoinType {
			currencyBitcoin.BuyPrice = currency.BuyPrice
			currencyBitcoin.SellPrice = currency.SellPrice
		}
		
		if currency.Type == constants.TetherType {
			currency.BuyPrice = currencyBitcoin.BuyPrice / currency.BuyPrice
			currency.SellPrice = currencyBitcoin.SellPrice / currency.SellPrice
		}
		
		currency.BuyPrice = math.Round(currency.BuyPrice*100) / 100
		currency.SellPrice = math.Round(currency.SellPrice*100) / 100
		
		log.Println(currency)
		
		currencies = append(currencies, currency)
	}
	return currencies, nil
}

func getEuroBlue() (models.Currency, error) {
	
	currency, err := GetWebCurrencies()
	if err != nil {
		return currency, err
	}
	
	return currency, nil
}

func getUSD() ([]models.Currency, error) {
	var results []models.Currency
	
	currencies, err := GetUSD()
	if err != nil {
		return results, nil
	}
	
	for _, currency := range currencies {
		results = append(results, currency)
	}
	
	return results, nil
}
