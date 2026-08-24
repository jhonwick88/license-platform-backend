package token

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)

func GenerateRSAKeyPair() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

type LicenseTokenClaims struct {
	LicenseID      string                 `json:"license_id"`
	ProductID      string                 `json:"product_id"`
	CustomerID     string                 `json:"customer_id"`
	PlanID         string                 `json:"plan_id"`
	InstallationID string                 `json:"installation_id"`
	Features       map[string]interface{} `json:"features"`
}

func SignLicenseToken(claims LicenseTokenClaims, privateKey *rsa.PrivateKey) (string, error) {
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	hashed := sha256.Sum256(payload)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)

	return encodedPayload + "." + encodedSignature, nil
}
