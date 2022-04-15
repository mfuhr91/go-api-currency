package repository

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/google/uuid"
	"go-api-currency/models"
	"log"
	"time"
)

type FirestoreRepository interface {
	Save(currency *models.Currency) (*models.Currency, error)
	FindAll() ([]models.Currency, error)
}

type firestoreRepo struct{}

func NewFirestoreRepository() FirestoreRepository {
	return &firestoreRepo{}
}

const (
	projectId  string = "go-currency-api"
	collection string = "currencies"
)

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

	_, _, err = client.Collection(collection).Add(ctx, map[string]interface{}{
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
	iterator := client.Collection(collection).Documents(ctx)

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

func firestoreConnect() (*firestore.Client, context.Context, error) {

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("failed to create a firestore client: %v", err.Error())
		return nil, nil, err
	}

	return client, ctx, nil
}
