package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/subuk/csrtool/pkg/csrtool"
)

// Version information set during build
var (
	version   string
	buildTime string
	gitCommit string
)

var cli struct {
	Generate struct {
		KeyType            string   `help:"Type of key to generate (rsa2048, rsa4096, ec256, ec384)" enum:"rsa2048,rsa4096,ec256,ec384" default:"rsa2048"`
		OutputKey          string   `help:"Output file for the private key" type:"path" default:"private.key"`
		OutputCSR          string   `help:"Output file for the CSR" type:"path" default:"request.csr"`
		CommonName         string   `help:"Common Name (CN) for the certificate" required:""`
		Organization       []string `help:"Organization (O) for the certificate"`
		OrganizationalUnit []string `help:"Organizational Unit (OU) for the certificate"`
		Country            []string `help:"Country (C) for the certificate"`
		Province           []string `help:"Province/State (ST) for the certificate"`
		Locality           []string `help:"Locality (L) for the certificate"`
		DNSNames           []string `help:"DNS names for the certificate"`
		ChallengePassword  string   `help:"Challenge password for the CSR"`
	} `cmd:"" help:"Generate a new private key and CSR"`

	Version struct{} `cmd:"" help:"Show version information"`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("csrtool"),
		kong.Description("A tool for generating private keys and CSRs"),
		kong.UsageOnError(),
	)

	switch ctx.Command() {
	case "generate":
		generate()
	case "version":
		fmt.Printf("csrtool version %s\n", version)
		fmt.Printf("Build time: %s\n", buildTime)
		fmt.Printf("Git commit: %s\n", gitCommit)
	default:
		ctx.FatalIfErrorf(fmt.Errorf("unknown command"))
	}
}

func generate() {
	// Generate private key
	var privateKey interface{}
	var err error

	switch strings.ToLower(cli.Generate.KeyType) {
	case "rsa2048":
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	case "rsa4096":
		privateKey, err = rsa.GenerateKey(rand.Reader, 4096)
	case "ec256":
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "ec384":
		privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	default:
		log.Fatalf("Unsupported key type: %s", cli.Generate.KeyType)
	}

	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// Save private key
	keyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Failed to marshal private key: %v", err)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})

	if err := os.WriteFile(cli.Generate.OutputKey, keyPEM, 0600); err != nil {
		log.Fatalf("Failed to write private key: %v", err)
	}

	// Create subject
	subject := pkix.Name{
		CommonName:         cli.Generate.CommonName,
		Organization:       cli.Generate.Organization,
		OrganizationalUnit: cli.Generate.OrganizationalUnit,
		Country:            cli.Generate.Country,
		Province:           cli.Generate.Province,
		Locality:           cli.Generate.Locality,
	}

	// Generate CSR
	csrBytes, err := csrtool.GenerateCSR(
		privateKey,
		subject,
		cli.Generate.DNSNames,
		cli.Generate.ChallengePassword,
	)
	if err != nil {
		log.Fatalf("Failed to generate CSR: %v", err)
	}

	// Save CSR
	if err := os.WriteFile(cli.Generate.OutputCSR, csrBytes, 0644); err != nil {
		log.Fatalf("Failed to write CSR: %v", err)
	}

	fmt.Printf("Private key saved to: %s\n", cli.Generate.OutputKey)
	fmt.Printf("CSR saved to: %s\n", cli.Generate.OutputCSR)
}
