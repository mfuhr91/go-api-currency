package services

import (
	"go-api-currency/models"
	"go-api-currency/repository"
	"go-api-currency/utils/constants"
	"math"
	"sort"
	"strings"
)

var (
	firestoreRepo = repository.NewFirestoreRepository()
)

type CurrencyService interface {
	Save() (*models.Currency, error)
	FindAll() ([]models.Currency, error)
	GetLastCurrencies() ([]models.Currency, error)
}

type service struct{}

func NewCurrencyService(repository repository.FirestoreRepository) CurrencyService {
	firestoreRepo = repository
	return &service{}
}

func (*service) GetLastCurrencies() ([]models.Currency, error) {
	repo := repository.NewFirestoreRepository()
	currencies, err := NewCurrencyService(repo).FindAll()
	if err != nil {
		return currencies, err
	}

	if len(currencies) == 0 {
		return currencies, err
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Date.After(currencies[j].Date)
	})

	var lastCurrencies []models.Currency

	for _, currency := range currencies {
		for _, currType := range models.GetCurrencyTypes() {

			if strings.EqualFold(currency.Type, currType.Type) {
				lastCurrencies = append(lastCurrencies, currency)
				break
			}
			if len(lastCurrencies) == len(models.GetCurrencyTypes()) {
				return lastCurrencies, nil
			}
		}
	}
	return []models.Currency{}, nil
}

func (*service) FindAll() ([]models.Currency, error) {
	currencies, err := firestoreRepo.FindAll()
	if err != nil {
		return currencies, err
	}

	if len(currencies) == 0 {
		return currencies, err
	}

	return currencies, nil
}

func (*service) Save() (*models.Currency, error) {
	resp, err := SaveCrypto()
	if err != nil {
		return resp, err
	}

	resp, err = SaveEuroBlue()
	if err != nil {
		return resp, err
	}

	err = SaveUSD()
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func SaveCrypto() (*models.Currency, error) {
	var resp *models.Currency

	currencyTypes := models.GetCurrencyTypes()

	var currencyBitcoin models.Currency
	for _, currType := range currencyTypes {

		if currType.Type == constants.EuroBlueType ||
			currType.Type == constants.DolarOficialType ||
			currType.Type == constants.DolarBlueType {
			continue
		}

		currency := GetCrypto(currType.Type)

		if currency.Type == constants.BitcoinType {
			currencyBitcoin.BuyPrice = currency.BuyPrice
			currencyBitcoin.SellPrice = currency.SellPrice
		}

		if currency.Type == constants.TetherType {
			currency.BuyPrice = currencyBitcoin.BuyPrice * currency.BuyPrice
			currency.SellPrice = currencyBitcoin.SellPrice * currency.SellPrice
		}

		currency.BuyPrice = math.Round(currency.BuyPrice*100) / 100
		currency.SellPrice = math.Round(currency.SellPrice*100) / 100

		var err error
		resp, err = firestoreRepo.Save(&currency)
		if err != nil {
			return &currency, err
		}
	}
	return resp, nil
}

func SaveEuroBlue() (*models.Currency, error) {
	var resp *models.Currency

	currency, err := GetWebCurrencies()
	if err != nil {
		return &currency, err
	}

	resp, err = firestoreRepo.Save(&currency)
	return resp, nil
}

func SaveUSD() (err error) {
	currencies, err := GetUSD()
	if err != nil {
		return err
	}

	for _, currency := range currencies {
		_, err = firestoreRepo.Save(&currency)
	}
	return nil
}
