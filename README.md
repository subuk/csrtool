# CSRTool

A tool for generating private keys and Certificate Signing Requests (CSRs) using ASN.1 directly.

## Features

- Generate private keys with RSA and ECDSA support
- Create CSRs with ASN.1 encoding
- Web interface for easy CSR generation
- CLI tool for automation and scripting
- PKCS#10 standard compliance

## Installation

### CLI Tool

```bash
go install github.com/subuk/csrtool/cmd/csrtool@latest
```

### Web Interface

1. Clone the repository:
```bash
git clone https://github.com/subuk/csrtool.git
cd csrtool
```

2. Build and run the web interface:
```bash
make web
```

The web interface will be available at http://localhost:3000

## Usage

### CLI Tool

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

Available key types:
- `rsa2048`: RSA with 2048-bit key
- `rsa4096`: RSA with 4096-bit key
- `ec256`: ECDSA with P-256 curve
- `ec384`: ECDSA with P-384 curve

### Web Interface

1. Open http://localhost:3000 in your browser
2. Fill out the CSR form with your desired information
3. Click "Generate CSR" to create your private key and CSR
4. The generated private key and CSR will be displayed on the page

The web interface uses WebAssembly to generate CSRs entirely in your browser, ensuring that private keys never leave your device.

## Building from Source

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

## Development

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

## License

MIT

---

<small>Built with [Cursor](https://cursor.sh)</small>
