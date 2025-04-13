# CSRTool

[![Go Report Card](https://goreportcard.com/badge/github.com/subuk/csrtool)](https://goreportcard.com/report/github.com/subuk/csrtool)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/subuk/csrtool.svg)](https://pkg.go.dev/github.com/subuk/csrtool)
[![Web Demo](https://img.shields.io/badge/Demo-Try%20Online-blue)](https://subuk.github.io/csrtool/)

CSRTool is a Go library and command-line tool for generating private keys and Certificate Signing Requests (CSRs) using ASN.1 directly, without relying on the standard library's x509 package. This library was created to address limitations in the Go standard library's PKCS#10 implementation, particularly around handling special attributes in CSRs (see [golang/go#15995](https://github.com/golang/go/issues/15995)).

## ✨ Features

- 🔐 Generate private keys with RSA and ECDSA support
- 📝 Create CSRs with ASN.1 encoding
- 🌐 Web interface for easy CSR generation
- 💻 CLI tool for automation and scripting
- 📋 PKCS#10 standard compliance
- 🔒 Private keys never leave your device

## 🚀 Quick Start

### Web Interface

Try the web interface directly in your browser: [https://subuk.github.io/csrtool/](https://subuk.github.io/csrtool/)

Or run it locally:
```bash
git clone https://github.com/subuk/csrtool.git
cd csrtool
make web
```
The web interface will be available at http://localhost:3000

### CLI Tool

Install the CLI tool:
```bash
go install github.com/subuk/csrtool/cmd/csrtool@latest
```

Generate a new private key and CSR:
```bash
csrtool generate example.com \
  --key-type rsa2048 \
  --country US \
  --state California \
  --locality San Francisco \
  --org "Example Inc" \
  --org-unit IT \
  --email admin@example.com
```

## 🔧 Key Types

| Type     | Description                |
|----------|----------------------------|
| rsa2048  | RSA with 2048-bit key      |
| rsa4096  | RSA with 4096-bit key      |
| ec256    | ECDSA with P-256 curve     |
| ec384    | ECDSA with P-384 curve     |

## 📚 Library Usage

Here's how to use the CSRTool library in your Go code:

```go
package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509/pkix"
    "fmt"
    "log"

    "github.com/subuk/csrtool/pkg/csrtool"
)

func main() {
    // Generate a private key
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
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
        "challenge-password", // Optional challenge password
    )
    if err != nil {
        log.Fatalf("Failed to generate CSR: %v", err)
    }

    fmt.Printf("Generated CSR:\n%s\n", string(csrBytes))
}
```

### Advanced Usage

You can also use ECDSA keys:

```go
package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/x509/pkix"
    "fmt"
    "log"

    "github.com/subuk/csrtool/pkg/csrtool"
)

func main() {
    // Generate an ECDSA private key
    privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        log.Fatalf("Failed to generate private key: %v", err)
    }

    // Create a subject with optional fields
    subject := pkix.Name{
        CommonName: "example.com",
    }

    // Generate a CSR with optional fields
    csrBytes, err := csrtool.GenerateCSR(
        privateKey,
        subject,
        []string{"example.com"},
        "", // No challenge password
    )
    if err != nil {
        log.Fatalf("Failed to generate CSR: %v", err)
    }

    fmt.Printf("Generated CSR:\n%s\n", string(csrBytes))
}
```

## 🛠️ Building from Source

### CLI Tool
```bash
make build
```
The binary will be created in the `bin/` directory.

### Web Interface
```bash
make web-build
```
The built files will be in the `web/dist/` directory.

## 🧪 Development

### Running the Web Development Server
```bash
make web
```
This will:
1. Build the WASM module
2. Install web dependencies
3. Start the development server at http://localhost:3000

### Cleaning Build Artifacts
```bash
make clean
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<small>Built with [Cursor](https://cursor.sh)</small>
