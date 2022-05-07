package config

import (
	"encoding/json"
	"go-api-currency/models"
	"go-api-currency/utils/constants"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func CreateCredsFile() {
	fileName := constants.JsonPath
	if _, err := os.Stat(fileName); err == nil {
		log.Printf("File %s already exists", fileName)
		return
	}
	
	creds := models.FirebaseCreds{
		Type:         os.Getenv("FSTORE_TYPE"),
		ProjectID:    os.Getenv("FSTORE_PROJECT_ID"),
		PrivateKeyID: os.Getenv("FSTORE_PRIVATE_KEY_ID"),
		PrivateKey:   os.Getenv("FSTORE_PRIVATE_KEY"),
		ClientEmail:  os.Getenv("FSTORE_CLIENT_EMAIL"),
		ClientID:     os.Getenv("FSTORE_CLIENT_ID"),
		AuthUri:      os.Getenv("FSTORE_AUTH_URI"),
		TokenUri:     os.Getenv("FSTORE_TOKEN_URI"),
		AuthURL:      os.Getenv("FSTORE_AUTH_URL"),
		ClientURL:    os.Getenv("FSTORE_CLIENT_URL"),
	}
	
	privateKey := creds.PrivateKey
	creds.PrivateKey = strings.ReplaceAll(privateKey, "\\n", "\n")
	
	jsonCreds, err := json.Marshal(creds)
	if err != nil {
		log.Fatalf("error when marshaling firebaseCreds - error: ", err)
	}
	err = ioutil.WriteFile(constants.JsonPath, jsonCreds, 0644)
	if err != nil {
		log.Fatalf("error when writing firebaseCreds - error: ", err)
	}
}
