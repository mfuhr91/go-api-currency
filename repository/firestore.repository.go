package repository

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/google/uuid"
	"go-api-currency/models"
	"go-api-currency/utils/constants"
	"google.golang.org/api/option"
	"log"
	"strings"
	"time"
)

type FirestoreRepository interface {
	Save(currencies []models.Currency) ([]models.Currency, error)
	UpdateCurrencies([]models.Currency) error
	FindAll() ([]models.Currency, error)
}

type firestoreRepo struct{}

func NewFirestoreRepository() FirestoreRepository {
	return &firestoreRepo{}
}

func (*firestoreRepo) Save(currencies []models.Currency) ([]models.Currency, error) {
	client, ctx, err := firestoreConnect()
	if err != nil {
		return nil, err
	}
	
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			return
		}
	}(client)
	
	var curreniesNotSaved []models.Currency
	for _, item := range currencies {
		
		item.ID = uuid.NewString()
		
		_, _, err = client.Collection(constants.Collection).Add(ctx, map[string]interface{}{
			"id":        item.ID,
			"type":      item.Type,
			"buyPrice":  item.BuyPrice,
			"sellPrice": item.SellPrice,
			"date":      item.Date,
		})
		if err != nil {
			curreniesNotSaved = append(curreniesNotSaved, item)
		}
	}
	if len(curreniesNotSaved) != 0 {
		for _, item := range curreniesNotSaved {
			log.Printf("failed saving the currency: %v", item)
		}
	}
	log.Printf("all currencies saved")
	return currencies, nil
}

func (*firestoreRepo) FindAll() ([]models.Currency, error) {
	client, ctx, err := firestoreConnect()
	if err != nil {
		return nil, err
	}
	
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			return
		}
	}(client)
	
	var currencies []models.Currency
	iterator := client.Collection(constants.Collection).Documents(ctx)
	
	for {
		doc, err := iterator.Next()
		
		if len(currencies) == 0 && err != nil {
			log.Printf("failed to iterate the list of currencies: %v", err.Error())
			break
		}
		if err != nil {
			log.Printf("all currencies listed")
			break
		}
		currency := models.Currency{
			ID:        doc.Data()["id"].(string),
			Type:      doc.Data()["type"].(string),
			BuyPrice:  doc.Data()["buyPrice"].(float64),
			SellPrice: doc.Data()["sellPrice"].(float64),
			Date:      doc.Data()["date"].(time.Time),
		}
		currencies = append(currencies, currency)
	}
	return currencies, nil
}

func (*firestoreRepo) UpdateCurrencies(currencies []models.Currency) error {
	client, ctx, err := firestoreConnect()
	if err != nil {
		return err
	}
	
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			return
		}
	}(client)
	
	var results []time.Time
	
	iterator := client.Collection(constants.Collection).Documents(ctx)
	
	for {
		doc, err := iterator.Next()
		
		if err != nil {
			break
		}
		var currency models.Currency
		for _, item := range currencies {
			if strings.EqualFold(item.Type, doc.Data()["type"].(string)) {
				currency = item
				break
			}
		}
		result, err := doc.Ref.Update(ctx, []firestore.Update{
			{
				Path:  "date",
				Value: currency.Date,
			},
			{
				Path:  "buyPrice",
				Value: currency.BuyPrice,
			},
			{
				Path:  "sellPrice",
				Value: currency.SellPrice,
			},
		})
		if err != nil {
			return err
		}
		results = append(results, result.UpdateTime)
	}
	
	if len(results) == 0 {
		log.Printf("currencies not updated")
	}
	
	log.Printf("all currencies updated")
	
	return nil
}

func firestoreConnect() (*firestore.Client, context.Context, error) {
	
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, constants.ProjectId, option.WithCredentialsFile(constants.JsonPath))
	if err != nil {
		log.Fatalf("failed to create a firestore client: %v", err.Error())
		return nil, nil, err
	}
	
	return client, ctx, nil
}
