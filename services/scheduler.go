package services

import (
	"github.com/go-co-op/gocron"
	"go-api-currency/repository"
	"log"
	"time"
)

var (
	firestoreRepository = repository.NewFirestoreRepository()
)

type SchedulerService interface {
	SaveCurrenciesTask()
}

type scheduler struct{}

func NewScheduler() SchedulerService {
	return &scheduler{}
}

func (*scheduler) SaveCurrenciesTask() {

	// https://pkg.go.dev/github.com/go-co-op/gocron#section-readme

	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(10).Minutes().Do(saveCurrenciesAutomatically)
	if err != nil {
		return
	}

	s.StartAsync()
}

func saveCurrenciesAutomatically() {

	log.Println("getting the currencies and saving it automatically...")
	NewCurrencyService(firestoreRepository).Save()
	log.Println("currencies saved!")

}
