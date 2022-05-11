package repository

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/google/uuid"
	"go-api-currency/models"
	"go-api-currency/utils/constants"
	"google.golang.org/api/option"
	"log"
	"time"
)

type FirestoreRepository interface {
	Save(currency *models.Currency) (*models.Currency, error)
	DeleteOldCurrencies() error
	FindAll() ([]models.Currency, error)
}

type firestoreRepo struct{}

func NewFirestoreRepository() FirestoreRepository {
	return &firestoreRepo{}
}

func (*firestoreRepo) Save(currency *models.Currency) (*models.Currency, error) {
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
	
	currency.ID = uuid.NewString()
	
	_, _, err = client.Collection(constants.Collection).Add(ctx, map[string]interface{}{
		"id":        currency.ID,
		"type":      currency.Type,
		"buyPrice":  currency.BuyPrice,
		"sellPrice": currency.SellPrice,
		"date":      currency.Date,
	})
	if err != nil {
		log.Printf("failed saving the currency: %v", err.Error())
		return nil, err
	}
	return currency, nil
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

func (*firestoreRepo) DeleteOldCurrencies() error {
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
	
	query := client.Collection(constants.Collection).Where("date", "<=", time.Now().Add(-time.Minute*30))
	for {
		doc, err := query.Documents(ctx).Next()
		if err != nil {
			log.Printf("currencies deleted")
			break
		}
		date := doc.Data()["date"].(time.Time)
		id := doc.Data()["id"].(string)
		
		print(date.String())
		print(id)
		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			log.Printf("currency cannot be deleted")
			continue
		}
	}
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
