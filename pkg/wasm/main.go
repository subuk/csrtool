//go:build js && wasm

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"syscall/js"

	"github.com/subuk/csrtool/pkg/csrtool"
)

// CSRRequest represents the JSON structure expected from the web interface
type CSRRequest struct {
	CommonName        string   `json:"commonName"`
	KeyType           string   `json:"keyType"`
	Country           string   `json:"country"`
	State             string   `json:"state"`
	Locality          string   `json:"locality"`
	Org               string   `json:"org"`
	OrgUnit           string   `json:"orgUnit"`
	Email             string   `json:"email"`
	DNSNames          []string `json:"dnsNames"`
	ChallengePassword string   `json:"challengePassword"`
}

// CSRResponse represents the JSON structure returned to the web interface
type CSRResponse struct {
	PrivateKey string `json:"privateKey"`
	CSR        string `json:"csr"`
	Error      string `json:"error,omitempty"`
}

func generateCSR(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return js.ValueOf(map[string]interface{}{
			"privateKey": "",
			"csr":        "",
			"error":      "invalid number of arguments",
		})
	}

	var req CSRRequest
	if err := json.Unmarshal([]byte(args[0].String()), &req); err != nil {
		return js.ValueOf(map[string]interface{}{
			"privateKey": "",
			"csr":        "",
			"error":      err.Error(),
		})
	}

	// Generate private key
	var privateKey interface{}
	var err error

	switch req.KeyType {
	case "rsa2048":
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	case "rsa4096":
		privateKey, err = rsa.GenerateKey(rand.Reader, 4096)
	case "ec256":
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "ec384":
		privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	default:
		return js.ValueOf(map[string]interface{}{
			"privateKey": "",
			"csr":        "",
			"error":      "unsupported key type",
		})
	}

	if err != nil {
		return js.ValueOf(map[string]interface{}{
			"privateKey": "",
			"csr":        "",
			"error":      err.Error(),
		})
	}

	// Create subject
	subject := pkix.Name{
		CommonName:         req.CommonName,
		Organization:       []string{req.Org},
		OrganizationalUnit: []string{req.OrgUnit},
		Country:            []string{req.Country},
		Province:           []string{req.State},
		Locality:           []string{req.Locality},
	}

	// Generate CSR
	csrBytes, err := csrtool.GenerateCSR(
		privateKey,
		subject,
		req.DNSNames,
		req.ChallengePassword,
	)
	if err != nil {
		return js.ValueOf(map[string]interface{}{
			"privateKey": "",
			"csr":        "",
			"error":      err.Error(),
		})
	}

	// Convert private key to PEM
	keyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return js.ValueOf(map[string]interface{}{
			"privateKey": "",
			"csr":        "",
			"error":      err.Error(),
		})
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})

	return js.ValueOf(map[string]interface{}{
		"privateKey": string(keyPEM),
		"csr":        string(csrBytes),
	})
}

func main() {
	js.Global().Set("generateCSR", js.FuncOf(generateCSR))
	// Keep the Go program alive
	select {}
}
