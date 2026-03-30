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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Format request tidak valid")
		return
	}

	amount, err := strconv.ParseFloat(req.Amount.Value, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Format amount tidak valid")
		return
	}

	if amount <= 0 {
		respondError(w, http.StatusBadRequest, "Jumlah pembayaran harus lebih dari 0")
		return
	}

	refNo := "REF" + strconv.FormatInt(time.Now().Unix(), 10)

	_, err = database.DB.Exec(`
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
		respondError(w, http.StatusInternalServerError, err.Error())
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