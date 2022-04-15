package models

type DolarsiResp struct {
	ValoresPrincipales ValoresPrincipales `xml:"valores_principales"`
}

type ValoresPrincipales struct {
	DolarBlue    Casa `xml:"casa310"`
	DolarCCL     Casa `xml:"casa312"`
	DolarMEP     Casa `xml:"casa313"`
	DolarOficial Casa `xml:"casa349"`
}

type Casa struct {
	Compra  string `xml:"compra"`
	Venta   string `xml:"venta"`
	Agencia int    `xml:"agencia"`
	Nombre  string `xml:"nombre"`
}
