package model

type Transaction struct {
	MerchantID         string
	Amount             float64
	TrxID              string
	PartnerReferenceNo string
	ReferenceNo        string
	Status             string
	TransactionDate    string
	PaidDate           string
}