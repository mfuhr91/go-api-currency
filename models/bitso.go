package models

type BitsoResponse struct {
	Success bool
	Payload []Book
}

type Book struct {
	Book       string
	CreatedAt  string `json:"created_at"`
	Amount     string
	MarkerSide string `json:"maker_side"`
	Price      string
	TID        float64 `json:"tid"`
}
