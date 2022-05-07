package models

type FirebaseCreds struct {
	Type         string `json:"type"`
	ProjectID    string `json:"project_id"`
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
	ClientID     string `json:"client_id"`
	AuthUri      string `json:"auth_uri"`
	TokenUri     string `json:"token_uri"`
	AuthURL      string `json:"auth_provider_x509_cert_url"`
	ClientURL    string `json:"client_x509_cert_url"`
}
