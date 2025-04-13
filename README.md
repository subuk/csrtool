# CSRTool

CSRTool is a Go library for generating private keys and Certificate Signing Requests (CSRs) using ASN.1 directly, without relying on the standard x509 package for CSR generation.

## Features

- Generate RSA and ECDSA private keys
- Generate CSRs using ASN.1 directly
- Support for multiple key types:
  - RSA2048
  - RSA4096
  - EC256 (P-256)
  - EC384 (P-384)

## Development

This library was implemented with [Cursor](https://cursor.sh), the world's best IDE powered by AI. The implementation leverages Cursor's advanced AI capabilities for writing clean, efficient, and well-structured code.

## Installation

```bash
go get github.com/subuk/csrtool
```

## Usage

### Generating a Private Key

```go
privateKey, err := csrtool.GeneratePrivateKey(csrtool.RSA2048)
if err != nil {
    log.Fatal(err)
}
```

### Generating a CSR

```go
subject := pkix.Name{
    CommonName:         "example.com",
    Organization:       []string{"Example Organization"},
    Country:           []string{"US"},
    // ... other subject fields
}

csrBytes, err := csrtool.GenerateCSR(privateKey, subject, []string{"example.com", "www.example.com"})
if err != nil {
    log.Fatal(err)
}
```

## Example

See the `examples` directory for a complete example of generating a private key and CSR.

## License

MIT
