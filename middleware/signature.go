package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"payment-gateway-test-manjo/utils"
)

func SignatureMiddleware(buildPayload func(map[string]interface{}) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			signature := r.Header.Get("X-SIGNATURE")
			if signature == "" {
				http.Error(w, "Missing signature", http.StatusUnauthorized)
				return
			}

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Invalid body", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			var body map[string]interface{}
			json.Unmarshal(bodyBytes, &body)

			payload := buildPayload(body)

			expected := utils.GenerateSignature(os.Getenv("SECRET_KEY"), payload)

			if signature != expected {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}