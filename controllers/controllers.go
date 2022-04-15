package controllers

import (
	"github.com/gin-gonic/gin"
	"go-api-currency/services"
	"net/http"
)

var (
	currencyService services.CurrencyService
)

type controller struct{}

type Controller interface {
	GetLastCurrencies(c *gin.Context)
	GetAllCurrencies(c *gin.Context)
	AddCurrency(c *gin.Context)
}

func NewController(service services.CurrencyService) Controller {
	currencyService = service
	return &controller{}
}

func (*controller) GetLastCurrencies(c *gin.Context) {

	currencies, err := currencyService.GetLastCurrencies()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(currencies) == 0 {
		c.JSON(http.StatusOK, gin.H{"msg": "There are not any currency"})
		return
	}
	c.JSON(http.StatusOK, currencies)

}

func (*controller) GetAllCurrencies(c *gin.Context) {

	currencies, err := currencyService.FindAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(currencies) == 0 {
		c.JSON(http.StatusOK, gin.H{"msg": "There are not any currency"})
		return
	}
	c.JSON(http.StatusOK, currencies)

}
func (*controller) AddCurrency(c *gin.Context) {

	currencyService.Save()

	c.JSON(http.StatusCreated, gin.H{"msg": "successfully saved!"})

}
