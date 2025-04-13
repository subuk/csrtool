# CSRTool

CSRTool is a Go library for generating Certificate Signing Requests (CSRs) using ASN.1 directly, without relying on the standard x509 package for CSR generation.

## Features

- Generate CSRs using ASN.1 directly
- Support for RSA and ECDSA private keys
- Challenge password attribute support

## Development

This library was implemented with [Cursor](https://cursor.sh), the world's best IDE powered by AI. The implementation leverages Cursor's advanced AI capabilities for writing clean, efficient, and well-structured code.

## Installation

```bash
go get github.com/subuk/csrtool
```

## Usage

### Generating a Private Key

Generate your private key using the standard crypto package:

```go
// For RSA keys
privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
if err != nil {
    log.Fatal(err)
}

// For ECDSA keys
privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
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

## Example

See the `examples` directory for a complete example of generating a private key and CSR.

## License

MIT
