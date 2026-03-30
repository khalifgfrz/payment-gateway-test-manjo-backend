package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"payment-gateway-test-manjo/utils"
)

type SignatureRequest struct {
	Type                string `json:"type"`
	PartnerReferenceNo  string `json:"partnerReferenceNo"`
	MerchantID          string `json:"merchantId"`
	OriginalReferenceNo string `json:"originalReferenceNo"`
	AmountValue         string `json:"amountValue"`
}

func GenerateSignatureHandler(w http.ResponseWriter, r *http.Request) {
	var req SignatureRequest
	json.NewDecoder(r.Body).Decode(&req)

	var payload string

	if req.Type == "generateQR" {
		payload = req.PartnerReferenceNo + req.MerchantID + req.AmountValue
	} else {
		payload = req.OriginalReferenceNo + req.AmountValue
	}

	signature := utils.GenerateSignature(os.Getenv("SECRET_KEY"), payload)

	json.NewEncoder(w).Encode(map[string]string{
		"payload":   payload,
		"signature": signature,
	})
}