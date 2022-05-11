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

func deleteOldCurrencies() error {
	err := firestoreRepo.DeleteOldCurrencies()
	if err != nil {
		return err
	}
	return nil
}

func (*service) Save() {
	saveCrypto()
	saveEuroBlue()
	saveUSD()
	deleteOldCurrencies()
}

func saveCrypto() (*models.Currency, error) {
	var resp *models.Currency
	
	currencyTypes := models.GetCurrencyTypes()
	
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
		
		var err error
		resp, err = firestoreRepo.Save(&currency)
		if err != nil {
			return &currency, err
		}
	}
	return resp, nil
}

func saveEuroBlue() (*models.Currency, error) {
	var resp *models.Currency
	
	currency, err := GetWebCurrencies()
	if err != nil {
		return &currency, err
	}
	
	resp, err = firestoreRepo.Save(&currency)
	return resp, nil
}

func saveUSD() (err error) {
	currencies, err := GetUSD()
	if err != nil {
		return err
	}
	
	for _, currency := range currencies {
		_, err = firestoreRepo.Save(&currency)
	}
	return nil
}
