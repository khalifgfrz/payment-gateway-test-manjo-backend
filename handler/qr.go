package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"payment-gateway-test-manjo/database"
)

type QRRequest struct {
	PartnerReferenceNo string `json:"partnerReferenceNo"`
	MerchantID         string `json:"merchantId"`
	TrxID              string `json:"trx_id"`
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
}

func GenerateQR(w http.ResponseWriter, r *http.Request) {
	var req QRRequest
	json.NewDecoder(r.Body).Decode(&req)

	amount, _ := strconv.ParseFloat(req.Amount.Value, 64)

	refNo := "REF" + strconv.FormatInt(time.Now().Unix(), 10)

	_, err := database.DB.Exec(`
	INSERT INTO transactions 
	(merchant_id, amount, trx_id, partner_reference_no, reference_no, status, transaction_date)
	VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		req.MerchantID,
		amount,
		req.TrxID,
		req.PartnerReferenceNo,
		refNo,
		"PENDING",
		time.Now(),
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res := map[string]interface{}{
		"responseCode":       "2004700",
		"responseMessage":    "Successful",
		"referenceNo":        refNo,
		"partnerReferenceNo": req.PartnerReferenceNo,
		"qrContent":          "QR-DUMMY-" + req.PartnerReferenceNo,
	}

	json.NewEncoder(w).Encode(res)
}