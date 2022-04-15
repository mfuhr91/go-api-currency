package main

import (
	"github.com/gin-gonic/gin"
	"go-api-currency/controllers"
	"go-api-currency/repository"
	"go-api-currency/services"
	"log"
	"net/http"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ping": "pong"})
}

func home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Welcome to currency-api"})
}

var (
	firestoreRepository = repository.NewFirestoreRepository()
	currencyService     = services.NewCurrencyService(firestoreRepository)
	controller          = controllers.NewController(currencyService)
)

func main() {

	r := gin.Default()
	r.GET("/", home)
	r.GET("/ping", ping)

	r.GET("/currencies", controller.GetLastCurrencies)

	r.GET("/all", controller.GetAllCurrencies)
	r.POST("/save", controller.AddCurrency)

	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Cannot start the server: %v ", err.Error())
		return
	}

}
