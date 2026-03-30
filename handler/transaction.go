package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Transaction struct {
	ID                 string         `json:"id"`
	MerchantID         string         `json:"merchant_id"`
	Amount             float64        `json:"amount"`
	TrxID              string         `json:"trx_id"`
	PartnerReferenceNo string         `json:"partner_reference_no"`
	ReferenceNo        string         `json:"reference_no"`
	Status             string         `json:"status"`
	TransactionDate    string         `json:"transaction_date"`
	PaidDate           sql.NullString `json:"paid_date"`
}

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type Response struct {
	Data []Transaction `json:"data"`
	Meta Meta          `json:"meta"`
}

func GetTransactions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")
		ref := r.URL.Query().Get("reference_no")
		sort := r.URL.Query().Get("sort")
		order := strings.ToUpper(r.URL.Query().Get("order"))

		page := 1
		limit := 10

		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}

		offset := (page - 1) * limit

		search := "%"
		if ref != "" {
			search = "%" + ref + "%"
		}

		validSortFields := map[string]string{
			"transaction_date": "transaction_date",
			"amount":           "amount",
			"status":           "status",
		}

		sortField, ok := validSortFields[sort]
		if !ok {
			sortField = "transaction_date"
		}

		if order != "ASC" && order != "DESC" {
			order = "DESC"
		}

		var total int
		err := db.QueryRow(`
			SELECT COUNT(*) FROM transactions
			WHERE reference_no ILIKE $1
		`, search).Scan(&total)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		totalPages := (total + limit - 1) / limit

		query := fmt.Sprintf(`
			SELECT id, merchant_id, amount, trx_id,
			       partner_reference_no, reference_no,
			       status, transaction_date, paid_date
			FROM transactions
			WHERE reference_no ILIKE $1
			ORDER BY %s %s
			LIMIT $2 OFFSET $3
		`, sortField, order)

		rows, err := db.Query(query, search, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var transactions []Transaction

		for rows.Next() {
			var t Transaction

			err := rows.Scan(
				&t.ID,
				&t.MerchantID,
				&t.Amount,
				&t.TrxID,
				&t.PartnerReferenceNo,
				&t.ReferenceNo,
				&t.Status,
				&t.TransactionDate,
				&t.PaidDate,
			)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			transactions = append(transactions, t)
		}

		response := Response{
			Data: transactions,
			Meta: Meta{
				Page:       page,
				Limit:      limit,
				Total:      total,
				TotalPages: totalPages,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GetTransactionByReference(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ref := chi.URLParam(r, "referenceNo")
		if ref == "" {
			http.Error(w, "referenceNo is required", http.StatusBadRequest)
			return
		}

		var t Transaction
		query := `
			SELECT id, merchant_id, amount, trx_id,
			       partner_reference_no, reference_no,
			       status, transaction_date, paid_date
			FROM transactions
			WHERE reference_no = $1
		`
		row := db.QueryRow(query, ref)
		err := row.Scan(
			&t.ID, &t.MerchantID, &t.Amount, &t.TrxID,
			&t.PartnerReferenceNo, &t.ReferenceNo,
			&t.Status, &t.TransactionDate, &t.PaidDate,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Transaction not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(t)
	}
}