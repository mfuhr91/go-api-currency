package models

import "time"

type BluelyticsResp struct {
	Oficial     Valor     `json:"oficial"`
	Blue        Valor     `json:"blue"`
	OficialEuro Valor     `json:"oficial_euro"`
	BlueEuro    Valor     `json:"blue_euro"`
	LastUpdate  time.Time `json:"last_update"`
}

type Valor struct {
	ValueAvg  float64 `json:"value_avg"`
	ValueSell float64 `json:"value_sell"`
	ValueBuy  float64 `json:"value_buy"`
}
