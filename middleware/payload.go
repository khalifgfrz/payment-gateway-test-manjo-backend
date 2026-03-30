package middleware

func GenerateQRPayload(body map[string]interface{}) string {
	amount := body["amount"].(map[string]interface{})
	return body["partnerReferenceNo"].(string) +
		body["merchantId"].(string) +
		amount["value"].(string)
}

func PaymentPayload(body map[string]interface{}) string {
	amount := body["amount"].(map[string]interface{})
	return body["originalReferenceNo"].(string) +
		amount["value"].(string)
}