package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-api-currency/controllers"
	"go-api-currency/repository"
	"go-api-currency/services"
	"go-api-currency/utils/config"
	"log"
	"net/http"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ping": "pong"})
}

func home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Welcome": "Welcome to go-api-currency",
		"Endpoints": gin.H{
			"/ping":       "check the api status",
			"/currencies": "get the latest currencies",
			"/all":        "get all saved currencies",
		},
		"Author": "Mariano Fuhr",
		"Email":  "mfuhr91@gmail.com",
	})
}

var (
	firestoreRepository = repository.NewFirestoreRepository()
	currencyService     = services.NewCurrencyService(firestoreRepository)
	controller          = controllers.NewController(currencyService)
	scheduler           = services.NewScheduler()
)

func main() {
	
	config.CreateCredsFile()
	scheduler.SaveCurrenciesTask()
	
	r := gin.Default()
	
	r.Use(cors.Default())
	r.NoRoute(home)
	r.GET("/", home)
	r.GET("/ping", ping)
	
	r.GET("/currencies", controller.GetLastCurrencies)
	
	r.GET("/all", controller.GetAllCurrencies)
	//r.POST("/save", controller.AddCurrency) // only for dev and testing
	
	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Cannot start the server: %v ", err.Error())
		return
	}
	
}
