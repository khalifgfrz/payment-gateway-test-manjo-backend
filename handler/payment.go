package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"payment-gateway-test-manjo/database"
)

type PaymentRequest struct {
	OriginalReferenceNo   string `json:"originalReferenceNo"`
	TransactionStatusDesc string `json:"transactionStatusDesc"`
	PaidTime              string `json:"paidTime"`
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": status,
		"message": message,
	})
}

func PaymentCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Format request tidak valid")
		return
	}

	var dbAmount float64
	var currentStatus string
	err := database.DB.QueryRow(`
		SELECT amount, status
		FROM transactions 
		WHERE reference_no = $1
	`, req.OriginalReferenceNo).Scan(&dbAmount, &currentStatus)

	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "Transaksi tidak ditemukan")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if currentStatus != "PENDING" {
        message := "Transaksi dengan nomor referensi " + req.OriginalReferenceNo + " sudah diproses dengan status " + currentStatus
        respondError(w, http.StatusBadRequest, message)
        return
    }

	reqAmount, err := strconv.ParseFloat(req.Amount.Value, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Format jumlah pembayaran tidak valid")
		return
	}

	if reqAmount != dbAmount {
		respondError(w, http.StatusBadRequest, "Jumlah pembayaran tidak sesuai dengan transaksi")
		return
	}

	result, err := database.DB.Exec(`
		UPDATE transactions 
		SET status=$1, paid_date=$2
		WHERE reference_no=$3
	`,
		req.TransactionStatusDesc,
		req.PaidTime,
		req.OriginalReferenceNo,
	)

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		respondError(w, http.StatusNotFound, "Transaksi tidak ditemukan")
		return
	}

	res := map[string]interface{}{
		"responseCode":          "2005100",
		"responseMessage":       "Successful",
		"transactionStatusDesc": req.TransactionStatusDesc,
	}

	json.NewEncoder(w).Encode(res)
}