# CSRTool

CSRTool is a Go library for generating Certificate Signing Requests (CSRs) using ASN.1 directly, without relying on the standard x509 package for CSR generation.

## Features

- Generate CSRs using ASN.1 directly
- Support for RSA and ECDSA private keys
- Challenge password attribute support
- Command-line interface for easy usage

## Development

This library was implemented with [Cursor](https://cursor.sh), the world's best IDE powered by AI. The implementation leverages Cursor's advanced AI capabilities for writing clean, efficient, and well-structured code.

## Installation

```bash
go get github.com/subuk/csrtool
```

## Usage

### Library Usage

```go
subject := pkix.Name{
    CommonName:         "example.com",
    Organization:       []string{"Example Organization"},
    Country:           []string{"US"},
    // ... other subject fields
}

csrBytes, err := csrtool.GenerateCSR(
    privateKey,
    subject,
    []string{"example.com", "www.example.com"},
    "challengePassword", // optional challenge password
)
if err != nil {
    log.Fatal(err)
}
```

### Command Line Interface

The CLI tool provides an easy way to generate private keys and CSRs:

```bash
# Generate a new RSA 2048-bit key and CSR
csrtool generate \
    --common-name "example.com" \
    --organization "Example Organization" \
    --country "US" \
    --dns-names "example.com" "www.example.com" \
    --challenge-password "hello321"

# Generate an ECDSA P-256 key and CSR
csrtool generate \
    --key-type ec256 \
    --common-name "example.com" \
    --organization "Example Organization" \
    --country "US" \
    --dns-names "example.com" "www.example.com"

# Show help
csrtool --help
```

Available options:
- `--key-type`: Type of key to generate (rsa2048, rsa4096, ec256, ec384)
- `--output-key`: Output file for the private key
- `--output-csr`: Output file for the CSR
- `--common-name`: Common Name (CN) for the certificate
- `--organization`: Organization (O) for the certificate
- `--organizational-unit`: Organizational Unit (OU) for the certificate
- `--country`: Country (C) for the certificate
- `--province`: Province/State (ST) for the certificate
- `--locality`: Locality (L) for the certificate
- `--dns-names`: DNS names for the certificate
- `--challenge-password`: Challenge password for the CSR

## Example

See the `examples` directory for a complete example of generating a private key and CSR.

## License

MIT
