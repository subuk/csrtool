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

type GenerateCmd struct {
	CommonName        string   `arg:"" help:"Common Name (CN) for the certificate"`
	KeyType           string   `short:"t" default:"rsa2048" enum:"rsa2048,rsa4096,ec256,ec384" help:"Key type"`
	OutputKey         string   `short:"k" default:"private.key" help:"Output path for private key"`
	OutputCSR         string   `short:"c" default:"csr.pem" help:"Output path for CSR"`
	Country           string   `short:"C" default:"US" help:"Country code"`
	State             string   `short:"S" default:"California" help:"State or province"`
	Locality          string   `short:"L" default:"San Francisco" help:"Locality or city"`
	Org               string   `short:"O" default:"Example Inc" help:"Organization name"`
	OrgUnit           string   `short:"U" default:"IT" help:"Organizational unit"`
	Email             string   `short:"E" default:"" help:"Email address"`
	DNSNames          []string `help:"DNS names for the certificate"`
	ChallengePassword string   `short:"p" help:"Challenge password for the CSR"`
}

var cli struct {
	Generate GenerateCmd `cmd:"" help:"Generate a new private key and CSR"`
	Version  struct{}    `cmd:"" help:"Show version information"`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("csrtool"),
		kong.Description("A tool for generating private keys and CSRs"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)

	switch ctx.Command() {
	case "generate <common-name>":
		if err := generate(cli.Generate); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "version":
		fmt.Printf("csrtool version %s\n", version)
		fmt.Printf("Build time: %s\n", buildTime)
		fmt.Printf("Git commit: %s\n", gitCommit)
	default:
		ctx.PrintUsage(true)
		os.Exit(1)
	}
}

func generate(cmd GenerateCmd) error {
	// Generate private key
	var privateKey interface{}
	var err error

	switch strings.ToLower(cmd.KeyType) {
	case "rsa2048":
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	case "rsa4096":
		privateKey, err = rsa.GenerateKey(rand.Reader, 4096)
	case "ec256":
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "ec384":
		privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	default:
		return fmt.Errorf("unsupported key type: %s", cmd.KeyType)
	}

	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	// Save private key
	keyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %v", err)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})

	if err := os.WriteFile(cmd.OutputKey, keyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write private key: %v", err)
	}

	// Create subject
	subject := pkix.Name{
		CommonName:         cmd.CommonName,
		Organization:       []string{cmd.Org},
		OrganizationalUnit: []string{cmd.OrgUnit},
		Country:            []string{cmd.Country},
		Province:           []string{cmd.State},
		Locality:           []string{cmd.Locality},
	}

	// Generate CSR
	csrBytes, err := csrtool.GenerateCSR(
		privateKey,
		subject,
		cmd.DNSNames,
		cmd.ChallengePassword,
	)
	if err != nil {
		return fmt.Errorf("failed to generate CSR: %v", err)
	}

	// Save CSR
	if err := os.WriteFile(cmd.OutputCSR, csrBytes, 0644); err != nil {
		return fmt.Errorf("failed to write CSR: %v", err)
	}

	fmt.Printf("Private key saved to: %s\n", cmd.OutputKey)
	fmt.Printf("CSR saved to: %s\n", cmd.OutputCSR)
	return nil
}
