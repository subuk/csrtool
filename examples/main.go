package main

import (
	"crypto/x509/pkix"
	"fmt"
	"log"
	"os"

	"github.com/subuk/csrtool/pkg/csrtool"
)

func main() {
	// Generate an RSA private key
	privateKey, err := csrtool.GeneratePrivateKey(csrtool.RSA2048)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// Create a subject for the CSR
	subject := pkix.Name{
		CommonName:         "example.com",
		Organization:       []string{"Example Organization"},
		OrganizationalUnit: []string{"Example Unit"},
		Country:            []string{"US"},
		Province:           []string{"California"},
		Locality:           []string{"San Francisco"},
	}

	// Generate a CSR
	csrBytes, err := csrtool.GenerateCSR(
		privateKey,
		subject,
		[]string{"example.com", "www.example.com"},
		"mySecurePassword123",
	)
	if err != nil {
		log.Fatalf("Failed to generate CSR: %v", err)
	}

	// Save the CSR to a file
	err = os.WriteFile("example.csr", csrBytes, 0644)
	if err != nil {
		log.Fatalf("Failed to write CSR to file: %v", err)
	}

	fmt.Println("CSR generated successfully and saved to example.csr")
}
